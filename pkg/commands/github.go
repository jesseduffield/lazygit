package commands

import (
	"encoding/json"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

func (c *GitCommand) GithubMostRecentPRs() (map[string]models.GithubPullRequest, error) {
	commandOutput, err := c.OSCommand.RunCommandWithOutput("gh pr list --limit 50 --state all --json state,url,number,headRefName,headRepositoryOwner")
	if err != nil {
		return nil, err
	}

	prs := []models.GithubPullRequest{}
	err = json.Unmarshal([]byte(commandOutput), &prs)
	if err != nil {
		return nil, err
	}

	res := map[string]models.GithubPullRequest{}
	for _, pr := range prs {
		res[pr.HeadRepositoryOwner.Login+":"+pr.HeadRefName] = pr
	}
	return res, nil
}

func (c *GitCommand) FoundBranchWithGithubPullRequest(prs map[string]models.GithubPullRequest, branches []*models.Branch) bool {
	if len(prs) == 0 {
		return false
	}

	remotesToOwnersMap, _ := c.GetRemotesToOwnersMap()
	if len(remotesToOwnersMap) == 0 {
		return false
	}

	foundBranchWithGithubPullRequest := false

	for _, branch := range branches {
		_, has_pr := presentation.GetPr(branch, remotesToOwnersMap, prs)

		if has_pr {
			foundBranchWithGithubPullRequest = true
		}
	}

	return foundBranchWithGithubPullRequest
}
