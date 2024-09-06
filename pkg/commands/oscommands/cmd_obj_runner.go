package oscommands

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type ICmdObjRunner interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunWithOutputs(cmdObj ICmdObj) (string, string, error)
	RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error
}

type cmdObjRunner struct {
	log   *logrus.Entry
	guiIO *guiIO
}

var _ ICmdObjRunner = &cmdObjRunner{}

func (self *cmdObjRunner) Run(cmdObj ICmdObj) error {
	if cmdObj.Mutex() != nil {
		cmdObj.Mutex().Lock()
		defer cmdObj.Mutex().Unlock()
	}

	if cmdObj.GetCredentialStrategy() != NONE {
		return self.runWithCredentialHandling(cmdObj)
	}

	if cmdObj.ShouldStreamOutput() {
		return self.runAndStream(cmdObj)
	}

	_, err := self.RunWithOutputAux(cmdObj)
	return err
}

func (self *cmdObjRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	if cmdObj.Mutex() != nil {
		cmdObj.Mutex().Lock()
		defer cmdObj.Mutex().Unlock()
	}

	if cmdObj.GetCredentialStrategy() != NONE {
		err := self.runWithCredentialHandling(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutput
		// but ignore the output, hence why we've got this check here.
		return "", err
	}

	if cmdObj.ShouldStreamOutput() {
		err := self.runAndStream(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutput
		// but ignore the output, hence why we've got this check here.
		return "", err
	}

	return self.RunWithOutputAux(cmdObj)
}

func (self *cmdObjRunner) RunWithOutputs(cmdObj ICmdObj) (string, string, error) {
	if cmdObj.Mutex() != nil {
		cmdObj.Mutex().Lock()
		defer cmdObj.Mutex().Unlock()
	}

	if cmdObj.GetCredentialStrategy() != NONE {
		err := self.runWithCredentialHandling(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutputs
		// but ignore the output, hence why we've got this check here.
		return "", "", err
	}

	if cmdObj.ShouldStreamOutput() {
		err := self.runAndStream(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutputs
		// but ignore the output, hence why we've got this check here.
		return "", "", err
	}

	return self.RunWithOutputsAux(cmdObj)
}

func (self *cmdObjRunner) RunWithOutputAux(cmdObj ICmdObj) (string, error) {
	self.log.WithField("command", cmdObj.ToString()).Debug("RunCommand")

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	t := time.Now()
	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		self.log.WithField("command", cmdObj.ToString()).Error(output)
	}

	self.log.Infof("%s (%s)", cmdObj.ToString(), time.Since(t))

	return output, err
}

func (self *cmdObjRunner) RunWithOutputsAux(cmdObj ICmdObj) (string, string, error) {
	self.log.WithField("command", cmdObj.ToString()).Debug("RunCommand")

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	t := time.Now()
	var outBuffer, errBuffer bytes.Buffer
	cmd := cmdObj.GetCmd()
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	err := cmd.Run()

	self.log.Infof("%s (%s)", cmdObj.ToString(), time.Since(t))

	stdout := outBuffer.String()
	stderr, err := sanitisedCommandOutput(errBuffer.Bytes(), err)
	if err != nil {
		self.log.WithField("command", cmdObj.ToString()).Error(stderr)
	}

	return stdout, stderr, err
}

func (self *cmdObjRunner) RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	if cmdObj.Mutex() != nil {
		cmdObj.Mutex().Lock()
		defer cmdObj.Mutex().Unlock()
	}

	if cmdObj.GetCredentialStrategy() != NONE {
		return errors.New("cannot call RunAndProcessLines with credential strategy. If you're seeing this then a contributor to Lazygit has accidentally called this method! Please raise an issue")
	}

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}
	t := time.Now()

	cmd := cmdObj.GetCmd()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(utils.ScanLinesAndTruncateWhenLongerThanBuffer(bufio.MaxScanTokenSize))
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
			_ = Kill(cmd)
			break
		}
	}

	if scanner.Err() != nil {
		_ = Kill(cmd)
		return scanner.Err()
	}

	_ = cmd.Wait()

	self.log.Infof("%s (%s)", cmdObj.ToString(), time.Since(t))

	return nil
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

func (self *cmdObjRunner) runAndStreamAux(
	cmdObj ICmdObj,
	onRun func(*cmdHandler, io.Writer),
) error {
	cmdWriter := self.guiIO.newCmdWriterFn()

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}
	self.log.WithField("command", cmdObj.ToString()).Debug("RunCommand")
	cmd := cmdObj.GetCmd()

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(cmdWriter, &stderr)

	handler, err := self.getCmdHandler(cmd)
	if err != nil {
		return err
	}

	var stdout bytes.Buffer
	handler.stdoutPipe = io.TeeReader(handler.stdoutPipe, &stdout)

	defer func() {
		if closeErr := handler.close(); closeErr != nil {
			self.log.Error(closeErr)
		}
	}()

	t := time.Now()

	onRun(handler, cmdWriter)

	err = cmd.Wait()

	self.log.Infof("%s (%s)", cmdObj.ToString(), time.Since(t))

	if err != nil {
		errStr := stderr.String()
		if errStr != "" {
			return errors.New(errStr)
		}

		if cmdObj.ShouldIgnoreEmptyError() {
			return nil
		}
		stdoutStr := stdout.String()
		if stdoutStr != "" {
			return errors.New(stdoutStr)
		}
		return errors.New("Command exited with non-zero exit code, but no output")
	}

	return nil
}

type CredentialType int

const (
	Password CredentialType = iota
	Username
	Passphrase
	PIN
	Token
)

// Whenever we're asked for a password we just enter a newline, which will
// eventually cause the command to fail.
var failPromptFn = func(CredentialType) <-chan string {
	ch := make(chan string)
	go func() {
		ch <- "\n"
	}()
	return ch
}

func (self *cmdObjRunner) runWithCredentialHandling(cmdObj ICmdObj) error {
	promptFn, err := self.getCredentialPromptFn(cmdObj)
	if err != nil {
		return err
	}

	return self.runAndDetectCredentialRequest(cmdObj, promptFn)
}

func (self *cmdObjRunner) getCredentialPromptFn(cmdObj ICmdObj) (func(CredentialType) <-chan string, error) {
	switch cmdObj.GetCredentialStrategy() {
	case PROMPT:
		return self.guiIO.promptForCredentialFn, nil
	case FAIL:
		return failPromptFn, nil
	default:
		// we should never land here
		return nil, errors.New("runWithCredentialHandling called but cmdObj does not have a credential strategy")
	}
}

// runAndDetectCredentialRequest detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fill in a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (self *cmdObjRunner) runAndDetectCredentialRequest(
	cmdObj ICmdObj,
	promptUserForCredential func(CredentialType) <-chan string,
) error {
	// setting the output to english so we can parse it for a username/password request
	cmdObj.AddEnvVars("LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")

	return self.runAndStreamAux(cmdObj, func(handler *cmdHandler, cmdWriter io.Writer) {
		tr := io.TeeReader(handler.stdoutPipe, cmdWriter)

		go utils.Safe(func() {
			self.processOutput(tr, handler.stdinPipe, promptUserForCredential, cmdObj.GetTask())
		})
	})
}

func (self *cmdObjRunner) processOutput(
	reader io.Reader,
	writer io.Writer,
	promptUserForCredential func(CredentialType) <-chan string,
	task gocui.Task,
) {
	checkForCredentialRequest := self.getCheckForCredentialRequestFunc()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		newBytes := scanner.Bytes()
		askFor, ok := checkForCredentialRequest(newBytes)
		if ok {
			responseChan := promptUserForCredential(askFor)
			if task != nil {
				task.Pause()
			}
			toInput := <-responseChan
			if task != nil {
				task.Continue()
			}
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
	prompts := map[string]CredentialType{
		`Password:`:                              Password,
		`.+'s password:`:                         Password,
		`Password\s*for\s*'.+':`:                 Password,
		`Username\s*for\s*'.+':`:                 Username,
		`Enter\s*passphrase\s*for\s*key\s*'.+':`: Passphrase,
		`Enter\s*PIN\s*for\s*.+\s*key\s*.+:`:     PIN,
		`.*2FA Token.*`:                          Token,
	}

	compiledPrompts := map[*regexp.Regexp]CredentialType{}
	for pattern, askFor := range prompts {
		compiledPattern := regexp.MustCompile(pattern)
		compiledPrompts[compiledPattern] = askFor
	}

	newlineRegex := regexp.MustCompile("\n")

	// this function takes each word of output from the command and builds up a string to see if we're being asked for a password
	return func(newBytes []byte) (CredentialType, bool) {
		_, err := ttyText.Write(newBytes)
		if err != nil {
			self.log.Error(err)
		}

		for pattern, askFor := range compiledPrompts {
			if match := pattern.Match([]byte(ttyText.String())); match {
				ttyText.Reset()
				return askFor, true
			}
		}

		if indices := newlineRegex.FindIndex([]byte(ttyText.String())); indices != nil {
			newText := []byte(ttyText.String()[indices[1]:])
			ttyText.Reset()
			ttyText.Write(newText)
		}
		return 0, false
	}
}
