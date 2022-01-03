package commands

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandStageFile(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"add", "--", "test.txt"}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.StageFile("test.txt"))
	runner.CheckForMissingCalls()
}

func TestGitCommandUnstageFile(t *testing.T) {
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
				ExpectGitArgs([]string{"rm", "--cached", "--force", "--", "test.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "Remove a tracked file from staging",
			reset:    true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "HEAD", "--", "test.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.UnStageFile([]string{"test.txt"}, s.reset))
		})
	}
}

// these tests don't cover everything, in part because we already have an integration
// test which does cover everything. I don't want to unnecessarily assert on the 'how'
// when the 'what' is what matters
func TestGitCommandDiscardAllFileChanges(t *testing.T) {
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
				ExpectGitArgs([]string{"reset", "--", "test"}, "", errors.New("error")),
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
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", errors.New("error")),
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
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
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
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
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
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
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
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil),
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			gitCmd.OSCommand.SetRemoveFile(s.removeFile)
			err := gitCmd.DiscardAllFileChanges(s.file)

			if s.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, s.expectedError, err.Error())
			}
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandDiff(t *testing.T) {
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--color=always", "--", "test.txt"}, expectedResult, nil),
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--color=always", "--cached", "--", "test.txt"}, expectedResult, nil),
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--color=never", "--", "test.txt"}, expectedResult, nil),
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--color=always", "--no-index", "--", "/dev/null", "test.txt"}, expectedResult, nil),
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--color=always", "--ignore-all-space", "--", "test.txt"}, expectedResult, nil),
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=17", "--color=always", "--", "test.txt"}, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			gitCmd.UserConfig.Git.DiffContextSize = s.contextSize
			result := gitCmd.WorktreeFileDiff(s.file, s.plain, s.cached, s.ignoreWhitespace)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandShowFileDiff(t *testing.T) {
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
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=3", "--no-renames", "--color=always", "1234567890", "0987654321", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName:    "Show diff with custom context size",
			from:        "1234567890",
			to:          "0987654321",
			reverse:     false,
			plain:       false,
			contextSize: 123,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"diff", "--submodule", "--no-ext-diff", "--unified=123", "--no-renames", "--color=always", "1234567890", "0987654321", "--", "test.txt"}, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			gitCmd.UserConfig.Git.DiffContextSize = s.contextSize
			result, err := gitCmd.ShowFileDiff(s.from, s.to, s.reverse, "test.txt", s.plain)
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandCheckoutFile(t *testing.T) {
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.CheckoutFile(s.commitSha, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandApplyPatch(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	expectFn := func(regexStr string, errToReturn error) func(cmdObj oscommands.ICmdObj) (string, error) {
		return func(cmdObj oscommands.ICmdObj) (string, error) {
			re := regexp.MustCompile(regexStr)
			matches := re.FindStringSubmatch(cmdObj.ToString())
			assert.Equal(t, 2, len(matches))

			filename := matches[1]

			content, err := ioutil.ReadFile(filename)
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.ApplyPatch("test", "cached"))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandDiscardOldFileChanges(t *testing.T) {
	type scenario struct {
		testName               string
		gitConfigMockResponses map[string]string
		commits                []*models.Commit
		commitIndex            int
		fileName               string
		runner                 *oscommands.FakeCmdObjRunner
		test                   func(error)
	}

	scenarios := []scenario{
		{
			testName:               "returns error when index outside of range of commits",
			gitConfigMockResponses: nil,
			commits:                []*models.Commit{},
			commitIndex:            0,
			fileName:               "test999.txt",
			runner:                 oscommands.NewFakeRunner(t),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:               "returns error when using gpg",
			gitConfigMockResponses: map[string]string{"commit.gpgsign": "true"},
			commits:                []*models.Commit{{Name: "commit", Sha: "123456"}},
			commitIndex:            0,
			fileName:               "test999.txt",
			runner:                 oscommands.NewFakeRunner(t),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testName:               "checks out file if it already existed",
			gitConfigMockResponses: nil,
			commits: []*models.Commit{
				{Name: "commit", Sha: "123456"},
				{Name: "commit2", Sha: "abcdef"},
			},
			commitIndex: 0,
			fileName:    "test999.txt",
			runner: oscommands.NewFakeRunner(t).
				Expect(`git rebase --interactive --autostash --keep-empty abcdef`, "", nil).
				Expect(`git cat-file -e HEAD^:"test999.txt"`, "", nil).
				Expect(`git checkout HEAD^ -- "test999.txt"`, "", nil).
				Expect(`git commit --amend --no-edit --allow-empty`, "", nil).
				Expect(`git rebase --continue`, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		// test for when the file was created within the commit requires a refactor to support proper mocks
		// currently we'd need to mock out the os.Remove function and that's gonna introduce tech debt
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			gitCmd.GitConfig = git_config.NewFakeGitConfig(s.gitConfigMockResponses)
			s.test(gitCmd.DiscardOldFileChanges(s.commits, s.commitIndex, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandDiscardUnstagedFileChanges(t *testing.T) {
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.DiscardUnstagedFileChanges(s.file))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandDiscardAnyUnstagedFileChanges(t *testing.T) {
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.DiscardAnyUnstagedFileChanges())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandRemoveUntrackedFiles(t *testing.T) {
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.RemoveUntrackedFiles())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestEditFileCmdStr(t *testing.T) {
	type scenario struct {
		filename                  string
		configEditCommand         string
		configEditCommandTemplate string
		runner                    *oscommands.FakeCmdObjRunner
		getenv                    func(string) string
		gitConfigMockResponses    map[string]string
		test                      func(string, error)
	}

	scenarios := []scenario{
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "", errors.New("error")),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.EqualError(t, err, "No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "nano",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `nano "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: map[string]string{"core.editor": "nano"},
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `nano "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				if env == "VISUAL" {
					return "nano"
				}

				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				if env == "EDITOR" {
					return "emacs"
				}

				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `emacs "test"`, cmdStr)
			},
		},
		{
			filename:                  "test",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "/usr/bin/vi", nil),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vi "test"`, cmdStr)
			},
		},
		{
			filename:                  "file/with space",
			configEditCommand:         "",
			configEditCommandTemplate: "{{editor}} {{filename}}",
			runner: oscommands.NewFakeRunner(t).
				Expect(`which vi`, "/usr/bin/vi", nil),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vi "file/with space"`, cmdStr)
			},
		},
		{
			filename:                  "open file/at line",
			configEditCommand:         "vim",
			configEditCommandTemplate: "{{editor}} +{{line}} {{filename}}",
			runner:                    oscommands.NewFakeRunner(t),
			getenv: func(env string) string {
				return ""
			},
			gitConfigMockResponses: nil,
			test: func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, `vim +1 "open file/at line"`, cmdStr)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := NewDummyGitCommandWithRunner(s.runner)
		gitCmd.UserConfig.OS.EditCommand = s.configEditCommand
		gitCmd.UserConfig.OS.EditCommandTemplate = s.configEditCommandTemplate
		gitCmd.OSCommand.Getenv = s.getenv
		gitCmd.GitConfig = git_config.NewFakeGitConfig(s.gitConfigMockResponses)
		s.test(gitCmd.EditFileCmdStr(s.filename, 1))
		s.runner.CheckForMissingCalls()
	}
}
