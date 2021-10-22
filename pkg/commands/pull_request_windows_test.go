//go:build windows
// +build windows

package commands

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestCreatePullRequestOnWindows is a function.
func TestCreatePullRequestOnWindows(t *testing.T) {
	type scenario struct {
		testName  string
		from      string
		to        string
		remoteUrl string
		command   func(string, ...string) *exec.Cmd
		test      func(url string, err error)
	}

	scenarios := []scenario{
		{
			testName:  "Opens a link to new pull request on bitbucket",
			from:      "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@bitbucket.org:johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page^&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with http remote url",
			from:      "feature/events",
			remoteUrl: "https://my_username@bitbucket.org/johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "https://my_username@bitbucket.org/johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/events^&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/events&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on github",
			from:      "feature/sum-operation",
			remoteUrl: "git@github.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@github.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://github.com/peter/calculator/compare/feature/sum-operation?expand=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature/sum-operation?expand=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with specific target branch",
			from:      "feature/profile-page/avatar",
			to:        "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@bitbucket.org:johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page/avatar^&dest=feature/profile-page^&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page/avatar&dest=feature/profile-page&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with http remote url with specified target branch",
			from:      "feature/remote-events",
			to:        "feature/events",
			remoteUrl: "https://my_username@bitbucket.org/johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "https://my_username@bitbucket.org/johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/remote-events^&dest=feature/events^&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/remote-events&dest=feature/events&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on github with specific target branch",
			from:      "feature/sum-operation",
			to:        "feature/operations",
			remoteUrl: "git@github.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@github.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://github.com/peter/calculator/compare/feature/operations...feature/sum-operation?expand=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature/operations...feature/sum-operation?expand=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab",
			from:      "feature/ui",
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/ui"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/ui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab in nested groups",
			from:      "feature/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature/ui"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature/ui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/commit-ui^&merge_request[target_branch]=epic/ui"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/commit-ui&merge_request[target_branch]=epic/ui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch in nested groups",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "cmd")
				assert.Equal(t, args, []string{"/c", "start", "", "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature/commit-ui^&merge_request[target_branch]=epic/ui"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature/commit-ui&merge_request[target_branch]=epic/ui", url)
			},
		},
		{
			testName:  "Throws an error if git service is unsupported",
			from:      "feature/divide-operation",
			remoteUrl: "git@something.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCommand := NewDummyGitCommand()
			gitCommand.OSCommand.Command = s.command
			gitCommand.OSCommand.Platform.OS = "windows"
			gitCommand.OSCommand.Platform.Shell = "cmd"
			gitCommand.OSCommand.Platform.ShellArg = "/c"
			gitCommand.OSCommand.Config.GetUserConfig().OS.OpenLinkCommand = `start "" {{link}}`
			gitCommand.OSCommand.Config.GetUserConfig().Services = map[string]string{
				// valid configuration for a custom service URL
				"git.work.com": "gitlab:code.work.com",
				// invalid configurations for a custom service URL
				"invalid.work.com":   "noservice:invalid.work.com",
				"noservice.work.com": "noservice.work.com",
			}
			gitCommand.GitConfig = git_config.NewFakeGitConfig(map[string]string{"remote.origin.url": s.remoteUrl})
			dummyPullRequest := NewPullRequest(gitCommand)
			s.test(dummyPullRequest.Create(s.from, s.to))
		})
	}
}
