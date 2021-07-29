package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (c *GitCommand) GithubMostRecentPRs() map[string]models.GithubPullRequest {
	commandOutput, err := c.OSCommand.RunCommandWithOutput("gh pr list --limit 50 --state all --json state,url,number,headRefName,headRepositoryOwner")
	if err != nil {
		fmt.Println(1, err)
		return nil
	}

	prs := []models.GithubPullRequest{}
	err = json.Unmarshal([]byte(commandOutput), &prs)
	if err != nil {
		fmt.Println(2, err)
		return nil
	}

	res := map[string]models.GithubPullRequest{}
	for _, pr := range prs {
		res[pr.HeadRepositoryOwner.Login+":"+pr.HeadRefName] = pr
	}
	return res
}

func (c *GitCommand) InjectGithubPullRequests(prs map[string]models.GithubPullRequest, branches []*models.Branch) bool {
	if len(prs) == 0 {
		return false
	}

	remotesToOwnersMap, _ := c.GetRemotesToOwnersMap()
	if len(remotesToOwnersMap) == 0 {
		return false
	}

	foundBranchWithGithubPullRequest := false

	for _, branch := range branches {
		if branch.UpstreamName == "" {
			continue
		}

		remoteAndName := strings.SplitN(branch.UpstreamName, "/", 2)
		owner, foundRemoteOwner := remotesToOwnersMap[remoteAndName[0]]
		if len(remoteAndName) != 2 || !foundRemoteOwner {
			continue
		}

		pr, hasPr := prs[owner+":"+remoteAndName[1]]
		if !hasPr {
			continue
		}

		foundBranchWithGithubPullRequest = true
		branch.PR = &pr
	}

	return foundBranchWithGithubPullRequest
}
