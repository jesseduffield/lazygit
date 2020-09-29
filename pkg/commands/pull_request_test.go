package commands

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/models"
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
		testName string
		branch   *models.Branch
		command  func(string, ...string) *exec.Cmd
		test     func(err error)
	}

	scenarios := []scenario{
		{
			"Opens a link to new pull request on bitbucket",
			&models.Branch{
				Name: "feature/profile-page",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "git@bitbucket.org:johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/profile-page&t=1"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Opens a link to new pull request on bitbucket with http remote url",
			&models.Branch{
				Name: "feature/events",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "https://my_username@bitbucket.org/johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature/events&t=1"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Opens a link to new pull request on github",
			&models.Branch{
				Name: "feature/sum-operation",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "git@github.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://github.com/peter/calculator/compare/feature/sum-operation?expand=1"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Opens a link to new pull request on gitlab",
			&models.Branch{
				Name: "feature/ui",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "git@gitlab.com:peter/calculator.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature/ui"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Throws an error if git service is unsupported",
			&models.Branch{
				Name: "feature/divide-operation",
			},
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("echo", "git@something.com:peter/calculator.git")
			},
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCommand := NewDummyGitCommand()
			gitCommand.OSCommand.command = s.command
			gitCommand.OSCommand.Config.GetUserConfig().Set("os.openLinkCommand", "open {{link}}")
			gitCommand.Config.GetUserConfig().Set("services", map[string]string{
				// valid configuration for a custom service URL
				"git.work.com": "gitlab:code.work.com",
				// invalid configurations for a custom service URL
				"invalid.work.com":   "noservice:invalid.work.com",
				"noservice.work.com": "noservice.work.com",
			})
			dummyPullRequest := NewPullRequest(gitCommand)
			s.test(dummyPullRequest.Create(s.branch))
		})
	}
}
