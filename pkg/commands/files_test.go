package commands_test

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os/exec"
// 	"runtime"
// 	"testing"

// 	"github.com/jesseduffield/lazygit/pkg/commands/models"
// 	"github.com/jesseduffield/lazygit/pkg/secureexec"
// 	"github.com/jesseduffield/lazygit/pkg/test"
// 	"github.com/stretchr/testify/assert"
// )

// // TestGitCommandDiff is a function.
// func TestGitCommandDiff(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		file     *models.File
// 		plain    bool
// 		cached   bool
// 	}

// 	scenarios := []scenario{
// 		{
// 			"Default case",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--", "test.txt"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			&models.File{
// 				Name:             "test.txt",
// 				HasStagedChanges: false,
// 				Tracked:          true,
// 			},
// 			false,
// 			false,
// 		},
// 		{
// 			"cached",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--cached", "--", "test.txt"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			&models.File{
// 				Name:             "test.txt",
// 				HasStagedChanges: false,
// 				Tracked:          true,
// 			},
// 			false,
// 			true,
// 		},
// 		{
// 			"plain",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=never", "--", "test.txt"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			&models.File{
// 				Name:             "test.txt",
// 				HasStagedChanges: false,
// 				Tracked:          true,
// 			},
// 			true,
// 			false,
// 		},
// 		{
// 			"File not tracked and file has no staged changes",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"diff", "--submodule", "--no-ext-diff", "--color=always", "--no-index", "--", "/dev/null", "test.txt"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			&models.File{
// 				Name:             "test.txt",
// 				HasStagedChanges: false,
// 				Tracked:          false,
// 			},
// 			false,
// 			false,
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			gitCmd.WorktreeFileDiff(s.file, s.plain, s.cached)
// 		})
// 	}
// }

// func TestGitCommandApplyPatch(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"valid case",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.Equal(t, "git", cmd)
// 				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
// 				filename := args[2]
// 				content, err := ioutil.ReadFile(filename)
// 				assert.NoError(t, err)

// 				assert.Equal(t, "test", string(content))

// 				return secureexec.Command("echo", "done")
// 			},
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 		{
// 			"command returns error",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.Equal(t, "git", cmd)
// 				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
// 				filename := args[2]
// 				// TODO: Ideally we want to mock out OSCommand here so that we're not
// 				// double handling testing it's CreateTempFile functionality,
// 				// but it is going to take a bit of work to make a proper mock for it
// 				// so I'm leaving it for another PR
// 				content, err := ioutil.ReadFile(filename)
// 				assert.NoError(t, err)

// 				assert.Equal(t, "test", string(content))

// 				return secureexec.Command("test")
// 			},
// 			func(err error) {
// 				assert.Error(t, err)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.ApplyPatch("test", "cached"))
// 		})
// 	}
// }

// // TestGitCommandDiscardOldFileChanges is a function.
// func TestGitCommandDiscardOldFileChanges(t *testing.T) {
// 	type scenario struct {
// 		testName          string
// 		getGitConfigValue func(string) (string, error)
// 		commits           []*models.Commit
// 		commitIndex       int
// 		fileName          string
// 		command           func(string, ...string) *exec.Cmd
// 		test              func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"returns error when index outside of range of commits",
// 			func(string) (string, error) {
// 				return "", nil
// 			},
// 			[]*models.Commit{},
// 			0,
// 			"test999.txt",
// 			nil,
// 			func(err error) {
// 				assert.Error(t, err)
// 			},
// 		},
// 		{
// 			"returns error when using gpg",
// 			func(string) (string, error) {
// 				return "true", nil
// 			},
// 			[]*models.Commit{{Name: "commit", Sha: "123456"}},
// 			0,
// 			"test999.txt",
// 			nil,
// 			func(err error) {
// 				assert.Error(t, err)
// 			},
// 		},
// 		{
// 			"checks out file if it already existed",
// 			func(string) (string, error) {
// 				return "", nil
// 			},
// 			[]*models.Commit{
// 				{Name: "commit", Sha: "123456"},
// 				{Name: "commit2", Sha: "abcdef"},
// 			},
// 			0,
// 			"test999.txt",
// 			test.CreateMockCommand(t, []*test.CommandSwapper{
// 				{
// 					Expect:  "git rebase --interactive --autostash --keep-empty abcdef",
// 					Replace: "echo",
// 				},
// 				{
// 					Expect:  "git cat-file -e HEAD^:test999.txt",
// 					Replace: "echo",
// 				},
// 				{
// 					Expect:  "git checkout HEAD^ test999.txt",
// 					Replace: "echo",
// 				},
// 				{
// 					Expect:  "git commit --amend --no-edit --allow-empty",
// 					Replace: "echo",
// 				},
// 				{
// 					Expect:  "git rebase --continue",
// 					Replace: "echo",
// 				},
// 			}),
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 		// test for when the file was created within the commit requires a refactor to support proper mocks
// 		// currently we'd need to mock out the os.Remove function and that's gonna introduce tech debt
// 	}

// 	gitCmd := NewDummyGit()

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd.GetOSCommand().Command = s.command
// 			gitCmd.getGitConfigValue = s.getGitConfigValue
// 			s.test(gitCmd.DiscardOldFileChanges(s.commits, s.commitIndex, s.fileName))
// 		})
// 	}
// }
