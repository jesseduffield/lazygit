package git_commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type GhCommands struct {
	*GitCommon
}

func NewGhCommand(gitCommon *GitCommon) *GhCommands {
	return &GhCommands{
		GitCommon: gitCommon,
	}
}

// https://github.com/cli/cli/issues/2300
func (self *GhCommands) BaseRepo() error {
	return self.cmd.New("git config --local --get-regexp .gh-resolved").Run()
}

// Ex: git config --local --add "remote.origin.gh-resolved" "jesseduffield/lazygit"
func (self *GhCommands) SetBaseRepo(repository string) (string, error) {
	return self.cmd.New(
		fmt.Sprintf("git config --local --add \"remote.origin.gh-resolved\" \"%s\"", repository),
	).RunWithOutput()
}

func (self *GhCommands) prList() (string, error) {
	return self.cmd.New(
		"gh pr list --limit 500 --state all --json state,url,number,headRefName,headRepositoryOwner",
	).RunWithOutput()
}

func (self *GhCommands) GithubMostRecentPRs() ([]*models.GithubPullRequest, error) {
	commandOutput, err := self.prList()
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

func GenerateGithubPullRequestMap(prs []*models.GithubPullRequest, branches []*models.Branch, remotes []*models.Remote) map[*models.Branch]*models.GithubPullRequest {
	res := map[*models.Branch]*models.GithubPullRequest{}

	if len(prs) == 0 {
		return res
	}

	remotesToOwnersMap := getRemotesToOwnersMap(remotes)

	if len(remotesToOwnersMap) == 0 {
		return res
	}

	// A PR can be identified by two things: the owner e.g. 'jesseduffield' and the
	// branch name e.g. 'feature/my-feature'. The owner might be different
	// to the owner of the repo if the PR is from a fork of that repo.
	type prKey struct {
		owner      string
		branchName string
	}

	prByKey := map[prKey]models.GithubPullRequest{}

	for _, pr := range prs {
		prByKey[prKey{owner: pr.UserName(), branchName: pr.BranchName()}] = *pr
	}

	for _, branch := range branches {
		if !branch.IsTrackingRemote() {
			continue
		}

		// TODO: support branches whose UpstreamRemote contains a full git
		// URL rather than just a remote name.
		owner, foundRemoteOwner := remotesToOwnersMap[branch.UpstreamRemote]
		if !foundRemoteOwner {
			continue
		}

		pr, hasPr := prByKey[prKey{owner: owner, branchName: branch.UpstreamBranch}]

		if !hasPr {
			continue
		}

		res[branch] = &pr
	}

	return res
}

func getRemotesToOwnersMap(remotes []*models.Remote) map[string]string {
	res := map[string]string{}
	for _, remote := range remotes {
		if len(remote.Urls) == 0 {
			continue
		}

		res[remote.Name] = GetRepoInfoFromURL(remote.Urls[0]).Owner
	}
	return res
}

type RepoInformation struct {
	Owner      string
	Repository string
}

// TODO: move this into hosting_service.go
func GetRepoInfoFromURL(url string) RepoInformation {
	isHTTP := strings.HasPrefix(url, "http")

	if isHTTP {
		splits := strings.Split(url, "/")
		owner := strings.Join(splits[3:len(splits)-1], "/")
		repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

		return RepoInformation{
			Owner:      owner,
			Repository: repo,
		}
	}

	tmpSplit := strings.Split(url, ":")
	splits := strings.Split(tmpSplit[1], "/")
	owner := strings.Join(splits[0:len(splits)-1], "/")
	repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

	return RepoInformation{
		Owner:      owner,
		Repository: repo,
	}
}
