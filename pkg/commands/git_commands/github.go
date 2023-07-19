package git_commands

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type GitHubCommands struct {
	*GitCommon
}

func NewGitHubCommand(gitCommon *GitCommon) *GitHubCommands {
	return &GitHubCommands{
		GitCommon: gitCommon,
	}
}

// https://github.com/cli/cli/issues/2300
func (self *GitHubCommands) BaseRepo() error {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--get-regexp", ".gh-resolved").
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().Run()
}

// Ex: git config --local --add "remote.origin.gh-resolved" "jesseduffield/lazygit"
func (self *GitHubCommands) SetBaseRepo(repository string) (string, error) {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--add", "remote.origin.gh-resolved", repository).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *GitHubCommands) FetchRecentPRs() ([]*models.GithubPullRequest, error) {
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

func (self *GitHubCommands) prList() (string, error) {
	cmdArgs := []string{"gh", "pr", "list", "--limit", "500", "--state", "all", "--json", "state,url,number,headRefName,headRepositoryOwner"}

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

// returns a map from branch name to pull request
func GenerateGithubPullRequestMap(
	prs []*models.GithubPullRequest,
	branches []*models.Branch,
	remotes []*models.Remote,
) map[string]*models.GithubPullRequest {
	res := map[string]*models.GithubPullRequest{}

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

		res[branch.Name] = &pr
	}

	return res
}

func getRemotesToOwnersMap(remotes []*models.Remote) map[string]string {
	res := map[string]string{}
	for _, remote := range remotes {
		if len(remote.Urls) == 0 {
			continue
		}

		res[remote.Name] = getRepoInfoFromURL(remote.Urls[0]).Owner
	}
	return res
}

type RepoInformation struct {
	Owner      string
	Repository string
}

// TODO: move this into hosting_service.go
func getRepoInfoFromURL(url string) RepoInformation {
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

// return <installed>, <valid version>
func (self *GitHubCommands) DetermineGitHubCliState() (bool, bool) {
	output, err := self.cmd.New([]string{"gh", "--version"}).DontLog().RunWithOutput()
	if err != nil {
		// assuming a failure here means that it's not installed
		return false, false
	}

	if !isGhVersionValid(output) {
		return true, false
	}

	return true, true
}

func isGhVersionValid(versionStr string) bool {
	// output should be something like:
	// gh version 2.0.0 (2021-08-23)
	// https://github.com/cli/cli/releases/tag/v2.0.0
	re := regexp.MustCompile(`[^\d]+([\d\.]+)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) == 0 {
		return false
	}

	ghVersion := matches[1]
	majorVersion, err := strconv.Atoi(ghVersion[0:1])
	if err != nil {
		return false
	}
	if majorVersion < 2 {
		return false
	}

	return true
}

func (self *GitHubCommands) InGithubRepo() bool {
	remotes, err := self.repo.Remotes()
	if err != nil {
		self.Log.Error(err)
		return false
	}

	if len(remotes) == 0 {
		return false
	}

	firstRemote := remotes[0]
	if len(firstRemote.Config().URLs) == 0 {
		return false
	}

	url := firstRemote.Config().URLs[0]
	return strings.Contains(url, "github.com")
}
