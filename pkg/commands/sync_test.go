package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandPush is a function.
func TestGitCommandPush(t *testing.T) {
	type scenario struct {
		testName          string
		getGitConfigValue func(string) (string, error)
		command           func(string, ...string) *exec.Cmd
		forcePush         bool
		test              func(error)
	}

	scenarios := []scenario{
		{
			"Push with force disabled, follow-tags on",
			func(string) (string, error) {
				return "", nil
			},
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--follow-tags"}, args)

				return secureexec.Command("echo")
			},
			false,
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with force enabled, follow-tags on",
			func(string) (string, error) {
				return "", nil
			},
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--follow-tags", "--force-with-lease"}, args)

				return secureexec.Command("echo")
			},
			true,
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with force disabled, follow-tags off",
			func(string) (string, error) {
				return "false", nil
			},
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push"}, args)

				return secureexec.Command("echo")
			},
			false,
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with an error occurring, follow-tags on",
			func(string) (string, error) {
				return "", nil
			},
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--follow-tags"}, args)
				return secureexec.Command("test")
			},
			false,
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGit()
			gitCmd.GetOSCommand().Command = s.command
			gitCmd.getGitConfigValue = s.getGitConfigValue
			err := gitCmd.Push("test", s.forcePush, "", "", func(passOrUname string) string {
				return "\n"
			})
			s.test(err)
		})
	}
}
