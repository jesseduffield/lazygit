package oscommands

import (
	"bufio"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/go-errors/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// for use in testing

type FakeCmdObjRunner struct {
	t *testing.T
	// commands can be run in any order; mimicking the concurrent behaviour of
	// production code.
	expectedCmds []CmdObjMatcher

	invokedCmdIndexes []int

	mutex sync.Mutex
}

type CmdObjMatcher struct {
	description string
	// returns true if the matcher matches the command object
	test func(ICmdObj) bool

	// output of the command
	output string
	// error of the command
	err error
}

var _ ICmdObjRunner = &FakeCmdObjRunner{}

func NewFakeRunner(t *testing.T) *FakeCmdObjRunner { //nolint:thelper
	return &FakeCmdObjRunner{t: t}
}

func (self *FakeCmdObjRunner) remainingExpectedCmds() []CmdObjMatcher {
	return lo.Filter(self.expectedCmds, func(_ CmdObjMatcher, i int) bool {
		return !lo.Contains(self.invokedCmdIndexes, i)
	})
}

func (self *FakeCmdObjRunner) Run(cmdObj ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *FakeCmdObjRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if len(self.remainingExpectedCmds()) == 0 {
		self.t.Errorf("ran too many commands. Unexpected command: `%s`", cmdObj.ToString())
		return "", errors.New("ran too many commands")
	}

	for i := range self.expectedCmds {
		if lo.Contains(self.invokedCmdIndexes, i) {
			continue
		}
		expectedCmd := self.expectedCmds[i]
		matched := expectedCmd.test(cmdObj)
		if matched {
			self.invokedCmdIndexes = append(self.invokedCmdIndexes, i)
			return expectedCmd.output, expectedCmd.err
		}
	}

	self.t.Errorf("Unexpected command: `%s`", cmdObj.ToString())
	return "", nil
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

func (self *FakeCmdObjRunner) ExpectFunc(description string, fn func(cmdObj ICmdObj) bool, output string, err error) *FakeCmdObjRunner {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.expectedCmds = append(self.expectedCmds, CmdObjMatcher{
		test:        fn,
		output:      output,
		err:         err,
		description: description,
	})

	return self
}

func (self *FakeCmdObjRunner) ExpectArgs(expectedArgs []string, output string, err error) *FakeCmdObjRunner {
	description := fmt.Sprintf("matches args %s", strings.Join(expectedArgs, " "))
	self.ExpectFunc(description, func(cmdObj ICmdObj) bool {
		args := cmdObj.GetCmd().Args

		if runtime.GOOS == "windows" {
			// thanks to the secureexec package, the first arg is something like
			// '"C:\\Program Files\\Git\\mingw64\\bin\\<command>.exe"
			// on windows so we'll just ensure it contains our program
			if !strings.Contains(args[0], expectedArgs[0]) {
				return false
			}
		} else {
			// first arg is the program name
			if expectedArgs[0] != args[0] {
				return false
			}
		}

		if !slices.Equal(expectedArgs[1:], args[1:]) {
			return false
		}

		return true
	}, output, err)

	return self
}

func (self *FakeCmdObjRunner) ExpectGitArgs(expectedArgs []string, output string, err error) *FakeCmdObjRunner {
	description := fmt.Sprintf("matches git args %s", strings.Join(expectedArgs, " "))
	self.ExpectFunc(description, func(cmdObj ICmdObj) bool {
		// first arg is 'git' on unix and something like '"C:\\Program Files\\Git\\mingw64\\bin\\git.exe" on windows so we'll just ensure it ends in either 'git' or 'git.exe'
		re := regexp.MustCompile(`git(\.exe)?$`)
		args := cmdObj.GetCmd().Args
		if !re.MatchString(args[0]) {
			return false
		}

		if !slices.Equal(expectedArgs, args[1:]) {
			return false
		}

		return true
	}, output, err)

	return self
}

func (self *FakeCmdObjRunner) CheckForMissingCalls() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	remaining := self.remainingExpectedCmds()
	if len(remaining) > 0 {
		self.t.Errorf(
			"expected %d more command(s) to be run. Remaining commands:\n%s",
			len(remaining),
			strings.Join(
				lo.Map(remaining, func(cmdObj CmdObjMatcher, _ int) string {
					return cmdObj.description
				}),
				"\n",
			),
		)
	}
}
