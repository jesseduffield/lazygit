package oscommands

import (
	"bufio"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

// for use in testing

type FakeCmdObjRunner struct {
	t                *testing.T
	expectedCmds     []func(ICmdObj) (string, error)
	expectedCmdIndex int
}

var _ ICmdObjRunner = &FakeCmdObjRunner{}

func NewFakeRunner(t *testing.T) *FakeCmdObjRunner { //nolint:thelper
	return &FakeCmdObjRunner{t: t}
}

func (self *FakeCmdObjRunner) Run(cmdObj ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *FakeCmdObjRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	if self.expectedCmdIndex > len(self.expectedCmds)-1 {
		self.t.Errorf("ran too many commands. Unexpected command: `%s`", cmdObj.ToString())
		return "", errors.New("ran too many commands")
	}

	expectedCmd := self.expectedCmds[self.expectedCmdIndex]
	output, err := expectedCmd(cmdObj)

	self.expectedCmdIndex++

	return output, err
}

func (self *FakeCmdObjRunner) RunWithOutputs(cmdObj ICmdObj) (string, string, error) {
	output, err := self.RunWithOutput(cmdObj)
	return output, "", err
}

func (self *FakeCmdObjRunner) RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	output, err := self.RunWithOutput(cmdObj)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		stop, err := onLine(line)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}

	return nil
}

func (self *FakeCmdObjRunner) ExpectFunc(fn func(cmdObj ICmdObj) (string, error)) *FakeCmdObjRunner {
	self.expectedCmds = append(self.expectedCmds, fn)

	return self
}

func (self *FakeCmdObjRunner) Expect(expectedCmdStr string, output string, err error) *FakeCmdObjRunner {
	self.ExpectFunc(func(cmdObj ICmdObj) (string, error) {
		cmdStr := cmdObj.ToString()
		assert.Equal(self.t, expectedCmdStr, cmdStr, fmt.Sprintf("expected command %d to be %s, but was %s", self.expectedCmdIndex+1, expectedCmdStr, cmdStr))

		return output, err
	})

	return self
}

func (self *FakeCmdObjRunner) ExpectArgs(expectedArgs []string, output string, err error) *FakeCmdObjRunner {
	self.ExpectFunc(func(cmdObj ICmdObj) (string, error) {
		args := cmdObj.GetCmd().Args

		if runtime.GOOS == "windows" {
			// thanks to the secureexec package, the first arg is something like
			// '"C:\\Program Files\\Git\\mingw64\\bin\\<command>.exe"
			// on windows so we'll just ensure it contains our program
			assert.Contains(self.t, args[0], expectedArgs[0])
		} else {
			// first arg is the program name
			assert.Equal(self.t, expectedArgs[0], args[0])
		}

		assert.EqualValues(self.t, expectedArgs[1:], args[1:], fmt.Sprintf("command %d did not match expectation", self.expectedCmdIndex+1))

		return output, err
	})

	return self
}

func (self *FakeCmdObjRunner) ExpectGitArgs(expectedArgs []string, output string, err error) *FakeCmdObjRunner {
	self.ExpectFunc(func(cmdObj ICmdObj) (string, error) {
		// first arg is 'git' on unix and something like '"C:\\Program Files\\Git\\mingw64\\bin\\git.exe" on windows so we'll just ensure it ends in either 'git' or 'git.exe'
		re := regexp.MustCompile(`git(\.exe)?$`)
		args := cmdObj.GetCmd().Args
		if !re.MatchString(args[0]) {
			self.t.Errorf("expected first arg to end in .git or .git.exe but was %s", args[0])
		}
		assert.EqualValues(self.t, expectedArgs, args[1:], fmt.Sprintf("command %d did not match expectation", self.expectedCmdIndex+1))

		return output, err
	})

	return self
}

func (self *FakeCmdObjRunner) CheckForMissingCalls() {
	if self.expectedCmdIndex < len(self.expectedCmds) {
		self.t.Errorf("expected command %d to be called, but was not", self.expectedCmdIndex+1)
	}
}
