package commands

import (
	"encoding/json"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (c *GitCommand) GithubMostRecentPRs() ([]*models.GithubPullRequest, error) {
	commandOutput, err := c.OSCommand.RunCommandWithOutput("gh pr list --limit 50 --state all --json state,url,number,headRefName,headRepositoryOwner")
	if err != nil {
		return nil, err
	}

	prs := []*models.GithubPullRequest{}
	err = json.Unmarshal([]byte(commandOutput), &prs)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (c *GitCommand) GenerateGithubPullRequestMap(prs []*models.GithubPullRequest, branches []*models.Branch) map[*models.Branch]*models.GithubPullRequest {
	res := map[*models.Branch]*models.GithubPullRequest{}

	if len(prs) == 0 {
		return res
	}

	remotesToOwnersMap, _ := c.GetRemotesToOwnersMap()
	if len(remotesToOwnersMap) == 0 {
		return res
	}

	prWithStringKey := map[string]models.GithubPullRequest{}

	for _, pr := range prs {
		prWithStringKey[pr.UserName()+":"+pr.BranchName()] = *pr
	}

	for _, branch := range branches {
		if !branch.IsTrackingRemote() {
			continue
		}

		owner, foundRemoteOwner := remotesToOwnersMap[branch.RemoteName()]
		if branch.BranchName() == "" || !foundRemoteOwner {
			continue
		}

		pr, hasPr := prWithStringKey[owner+":"+branch.BranchName()]
		if !hasPr {
			continue
		}

		res[branch] = &pr
	}

	return res
}
