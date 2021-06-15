package commands

// import (
// 	"os/exec"
// 	"testing"

// 	"github.com/jesseduffield/lazygit/pkg/test"
// 	"github.com/stretchr/testify/assert"
// )

// // TestGitCommandRebaseBranch is a function.
// func TestGitCommandRebaseBranch(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		arg      string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"successful rebase",
// 			"master",
// 			test.CreateMockCommand(t, []*test.CommandSwapper{
// 				{
// 					Expect:  "git rebase --interactive --autostash --keep-empty master",
// 					Replace: "echo",
// 				},
// 			}),
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 		{
// 			"unsuccessful rebase",
// 			"master",
// 			test.CreateMockCommand(t, []*test.CommandSwapper{
// 				{
// 					Expect:  "git rebase --interactive --autostash --keep-empty master",
// 					Replace: "test",
// 				},
// 			}),
// 			func(err error) {
// 				assert.Error(t, err)
// 			},
// 		},
// 	}

// 	gitCmd := NewDummyGit()

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.RebaseBranch(s.arg))
// 		})
// 	}
// }

// // TestGitCommandResetToCommit is a function.
// func TestGitCommandResetToRef(t *testing.T) {
// 	gitCmd := NewDummyGit()
// 	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
// 		assert.EqualValues(t, "git", cmd)
// 		assert.EqualValues(t, []string{"reset", "--hard", "78976bc"}, args)

// 		return secureexec.Command("echo")
// 	}

// 	assert.NoError(t, gitCmd.ResetToRef("78976bc", "hard", oscommands.RunCommandOptions{}))
// }
