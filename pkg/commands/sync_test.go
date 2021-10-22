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
		testName string
		command  func(string, ...string) *exec.Cmd
		opts     PushOpts
		test     func(error)
	}

	prompt := func(passOrUname string) string {
		return "\n"
	}

	scenarios := []scenario{
		{
			"Push with force disabled",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{Force: false, PromptUserForCredential: prompt},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with force enabled",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--force-with-lease"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{Force: true, PromptUserForCredential: prompt},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with an error occurring",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push"}, args)
				return secureexec.Command("test")
			},
			PushOpts{Force: false, PromptUserForCredential: prompt},
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"Push with force disabled, upstream supplied",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "origin", "master"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{
				Force:                   false,
				UpstreamRemote:          "origin",
				UpstreamBranch:          "master",
				PromptUserForCredential: prompt,
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with force disabled, setting upstream",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--set-upstream", "origin", "master"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{
				Force:                   false,
				UpstreamRemote:          "origin",
				UpstreamBranch:          "master",
				PromptUserForCredential: prompt,
				SetUpstream:             true,
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with force enabled, setting upstream",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "--force-with-lease", "--set-upstream", "origin", "master"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{
				Force:                   true,
				UpstreamRemote:          "origin",
				UpstreamBranch:          "master",
				PromptUserForCredential: prompt,
				SetUpstream:             true,
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Push with remote branch but no origin",
			func(cmd string, args ...string) *exec.Cmd {
				return nil
			},
			PushOpts{
				Force:                   true,
				UpstreamRemote:          "",
				UpstreamBranch:          "master",
				PromptUserForCredential: prompt,
				SetUpstream:             true,
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "Must specify a remote if specifying a branch", err.Error())
			},
		},
		{
			"Push with force disabled, upstream supplied",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"push", "origin", "master"}, args)

				return secureexec.Command("echo")
			},
			PushOpts{
				Force:                   false,
				UpstreamRemote:          "origin",
				UpstreamBranch:          "master",
				PromptUserForCredential: prompt,
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			err := gitCmd.Push(s.opts)
			s.test(err)
		})
	}
}
