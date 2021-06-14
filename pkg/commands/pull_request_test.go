package commands

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGetRepoInfoFromURL is a function.
func TestGetRepoInfoFromURL(t *testing.T) {
	type scenario struct {
		testName string
		repoURL  string
		test     func(*RepoInformation)
	}

	scenarios := []scenario{
		{
			"Returns repository information for git remote url",
			"git@github.com:petersmith/super_calculator",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "petersmith")
				assert.EqualValues(t, repoInfo.Repository, "super_calculator")
			},
		},
		{
			"Returns repository information for http remote url",
			"https://my_username@bitbucket.org/johndoe/social_network.git",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "johndoe")
				assert.EqualValues(t, repoInfo.Repository, "social_network")
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(getRepoInfoFromURL(s.repoURL))
		})
	}
}

// TestCreatePullRequest is a function.
func TestCreatePullRequest(t *testing.T) {
	type scenario struct {
		testName  string
		branch    *models.Branch
		remoteUrl string
		command   func(string, ...string) *exec.Cmd
		test      func(url string, err error)
	}

	scenarios := []scenario{
		{
			testName: "Opens a link to new pull request on bitbucket",
			branch: &models.Branch{
				Name: "feature/profile-page",
			},
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@bitbucket.org:johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page&t=1", url)
			},
		},
		{
			testName: "Opens a link to new pull request on bitbucket with http remote url",
			branch: &models.Branch{
				Name: "feature/events",
			},
			remoteUrl: "https://my_username@bitbucket.org/johndoe/social_network.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "https://my_username@bitbucket.org/johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/events&t=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/events&t=1", url)
			},
		},
		{
			testName: "Opens a link to new pull request on github",
			branch: &models.Branch{
				Name: "feature/sum-operation",
			},
			remoteUrl: "git@github.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@github.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://github.com/peter/calculator/compare/feature/sum-operation?expand=1"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature/sum-operation?expand=1", url)
			},
		},
		{
			testName: "Opens a link to new pull request on gitlab",
			branch: &models.Branch{
				Name: "feature/ui",
			},
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			command: func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return secureexec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/ui"})
				return secureexec.Command("echo")
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/ui", url)
			},
		},
		{
			testName: "Throws an error if git service is unsupported",
			branch: &models.Branch{
				Name: "feature/divide-operation",
			},
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
			gitCommand := NewDummyGit()
			gitCommand.GetOSCommand().Command = s.command
			gitCommand.GetOSCommand().Config.GetUserConfig().OS.OpenLinkCommand = "open {{link}}"
			gitCommand.GetOSCommand().Config.GetUserConfig().Services = map[string]string{
				// valid configuration for a custom service URL
				"git.work.com": "gitlab:code.work.com",
				// invalid configurations for a custom service URL
				"invalid.work.com":   "noservice:invalid.work.com",
				"noservice.work.com": "noservice.work.com",
			}
			gitCommand.getGitConfigValue = func(path string) (string, error) {
				assert.Equal(t, path, "remote.origin.url")
				return s.remoteUrl, nil
			}
			dummyPullRequest := NewPullRequest(gitCommand)
			s.test(dummyPullRequest.Create(s.branch))
		})
	}
}
