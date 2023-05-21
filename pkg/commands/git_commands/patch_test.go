package git_commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestPatchApplyPatch(t *testing.T) {
	type scenario struct {
		testName string
		opts     ApplyPatchOpts
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	// expectedArgs excludes the last argument which is an indeterminate filename
	expectFn := func(expectedArgs []string, errToReturn error) func(cmdObj oscommands.ICmdObj) (string, error) {
		return func(cmdObj oscommands.ICmdObj) (string, error) {
			args := cmdObj.Args()

			assert.Equal(t, len(args), len(expectedArgs)+1, fmt.Sprintf("unexpected command: %s", cmdObj.ToString()))

			filename := args[len(args)-1]

			content, err := os.ReadFile(filename)
			assert.NoError(t, err)

			assert.Equal(t, "test", string(content))

			return "", errToReturn
		}
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			opts:     ApplyPatchOpts{Cached: true},
			runner: oscommands.NewFakeRunner(t).
				ExpectFunc(expectFn([]string{"git", "apply", "--cached"}, nil)),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "command returns error",
			opts:     ApplyPatchOpts{Cached: true},
			runner: oscommands.NewFakeRunner(t).
				ExpectFunc(expectFn([]string{"git", "apply", "--cached"}, errors.New("error"))),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildPatchCommands(commonDeps{runner: s.runner})
			s.test(instance.ApplyPatch("test", s.opts))
			s.runner.CheckForMissingCalls()
		})
	}
}
