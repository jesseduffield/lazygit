package commands

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandStageFile is a function.
func TestGitCommandStageFile(t *testing.T) {
	gitCmd := NewDummyGitCommand()
	gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"add", "--", "test.txt"}, args)

		return secureexec.Command("echo")
	}

	assert.NoError(t, gitCmd.StageFile("test.txt"))
}

// TestGitCommandUnstageFile is a function.
func TestGitCommandUnstageFile(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
		reset    bool
	}

	scenarios := []scenario{
		{
			"Remove an untracked file from staging",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"rm", "--cached", "--force", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
			false,
		},
		{
			"Remove a tracked file from staging",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"reset", "HEAD", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.UnStageFile([]string{"test.txt"}, s.reset))
		})
	}
}

// TestGitCommandDiscardAllFileChanges is a function.
// these tests don't cover everything, in part because we already have an integration
// test which does cover everything. I don't want to unnecessarily assert on the 'how'
// when the 'what' is what matters
func TestGitCommandDiscardAllFileChanges(t *testing.T) {
	type scenario struct {
		testName   string
		command    func() (func(string, ...string) *exec.Cmd, *[][]string)
		test       func(*[][]string, error)
		file       *models.File
		removeFile func(string) error
	}

	scenarios := []scenario{
		{
			"An error occurred when resetting",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("test")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.Error(t, err)
				assert.Len(t, *cmdsCalled, 1)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"reset", "--", "test"},
				})
			},
			&models.File{
				Name:             "test",
				HasStagedChanges: true,
			},
			func(string) error {
				return nil
			},
		},
		{
			"An error occurred when removing file",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("test")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "an error occurred when removing file")
				assert.Len(t, *cmdsCalled, 0)
			},
			&models.File{
				Name:    "test",
				Tracked: false,
				Added:   true,
			},
			func(string) error {
				return fmt.Errorf("an error occurred when removing file")
			},
		},
		{
			"An error occurred with checkout",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("test")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.Error(t, err)
				assert.Len(t, *cmdsCalled, 1)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"checkout", "--", "test"},
				})
			},
			&models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			func(string) error {
				return nil
			},
		},
		{
			"Checkout only",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("echo")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cmdsCalled, 1)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"checkout", "--", "test"},
				})
			},
			&models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			func(string) error {
				return nil
			},
		},
		{
			"Reset and checkout staged changes",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("echo")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cmdsCalled, 2)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"reset", "--", "test"},
					{"checkout", "--", "test"},
				})
			},
			&models.File{
				Name:             "test",
				Tracked:          true,
				HasStagedChanges: true,
			},
			func(string) error {
				return nil
			},
		},
		{
			"Reset and checkout merge conflicts",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("echo")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cmdsCalled, 2)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"reset", "--", "test"},
					{"checkout", "--", "test"},
				})
			},
			&models.File{
				Name:              "test",
				Tracked:           true,
				HasMergeConflicts: true,
			},
			func(string) error {
				return nil
			},
		},
		{
			"Reset and remove",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("echo")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cmdsCalled, 1)
				assert.EqualValues(t, *cmdsCalled, [][]string{
					{"reset", "--", "test"},
				})
			},
			&models.File{
				Name:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: true,
			},
			func(filename string) error {
				assert.Equal(t, "test", filename)
				return nil
			},
		},
		{
			"Remove only",
			func() (func(string, ...string) *exec.Cmd, *[][]string) {
				cmdsCalled := [][]string{}
				return func(cmd string, args ...string) *exec.Cmd {
					cmdsCalled = append(cmdsCalled, args)

					return secureexec.Command("echo")
				}, &cmdsCalled
			},
			func(cmdsCalled *[][]string, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cmdsCalled, 0)
			},
			&models.File{
				Name:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: false,
			},
			func(filename string) error {
				assert.Equal(t, "test", filename)
				return nil
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			var cmdsCalled *[][]string
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command, cmdsCalled = s.command()
			gitCmd.OSCommand.SetRemoveFile(s.removeFile)
			s.test(cmdsCalled, gitCmd.DiscardAllFileChanges(s.file))
		})
	}
}

// TestGitCommandDiff is a function.
func TestGitCommandDiff(t *testing.T) {
	type scenario struct {
		testName         string
		command          func(string, ...string) *exec.Cmd
		file             *models.File
		plain            bool
		cached           bool
		ignoreWhitespace bool
	}

	scenarios := []scenario{
		{
			"Default case",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			&models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			false,
			false,
			false,
		},
		{
			"cached",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--cached", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			&models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			false,
			true,
			false,
		},
		{
			"plain",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=never", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			&models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			true,
			false,
			false,
		},
		{
			"File not tracked and file has no staged changes",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--no-index", "--", "/dev/null", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			&models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          false,
			},
			false,
			false,
			false,
		},
		{
			"Default case (ignore whitespace)",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--ignore-all-space", "--", "test.txt"}, args)

				return secureexec.Command("echo")
			},
			&models.File{
				Name:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			false,
			false,
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			gitCmd.WorktreeFileDiff(s.file, s.plain, s.cached, s.ignoreWhitespace)
		})
	}
}

// TestGitCommandCheckoutFile is a function.
func TestGitCommandCheckoutFile(t *testing.T) {
	type scenario struct {
		testName  string
		commitSha string
		fileName  string
		command   func(string, ...string) *exec.Cmd
		test      func(error)
	}

	scenarios := []scenario{
		{
			"typical case",
			"11af912",
			"test999.txt",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git checkout 11af912 -- test999.txt",
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"returns error if there is one",
			"11af912",
			"test999.txt",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git checkout 11af912 -- test999.txt",
					Replace: "test",
				},
			}),
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.CheckoutFile(s.commitSha, s.fileName))
		})
	}
}

func TestGitCommandApplyPatch(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			func(cmd string, args ...string) *exec.Cmd {
				assert.Equal(t, "git", cmd)
				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
				filename := args[2]
				content, err := ioutil.ReadFile(filename)
				assert.NoError(t, err)

				assert.Equal(t, "test", string(content))

				return secureexec.Command("echo", "done")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"command returns error",
			func(cmd string, args ...string) *exec.Cmd {
				assert.Equal(t, "git", cmd)
				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
				filename := args[2]
				// TODO: Ideally we want to mock out OSCommand here so that we're not
				// double handling testing it's CreateTempFile functionality,
				// but it is going to take a bit of work to make a proper mock for it
				// so I'm leaving it for another PR
				content, err := ioutil.ReadFile(filename)
				assert.NoError(t, err)

				assert.Equal(t, "test", string(content))

				return secureexec.Command("test")
			},
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.ApplyPatch("test", "cached"))
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
		command                func(string, ...string) *exec.Cmd
		test                   func(error)
	}

	scenarios := []scenario{
		{
			"returns error when index outside of range of commits",
			nil,
			[]*models.Commit{},
			0,
			"test999.txt",
			nil,
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"returns error when using gpg",
			map[string]string{"commit.gpgsign": "true"},
			[]*models.Commit{{Name: "commit", Sha: "123456"}},
			0,
			"test999.txt",
			nil,
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"checks out file if it already existed",
			nil,
			[]*models.Commit{
				{Name: "commit", Sha: "123456"},
				{Name: "commit2", Sha: "abcdef"},
			},
			0,
			"test999.txt",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git rebase --interactive --autostash --keep-empty abcdef",
					Replace: "echo",
				},
				{
					Expect:  "git cat-file -e HEAD^:test999.txt",
					Replace: "echo",
				},
				{
					Expect:  "git checkout HEAD^ -- test999.txt",
					Replace: "echo",
				},
				{
					Expect:  "git commit --amend --no-edit --allow-empty",
					Replace: "echo",
				},
				{
					Expect:  "git rebase --continue",
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
		// test for when the file was created within the commit requires a refactor to support proper mocks
		// currently we'd need to mock out the os.Remove function and that's gonna introduce tech debt
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			gitCmd.GitConfig = git_config.NewFakeGitConfig(s.gitConfigMockResponses)
			s.test(gitCmd.DiscardOldFileChanges(s.commits, s.commitIndex, s.fileName))
		})
	}
}

// TestGitCommandDiscardUnstagedFileChanges is a function.
func TestGitCommandDiscardUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		file     *models.File
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			&models.File{Name: "test.txt"},
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git checkout -- "test.txt"`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.DiscardUnstagedFileChanges(s.file))
		})
	}
}

// TestGitCommandDiscardAnyUnstagedFileChanges is a function.
func TestGitCommandDiscardAnyUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git checkout -- .`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.DiscardAnyUnstagedFileChanges())
		})
	}
}

// TestGitCommandRemoveUntrackedFiles is a function.
func TestGitCommandRemoveUntrackedFiles(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git clean -fd`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.RemoveUntrackedFiles())
		})
	}
}

// TestEditFileCmdStr is a function.
func TestEditFileCmdStr(t *testing.T) {
	gitCmd := NewDummyGitCommand()

	type scenario struct {
		filename                  string
		configEditCommand         string
		configEditCommandTemplate string
		command                   func(string, ...string) *exec.Cmd
		getenv                    func(string) string
		gitConfigMockResponses    map[string]string
		test                      func(string, error)
	}

	scenarios := []scenario{
		{
			"test",
			"",
			"{{editor}} {{filename}}",
			func(name string, arg ...string) *exec.Cmd {
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.EqualError(t, err, "No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
			},
		},
		{
			"test",
			"nano",
			"{{editor}} {{filename}}",
			func(name string, args ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "nano "+gitCmd.OSCommand.Quote("test"), cmdStr)
			},
		},
		{
			"test",
			"",
			"{{editor}} {{filename}}",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				return ""
			},
			map[string]string{"core.editor": "nano"},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "nano "+gitCmd.OSCommand.Quote("test"), cmdStr)
			},
		},
		{
			"test",
			"",
			"{{editor}} {{filename}}",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				if env == "VISUAL" {
					return "nano"
				}

				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"test",
			"",
			"{{editor}} {{filename}}",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				if env == "EDITOR" {
					return "emacs"
				}

				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "emacs "+gitCmd.OSCommand.Quote("test"), cmdStr)
			},
		},
		{
			"test",
			"",
			"{{editor}} {{filename}}",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vi "+gitCmd.OSCommand.Quote("test"), cmdStr)
			},
		},
		{
			"file/with space",
			"",
			"{{editor}} {{filename}}",
			func(name string, args ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vi "+gitCmd.OSCommand.Quote("file/with space"), cmdStr)
			},
		},
		{
			"open file/at line",
			"vim",
			"{{editor}} +{{line}} {{filename}}",
			func(name string, args ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			nil,
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vim +1 "+gitCmd.OSCommand.Quote("open file/at line"), cmdStr)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd.Config.GetUserConfig().OS.EditCommand = s.configEditCommand
		gitCmd.Config.GetUserConfig().OS.EditCommandTemplate = s.configEditCommandTemplate
		gitCmd.OSCommand.Command = s.command
		gitCmd.OSCommand.Getenv = s.getenv
		gitCmd.GitConfig = git_config.NewFakeGitConfig(s.gitConfigMockResponses)
		s.test(gitCmd.EditFileCmdStr(s.filename, 1))
	}
}
