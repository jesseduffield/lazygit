package oscommands

import (
	"bufio"
	"fmt"
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

func NewFakeRunner(t *testing.T) *FakeCmdObjRunner {
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
	self.expectedCmdIndex++

	return expectedCmd(cmdObj)
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
		if cmdStr != expectedCmdStr {
			assert.Equal(self.t, expectedCmdStr, cmdStr, fmt.Sprintf("expected command %d to be %s, but was %s", self.expectedCmdIndex+1, expectedCmdStr, cmdStr))
			return "", errors.New("expected cmd")
		}

		return output, err
	})

	return self
}
