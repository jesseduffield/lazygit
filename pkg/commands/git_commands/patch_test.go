package git_commands

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestWorkingTreeApplyPatch(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	expectFn := func(regexStr string, errToReturn error) func(cmdObj oscommands.ICmdObj) (string, error) {
		return func(cmdObj oscommands.ICmdObj) (string, error) {
			re := regexp.MustCompile(regexStr)
			cmdStr := cmdObj.ToString()
			matches := re.FindStringSubmatch(cmdStr)
			assert.Equal(t, 2, len(matches), fmt.Sprintf("unexpected command: %s", cmdStr))

			filename := matches[1]

			content, err := os.ReadFile(filename)
			assert.NoError(t, err)

			assert.Equal(t, "test", string(content))

			return "", errToReturn
		}
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			runner: oscommands.NewFakeRunner(t).
				ExpectFunc(expectFn(`git apply --cached "(.*)"`, nil)),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "command returns error",
			runner: oscommands.NewFakeRunner(t).
				ExpectFunc(expectFn(`git apply --cached "(.*)"`, errors.New("error"))),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildPatchCommands(commonDeps{runner: s.runner})
			s.test(instance.ApplyPatch("test", ApplyPatchOpts{Cached: true}))
			s.runner.CheckForMissingCalls()
		})
	}
}
