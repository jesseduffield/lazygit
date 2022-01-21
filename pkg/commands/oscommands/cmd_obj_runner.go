package oscommands

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type ICmdObjRunner interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error
}

type CredentialType int

const (
	Password CredentialType = iota
	Username
	Passphrase
)

type cmdObjRunner struct {
	log   *logrus.Entry
	guiIO *guiIO
}

var _ ICmdObjRunner = &cmdObjRunner{}

func (self *cmdObjRunner) Run(cmdObj ICmdObj) error {
	if cmdObj.GetCredentialStrategy() != NONE {
		return self.runWithCredentialHandling(cmdObj)
	}

	if cmdObj.ShouldStreamOutput() {
		return self.runAndStream(cmdObj)
	}

	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *cmdObjRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	if cmdObj.ShouldStreamOutput() {
		err := self.runAndStream(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutput
		// but ignore the output, hence why we've got this check here.
		return "", err
	}

	if cmdObj.GetCredentialStrategy() != NONE {
		err := self.runWithCredentialHandling(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutput
		// but ignore the output, hence why we've got this check here.
		return "", err
	}

	self.log.WithField("command", cmdObj.ToString()).Debug("RunCommand")

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		self.log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

func (self *cmdObjRunner) RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	if cmdObj.GetCredentialStrategy() != NONE {
		return errors.New("cannot call RunAndProcessLines with credential strategy. If you're seeing this then a contributor to Lazygit has accidentally called this method! Please raise an issue")
	}

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	cmd := cmdObj.GetCmd()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(bufio.ScanLines)
	if err := cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		stop, err := onLine(line)
		if err != nil {
			return err
		}
		if stop {
			_ = cmd.Process.Kill()
			break
		}
	}

	_ = cmd.Wait()

	return nil
}

// Whenever we're asked for a password we just enter a newline, which will
// eventually cause the command to fail.
var failPromptFn = func(CredentialType) string { return "\n" }

func (self *cmdObjRunner) runWithCredentialHandling(cmdObj ICmdObj) error {
	var promptFn func(CredentialType) string

	switch cmdObj.GetCredentialStrategy() {
	case PROMPT:
		promptFn = self.guiIO.promptForCredentialFn
	case FAIL:
		promptFn = failPromptFn
	case NONE:
		// we should never land here
		return errors.New("runWithCredentialHandling called but cmdObj does not have a a credential strategy")
	}

	return self.runAndDetectCredentialRequest(cmdObj, promptFn)
}

func (self *cmdObjRunner) logCmdObj(cmdObj ICmdObj) {
	self.guiIO.logCommandFn(cmdObj.ToString(), true)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		if outputString == "" {
			return "", utils.WrapError(err)
		}
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}

type cmdHandler struct {
	stdoutPipe io.Reader
	stdinPipe  io.Writer
	close      func() error
}

func (self *cmdObjRunner) runAndStream(cmdObj ICmdObj) error {
	return self.runAndStreamAux(cmdObj, func(handler *cmdHandler, cmdWriter io.Writer) {
		go func() {
			_, _ = io.Copy(cmdWriter, handler.stdoutPipe)
		}()
	})
}

// runAndDetectCredentialRequest detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fillin a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (self *cmdObjRunner) runAndDetectCredentialRequest(
	cmdObj ICmdObj,
	promptUserForCredential func(CredentialType) string,
) error {
	// setting the output to english so we can parse it for a username/password request
	cmdObj.AddEnvVars("LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")

	return self.runAndStreamAux(cmdObj, func(handler *cmdHandler, cmdWriter io.Writer) {
		tr := io.TeeReader(handler.stdoutPipe, cmdWriter)

		go utils.Safe(func() {
			self.processOutput(tr, handler.stdinPipe, promptUserForCredential)
		})
	})
}

func (self *cmdObjRunner) runAndStreamAux(
	cmdObj ICmdObj,
	onRun func(*cmdHandler, io.Writer),
) error {
	cmdWriter := self.guiIO.newCmdWriterFn()

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}
	self.log.WithField("command", cmdObj.ToString()).Info("RunCommand")
	cmd := cmdObj.GetCmd()

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(cmdWriter, &stderr)

	handler, err := self.getCmdHandler(cmd)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := handler.close(); closeErr != nil {
			self.log.Error(closeErr)
		}
	}()

	onRun(handler, cmdWriter)

	err = cmd.Wait()
	if err != nil {
		errStr := stderr.String()
		if cmdObj.ShouldIgnoreEmptyError() && errStr == "" {
			return nil
		}
		return errors.New(stderr.String())
	}

	return nil
}

func (self *cmdObjRunner) processOutput(reader io.Reader, writer io.Writer, promptUserForCredential func(CredentialType) string) {
	checkForCredentialRequest := self.getCheckForCredentialRequestFunc()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		newBytes := scanner.Bytes()
		askFor, ok := checkForCredentialRequest(newBytes)
		if ok {
			toInput := promptUserForCredential(askFor)
			// If the return data is empty we don't write anything to stdin
			if toInput != "" {
				_, _ = writer.Write([]byte(toInput))
			}
		}
	}
}

// having a function that returns a function because we need to maintain some state inbetween calls hence the closure
func (self *cmdObjRunner) getCheckForCredentialRequestFunc() func([]byte) (CredentialType, bool) {
	var ttyText strings.Builder
	// this function takes each word of output from the command and builds up a string to see if we're being asked for a password
	return func(newBytes []byte) (CredentialType, bool) {
		_, err := ttyText.Write(newBytes)
		if err != nil {
			self.log.Error(err)
		}

		prompts := map[string]CredentialType{
			`Password:`:                              Password,
			`.+'s password:`:                         Password,
			`Password\s*for\s*'.+':`:                 Password,
			`Username\s*for\s*'.+':`:                 Username,
			`Enter\s*passphrase\s*for\s*key\s*'.+':`: Passphrase,
		}

		for pattern, askFor := range prompts {
			if match, _ := regexp.MatchString(pattern, ttyText.String()); match {
				ttyText.Reset()
				return askFor, true
			}
		}

		return 0, false
	}
}
