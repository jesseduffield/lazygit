package commands

// import (
// 	"os/exec"
// 	"testing"

// 	"github.com/jesseduffield/lazygit/pkg/secureexec"
// 	"github.com/jesseduffield/lazygit/pkg/test"
// 	"github.com/stretchr/testify/assert"
// )

// // TestGitCommandGetCommitDifferences is a function.
// func TestGitCommandGetCommitDifferences(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(string, string)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"Can't retrieve pushable count",
// 			func(string, ...string) *exec.Cmd {
// 				return secureexec.Command("test")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "?", pushableCount)
// 				assert.EqualValues(t, "?", pullableCount)
// 			},
// 		},
// 		{
// 			"Can't retrieve pullable count",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				if args[1] == "HEAD..@{u}" {
// 					return secureexec.Command("test")
// 				}

// 				return secureexec.Command("echo")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "?", pushableCount)
// 				assert.EqualValues(t, "?", pullableCount)
// 			},
// 		},
// 		{
// 			"Retrieve pullable and pushable count",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				if args[1] == "HEAD..@{u}" {
// 					return secureexec.Command("echo", "10")
// 				}

// 				return secureexec.Command("echo", "11")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "11", pushableCount)
// 				assert.EqualValues(t, "10", pullableCount)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.GetCommitDifferences("HEAD", "@{u}"))
// 		})
// 	}
// }

// // TestGitCommandNewBranch is a function.
// func TestGitCommandNewBranch(t *testing.T) {
// 	gitCmd := NewDummyGit()
// 	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
// 		assert.EqualValues(t, "git", cmd)
// 		assert.EqualValues(t, []string{"checkout", "-b", "test", "master"}, args)

// 		return secureexec.Command("echo")
// 	}

// 	assert.NoError(t, gitCmd.NewBranch("test", "master"))
// }

// // TestGitCommandDeleteBranch is a function.
// func TestGitCommandDeleteBranch(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		branch   string
// 		force    bool
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"Delete a branch",
// 			"test",
// 			false,
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"branch", "-d", "test"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 		{
// 			"Force delete a branch",
// 			"test",
// 			true,
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"branch", "-D", "test"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.DeleteBranch(s.branch, s.force))
// 		})
// 	}
// }

// // TestGitCommandMerge is a function.
// func TestGitCommandMerge(t *testing.T) {
// 	gitCmd := NewDummyGit()
// 	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
// 		assert.EqualValues(t, "git", cmd)
// 		assert.EqualValues(t, []string{"merge", "--no-edit", "test"}, args)

// 		return secureexec.Command("echo")
// 	}

// 	assert.NoError(t, gitCmd.Merge("test", MergeOpts{}))
// }

// // TestGitCommandCheckout is a function.
// func TestGitCommandCheckout(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 		force    bool
// 	}

// 	scenarios := []scenario{
// 		{
// 			"Checkout",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"checkout", "test"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 			false,
// 		},
// 		{
// 			"Checkout forced",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)
// 				assert.EqualValues(t, []string{"checkout", "--force", "test"}, args)

// 				return secureexec.Command("echo")
// 			},
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 			true,
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.Checkout("test", CheckoutOptions{Force: s.force}))
// 		})
// 	}
// }

// // TestGitCommandGetBranchGraph is a function.
// func TestGitCommandGetBranchGraph(t *testing.T) {
// 	gitCmd := NewDummyGit()
// 	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
// 		assert.EqualValues(t, "git", cmd)
// 		assert.EqualValues(t, []string{"log", "--graph", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium", "test", "--"}, args)
// 		return secureexec.Command("echo")
// 	}
// 	_, err := gitCmd.GetBranchGraph("test")
// 	assert.NoError(t, err)
// }

// func TestGitCommandGetAllBranchGraph(t *testing.T) {
// 	gitCmd := NewDummyGit()
// 	gitCmd.GetOSCommand().Command = func(cmd string, args ...string) *exec.Cmd {
// 		assert.EqualValues(t, "git", cmd)
// 		assert.EqualValues(t, []string{"log", "--graph", "--all", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium"}, args)
// 		return secureexec.Command("echo")
// 	}
// 	cmdStr := gitCmd.config.GetUserConfig().Git.AllBranchesLogCmd
// 	_, err := gitCmd.GetOSCommand().RunCommandWithOutput(cmdStr)
// 	assert.NoError(t, err)
// }

// // TestGitCommandCurrentBranchName is a function.
// func TestGitCommandCurrentBranchName(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(string, string, error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"says we are on the master branch if we are",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.Equal(t, "git", cmd)
// 				return secureexec.Command("echo", "master")
// 			},
// 			func(name string, displayname string, err error) {
// 				assert.NoError(t, err)
// 				assert.EqualValues(t, "master", name)
// 				assert.EqualValues(t, "master", displayname)
// 			},
// 		},
// 		{
// 			"falls back to git `git branch --contains` if symbolic-ref fails",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)

// 				switch args[0] {
// 				case "symbolic-ref":
// 					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
// 					return secureexec.Command("test")
// 				case "branch":
// 					assert.EqualValues(t, []string{"branch", "--contains"}, args)
// 					return secureexec.Command("echo", "* master")
// 				}

// 				return nil
// 			},
// 			func(name string, displayname string, err error) {
// 				assert.NoError(t, err)
// 				assert.EqualValues(t, "master", name)
// 				assert.EqualValues(t, "master", displayname)
// 			},
// 		},
// 		{
// 			"handles a detached head",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.EqualValues(t, "git", cmd)

// 				switch args[0] {
// 				case "symbolic-ref":
// 					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
// 					return secureexec.Command("test")
// 				case "branch":
// 					assert.EqualValues(t, []string{"branch", "--contains"}, args)
// 					return secureexec.Command("echo", "* (HEAD detached at 123abcd)")
// 				}

// 				return nil
// 			},
// 			func(name string, displayname string, err error) {
// 				assert.NoError(t, err)
// 				assert.EqualValues(t, "123abcd", name)
// 				assert.EqualValues(t, "(HEAD detached at 123abcd)", displayname)
// 			},
// 		},
// 		{
// 			"bubbles up error if there is one",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				assert.Equal(t, "git", cmd)
// 				return secureexec.Command("test")
// 			},
// 			func(name string, displayname string, err error) {
// 				assert.Error(t, err)
// 				assert.EqualValues(t, "", name)
// 				assert.EqualValues(t, "", displayname)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.CurrentBranchName())
// 		})
// 	}
// }

// // TestGitCommandResetHard is a function.
// func TestGitCommandResetHard(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		ref      string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"valid case",
// 			"HEAD",
// 			test.CreateMockCommand(t, []*test.CommandSwapper{
// 				{
// 					Expect:  `git reset --hard HEAD`,
// 					Replace: "echo",
// 				},
// 			}),
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 	}

// 	gitCmd := NewDummyGit()

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.ResetHard(s.ref))
// 		})
// 	}
// }
