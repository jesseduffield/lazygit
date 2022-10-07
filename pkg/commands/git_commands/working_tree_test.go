package git_commands

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestWorkingTreeStageFile(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git add -- "test.txt"`, "", nil)

	instance := buildWorkingTreeCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.StageFile("test.txt"))
	runner.CheckForMissingCalls()
}

func TestWorkingTreeStageFiles(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git add -- "test.txt" "test2.txt"`, "", nil)

	instance := buildWorkingTreeCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.StageFiles([]string{"test.txt", "test2.txt"}))
	runner.CheckForMissingCalls()
}

func TestWorkingTreeUnstageFile(t *testing.T) {
	type scenario struct {
		testName string
		reset    bool
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "Remove an untracked file from staging",
			reset:    false,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git rm --cached --force -- "test.txt"`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "Remove a tracked file from staging",
			reset:    true,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git reset HEAD -- "test.txt"`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.UnStageFile([]string{"test.txt"}, s.reset))
		})
	}
}

// these tests don't cover everything, in part because we already have an integration
// test which does cover everything. I don't want to unnecessarily assert on the 'how'
// when the 'what' is what matters
func TestWorkingTreeDiscardAllFileChanges(t *testing.T) {
	type scenario struct {
		testName      string
		file          *models.File
		removeFile    func(string) error
		runner        *oscommands.FakeCmdObjRunner
		expectedError string
	}

	scenarios := []scenario{
		{
			testName: "An error occurred when resetting",
			file: &models.File{
				Name:             "test",
				HasStagedChanges: true,
			},
			removeFile: func(string) error { return nil },
			runner: oscommands.NewFakeRunner(t).
				Expect(`git reset -- "test"`, "", errors.New("error")),
			expectedError: "error",
		},
		{
			testName: "An error occurred when removing file",
			file: &models.File{
				Name:    "test",
				Tracked: false,
				Added:   true,
			},
			removeFile: func(string) error {
				return fmt.Errorf("an error occurred when removing file")
			},
			runner:        oscommands.NewFakeRunner(t),
			expectedError: "an error occurred when removing file",
		},
		{
			testName: "An error occurred with checkout",
			file: &models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			removeFile: func(string) error { return nil },
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout -- "test"`, "", errors.New("error")),
			expectedError: "error",
		},
		{
			testName: "Checkout only",
			file: &models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			removeFile: func(string) error { return nil },
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout -- "test"`, "", nil),
			expectedError: "",
		},
		{
			testName: "Reset and checkout staged changes",
			file: &models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: true,
			},
			removeFile: func(string) error { return nil },
			runner: oscommands.NewFakeRunner(t).
				Expect(`git reset -- "test"`, "", nil).
				Expect(`git checkout -- "test"`, "", nil),
			expectedError: "",
		},
		{
			testName: "Reset and checkout merge conflicts",
			file: &models.File{
				Name:              "test",
				Tracked:           true,
				HasMergeConflicts: true,
			},
			removeFile: func(string) error { return nil },
			runner: oscommands.NewFakeRunner(t).
				Expect(`git reset -- "test"`, "", nil).
				Expect(`git checkout -- "test"`, "", nil),
			expectedError: "",
		},
		{
			testName: "Reset and remove",
			file: &models.File{
				Name:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: true,
			},
			removeFile: func(filename string) error {
				assert.Equal(t, "test", filename)
				return nil
			},
			runner: oscommands.NewFakeRunner(t).
				Expect(`git reset -- "test"`, "", nil),
			expectedError: "",
		},
		{
			testName: "Remove only",
			file: &models.File{
				Name:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: false,
			},
			removeFile: func(filename string) error {
				assert.Equal(t, "test", filename)
				return nil
			},
			runner:        oscommands.NewFakeRunner(t),
			expectedError: "",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, removeFile: s.removeFile})
			err := instance.DiscardAllFileChanges(s.file)

			if s.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, s.expectedError, err.Error())
			}
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiff(t *testing.T) {
	type scenario struct {
		testName         string
		file             *models.File
		plain            bool
		cached           bool
		ignoreWhitespace bool
		contextSize      int
		runner           *oscommands.FakeCmdObjRunner
	}

	const expectedResult = "pretend this is an actual git diff"

	scenarios := []scenario{
		{
			testName: "Default case",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:            false,
			cached:           false,
			ignoreWhitespace: false,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --color=always -- "test.txt"`, expectedResult, nil),
		},
		{
			testName: "cached",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:            false,
			cached:           true,
			ignoreWhitespace: false,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --color=always --cached -- "test.txt"`, expectedResult, nil),
		},
		{
			testName: "plain",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:            true,
			cached:           false,
			ignoreWhitespace: false,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --color=never -- "test.txt"`, expectedResult, nil),
		},
		{
			testName: "File not tracked and file has no staged changes",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          false,
			},
			plain:            false,
			cached:           false,
			ignoreWhitespace: false,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --color=always --no-index -- /dev/null "test.txt"`, expectedResult, nil),
		},
		{
			testName: "Default case (ignore whitespace)",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:            false,
			cached:           false,
			ignoreWhitespace: true,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --color=always --ignore-all-space -- "test.txt"`, expectedResult, nil),
		},
		{
			testName: "Show diff with custom context size",
			file: &models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:            false,
			cached:           false,
			ignoreWhitespace: false,
			contextSize:      17,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=17 --color=always -- "test.txt"`, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.DiffContextSize = s.contextSize

			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, userConfig: userConfig})
			result := instance.WorktreeFileDiff(s.file, s.plain, s.cached, s.ignoreWhitespace)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeShowFileDiff(t *testing.T) {
	type scenario struct {
		testName    string
		from        string
		to          string
		reverse     bool
		plain       bool
		contextSize int
		runner      *oscommands.FakeCmdObjRunner
	}

	const expectedResult = "pretend this is an actual git diff"

	scenarios := []scenario{
		{
			testName:    "Default case",
			from:        "1234567890",
			to:          "0987654321",
			reverse:     false,
			plain:       false,
			contextSize: 3,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=3 --no-renames --color=always 1234567890 0987654321 -- "test.txt"`, expectedResult, nil),
		},
		{
			testName:    "Show diff with custom context size",
			from:        "1234567890",
			to:          "0987654321",
			reverse:     false,
			plain:       false,
			contextSize: 123,
			runner: oscommands.NewFakeRunner(t).
				Expect(`git diff --submodule --no-ext-diff --unified=123 --no-renames --color=always 1234567890 0987654321 -- "test.txt"`, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.DiffContextSize = s.contextSize

			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, userConfig: userConfig})

			result, err := instance.ShowFileDiff(s.from, s.to, s.reverse, "test.txt", s.plain)
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeCheckoutFile(t *testing.T) {
	type scenario struct {
		testName  string
		commitSha string
		fileName  string
		runner    *oscommands.FakeCmdObjRunner
		test      func(error)
	}

	scenarios := []scenario{
		{
			testName:  "typical case",
			commitSha: "11af912",
			fileName:  "test999.txt",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout 11af912 -- "test999.txt"`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName:  "returns error if there is one",
			commitSha: "11af912",
			fileName:  "test999.txt",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout 11af912 -- "test999.txt"`, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})

			s.test(instance.CheckoutFile(s.commitSha, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}

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
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.ApplyPatch("test", "cached"))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiscardUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		file     *models.File
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			file:     &models.File{Name: "test.txt"},
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout -- "test.txt"`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.DiscardUnstagedFileChanges(s.file))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiscardAnyUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git checkout -- .`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.DiscardAnyUnstagedFileChanges())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeRemoveUntrackedFiles(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git clean -fd`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.RemoveUntrackedFiles())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeResetHard(t *testing.T) {
	type scenario struct {
		testName string
		ref      string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			"HEAD",
			oscommands.NewFakeRunner(t).
				Expect(`git reset --hard "HEAD"`, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.ResetHard(s.ref))
		})
	}
}
