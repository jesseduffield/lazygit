package commands

import (
	"os/exec"
	"strings"
	"testing"

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
		branch   *Branch
		command  func(string, ...string) *exec.Cmd
		test     func(err error)
	}

	scenarios := []scenario{
		{
			"Opens a link to new pull request on bitbucket",
			&Branch{
				Name: "feature/profile-page",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "git@bitbucket.org:johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?t=feature/profile-page"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Opens a link to new pull request on bitbucket with http remote url",
			&Branch{
				Name: "feature/events",
			},
			func(cmd string, args ...string) *exec.Cmd {
				// Handle git remote url call
				if strings.HasPrefix(cmd, "git") {
					return exec.Command("echo", "https://my_username@bitbucket.org/johndoe/social_network.git")
				}

				assert.Equal(t, cmd, "open")
				assert.Equal(t, args, []string{"https://bitbucket.org/johndoe/social_network/pull-requests/new?t=feature/events"})
				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Opens a link to new pull request on github",
			&Branch{
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
			&Branch{
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
			&Branch{
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
			dummyPullRequest := NewPullRequest(gitCommand)
			s.test(dummyPullRequest.Create(s.branch))
		})
	}
}
