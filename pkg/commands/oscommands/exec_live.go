package oscommands

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type CredentialName string

const (
	Password   CredentialName = "password"
	Username                  = "username"
	Passphrase                = "passphrase"
)

// RunAndDetectCredentialRequest detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fillin a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (self *cmdObjRunner) RunAndDetectCredentialRequest(cmdObj ICmdObj, promptUserForCredential func(CredentialName) string) error {
	ttyText := ""
	err := self.RunCommandWithOutputLive(cmdObj, func(word string) string {
		ttyText = ttyText + " " + word

		prompts := map[string]CredentialName{
			`.+'s password:`:                         "password",
			`Password\s*for\s*'.+':`:                 "password",
			`Username\s*for\s*'.+':`:                 "username",
			`Enter\s*passphrase\s*for\s*key\s*'.+':`: "passphrase",
		}

		for pattern, askFor := range prompts {
			if match, _ := regexp.MatchString(pattern, ttyText); match {
				ttyText = ""
				return promptUserForCredential(askFor)
			}
		}

		return ""
	})
	return err
}

type cmdHandler struct {
	stdoutPipe io.Reader
	stdinPipe  io.Writer
	close      func() error
}

// RunCommandWithOutputLiveAux runs a command and return every word that gets written in stdout
// Output is a function that executes by every word that gets read by bufio
// As return of output you need to give a string that will be written to stdin
// NOTE: If the return data is empty it won't write anything to stdin
func (self *cmdObjRunner) RunCommandWithOutputLiveAux(
	cmdObj ICmdObj,
	// handleOutput takes a word from stdout and returns a string to be written to stdin.
	// See RunAndDetectCredentialRequest above for how this is used to check for a username/password request
	handleOutput func(string) string,
	startCmd func(cmd *exec.Cmd) (*cmdHandler, error),
) error {
	cmdWriter := self.guiIO.newCmdWriterFn()
	self.log.WithField("command", cmdObj.ToString()).Info("RunCommand")
	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}
	cmd := cmdObj.AddEnvVars("LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8").GetCmd()

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(cmdWriter, &stderr)

	handler, err := startCmd(cmd)
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
		scanner := bufio.NewScanner(tr)
		scanner.Split(scanWordsWithNewLines)
		for scanner.Scan() {
			text := scanner.Text()
			output := strings.Trim(text, " ")
			toInput := handleOutput(output)
			if toInput != "" {
				_, _ = handler.stdinPipe.Write([]byte(toInput))
			}
		}
	})

	err = cmd.Wait()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

// scanWordsWithNewLines is a copy of bufio.ScanWords but this also captures new lines
// For specific comments about this function take a look at: bufio.ScanWords
func scanWordsWithNewLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}

// isSpace is also copied from the bufio package and has been modified to also captures new lines
// For specific comments about this function take a look at: bufio.isSpace
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		switch r {
		case ' ', '\t', '\v', '\f':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}
