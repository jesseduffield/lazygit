package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

func TestGithubMostRecentPRs(t *testing.T) {
	scenarios := []struct {
		testName string
		response string
		expect   map[string]models.GithubPullRequest
	}{
		{
			"no response",
			"",
			nil,
		},
		{
			"error response",
			"none of the git remotes configured for this repository point to a known GitHub host. To tell gh about a new GitHub host, please use `gh auth login`",
			nil,
		},
		{
			"empty response",
			"[]",
			map[string]models.GithubPullRequest{},
		},
		{
			"response with data",
			`[{
				"headRefName": "command-log-2",
				"number": 1249,
				"state": "MERGED",
				"url": "https://github.com/jesseduffield/lazygit/pull/1249",
				"headRepositoryOwner": {
					"id": "MDQ6VXNlcjg0NTY2MzM=",
					"name": "Jesse Duffield",
					"login": "jesseduffield"
				}
			}]`,
			map[string]models.GithubPullRequest{
				"jesseduffield:command-log-2": {
					HeadRefName: "command-log-2",
					Number:      1249,
					State:       "MERGED",
					Url:         "https://github.com/jesseduffield/lazygit/pull/1249",
					HeadRepositoryOwner: models.GithubRepositoryOwner{
						ID:    "MDQ6VXNlcjg0NTY2MzM=",
						Name:  "Jesse Duffield",
						Login: "jesseduffield",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "gh", cmd)
				assert.EqualValues(t, []string{"pr", "list", "--limit", "50", "--state", "all", "--json", "state,url,number,headRefName,headRepositoryOwner"}, args)
				return secureexec.Command("echo", s.response)
			}

			res := gitCmd.GithubMostRecentPRs()
			assert.Equal(t, s.expect, res)
		})
	}
}
