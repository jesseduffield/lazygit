package oscommands

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type CredentialType int

const (
	Password CredentialType = iota
	Username
	Passphrase
)

type cmdHandler struct {
	stdoutPipe io.Reader
	stdinPipe  io.Writer
	close      func() error
}

// RunAndDetectCredentialRequest detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fillin a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (self *cmdObjRunner) RunAndDetectCredentialRequest(
	cmdObj ICmdObj,
	promptUserForCredential func(CredentialType) string,
) error {
	self.log.Warn("HERE")
	cmdWriter := self.guiIO.newCmdWriterFn()
	self.log.WithField("command", cmdObj.ToString()).Info("RunCommand")
	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}
	cmd := cmdObj.AddEnvVars("LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8").GetCmd()

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

	tr := io.TeeReader(handler.stdoutPipe, cmdWriter)

	go utils.Safe(func() {
		self.processOutput(tr, handler.stdinPipe, promptUserForCredential)
	})

	err = cmd.Wait()
	if err != nil {
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
