package oscommands

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

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
	scanner.Split(scanWordsWithNewLines)
	for scanner.Scan() {
		text := scanner.Text()
		self.log.Info(text)
		output := strings.Trim(text, " ")
		askFor, ok := checkForCredentialRequest(output)
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
func (self *cmdObjRunner) getCheckForCredentialRequestFunc() func(string) (CredentialType, bool) {
	ttyText := ""
	// this function takes each word of output from the command and builds up a string to see if we're being asked for a password
	return func(word string) (CredentialType, bool) {
		ttyText = ttyText + " " + word

		prompts := map[string]CredentialType{
			`.+'s password:`:                         Password,
			`Password\s*for\s*'.+':`:                 Password,
			`Username\s*for\s*'.+':`:                 Username,
			`Enter\s*passphrase\s*for\s*key\s*'.+':`: Passphrase,
		}

		for pattern, askFor := range prompts {
			if match, _ := regexp.MatchString(pattern, ttyText); match {
				ttyText = ""
				return askFor, true
			}
		}

		return 0, false
	}
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
