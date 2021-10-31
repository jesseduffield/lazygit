package commands

import (
	"encoding/json"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
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

func (c *GitCommand) GenerateGithubPullRequestMap(prs map[string]models.GithubPullRequest, branches []*models.Branch) (map[*models.Branch]*models.GithubPullRequest, bool) {
	res := map[*models.Branch]*models.GithubPullRequest{}

	if len(prs) == 0 {
		return res, false
	}

	remotesToOwnersMap, _ := c.GetRemotesToOwnersMap()
	if len(remotesToOwnersMap) == 0 {
		return res, false
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

		res[branch] = &pr
	}

	return res, foundBranchWithGithubPullRequest
}
