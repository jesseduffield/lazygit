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
		opts              PushOpts
		test              func(error)
	}

	prompt := func(passOrUname string) string {
		return "\n"
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
			PushOpts{Force: false, PromptUserForCredential: prompt},
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
			PushOpts{Force: true, PromptUserForCredential: prompt},
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
			PushOpts{Force: false, PromptUserForCredential: prompt},
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
			PushOpts{Force: false, PromptUserForCredential: prompt},
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"Push with force disabled, follow-tags off, upstream supplied",
			func(string) (string, error) {
				return "false", nil
			},
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
			"Push with force disabled, follow-tags off, setting upstream",
			func(string) (string, error) {
				return "false", nil
			},
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
			"Push with force enabled, follow-tags off, setting upstream",
			func(string) (string, error) {
				return "false", nil
			},
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
			func(string) (string, error) {
				return "false", nil
			},
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
			"Push with force disabled, follow-tags off, upstream supplied",
			func(string) (string, error) {
				return "false", nil
			},
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
			"Push with force disabled, follow-tags off, setting upstream",
			func(string) (string, error) {
				return "false", nil
			},
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
			"Push with force enabled, follow-tags off, setting upstream",
			func(string) (string, error) {
				return "false", nil
			},
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
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			gitCmd.getGitConfigValue = s.getGitConfigValue
			err := gitCmd.Push(s.opts)
			s.test(err)
		})
	}
}

type getPullModeScenario struct {
	testName              string
	getGitConfigValueMock func(string) (string, error)
	configPullModeValue   string
	test                  func(string)
}

func TestGetPullMode(t *testing.T) {

	scenarios := getPullModeScenarios(t)

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.getGitConfigValue = s.getGitConfigValueMock
			s.test(gitCmd.GetPullMode(s.configPullModeValue))
		})
	}
}

func getPullModeScenarios(t *testing.T) []getPullModeScenario {
	return []getPullModeScenario{
		{
			testName: "Merge is default",
			getGitConfigValueMock: func(s string) (string, error) {
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "merge", actual)
			},
		}, {
			testName: "Reads rebase when pull.rebase is true",
			getGitConfigValueMock: func(s string) (string, error) {
				if s == "pull.rebase" {
					return "true", nil
				}
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "rebase", actual)
			},
		}, {
			testName: "Reads ff-only when pull.ff is only",
			getGitConfigValueMock: func(s string) (string, error) {
				if s == "pull.ff" {
					return "only", nil
				}
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "ff-only", actual)
			},
		}, {
			testName: "Reads rebase when rebase is true and ff is only",
			getGitConfigValueMock: func(s string) (string, error) {
				if s == "pull.rebase" {
					return "true", nil
				}
				if s == "pull.ff" {
					return "only", nil
				}
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "rebase", actual)
			},
		}, {
			testName: "Reads rebase when pull.rebase is true",
			getGitConfigValueMock: func(s string) (string, error) {
				if s == "pull.rebase" {
					return "true", nil
				}
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "rebase", actual)
			},
		}, {
			testName: "Reads ff-only when pull.ff is only",
			getGitConfigValueMock: func(s string) (string, error) {
				if s == "pull.ff" {
					return "only", nil
				}
				return "", nil
			},
			configPullModeValue: "auto",
			test: func(actual string) {
				assert.Equal(t, "ff-only", actual)
			},
		}, {
			testName: "Respects merge config",
			getGitConfigValueMock: func(s string) (string, error) {
				return "", nil
			},
			configPullModeValue: "merge",
			test: func(actual string) {
				assert.Equal(t, "merge", actual)
			},
		}, {
			testName: "Respects rebase config",
			getGitConfigValueMock: func(s string) (string, error) {
				return "", nil
			},
			configPullModeValue: "rebase",
			test: func(actual string) {
				assert.Equal(t, "rebase", actual)
			},
		}, {
			testName: "Respects ff-only config",
			getGitConfigValueMock: func(s string) (string, error) {
				return "", nil
			},
			configPullModeValue: "ff-only",
			test: func(actual string) {
				assert.Equal(t, "ff-only", actual)
			},
		},
	}
}
