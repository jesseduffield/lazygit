package git_commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/auth"
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
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

type Response struct {
	Data RepositoryQuery `json:"data"`
}

type RepositoryQuery struct {
	Repository map[string]PullRequest `json:"repository"`
}

type PullRequest struct {
	Edges []PullRequestEdge `json:"edges"`
}

type PullRequestEdge struct {
	Node PullRequestNode `json:"node"`
}

type PullRequestNode struct {
	Title               string                `json:"title"`
	HeadRefName         string                `json:"headRefName"`
	Number              int                   `json:"number"`
	Url                 string                `json:"url"`
	HeadRepositoryOwner GithubRepositoryOwner `json:"headRepositoryOwner"`
	State               string                `json:"state"`
}

type GithubRepositoryOwner struct {
	Login string `json:"login"`
}

func fetchPullRequestsQuery(branches []string, owner string, repo string) string {
	var queries []string
	for i, branch := range branches {
		// We're making a sub-query per branch, and arbitrarily labelling each subquery
		// as a1, a2, etc.
		fieldName := fmt.Sprintf("a%d", i+1)
		// TODO: scope down by remote too if we can (right now if you search for master, you can get multiple results back, and all from forks)
		queries = append(queries, fmt.Sprintf(`%s: pullRequests(first: 1, headRefName: "%s") {
      edges {
        node {
          title
          headRefName
          state
          number
          url
          headRepositoryOwner {
            login
          }
        }
      }
    }`, fieldName, branch))
	}

	queryString := fmt.Sprintf(`{
  repository(owner: "%s", name: "%s") {
    %s
  }
}`, owner, repo, strings.Join(queries, "\n"))

	return queryString
}

// FetchRecentPRs fetches recent pull requests using GraphQL.
func (self *GitHubCommands) FetchRecentPRs(branches []string) ([]*models.GithubPullRequest, error) {
	repoOwner, repoName, err := self.GetBaseRepoOwnerAndName()
	if err != nil {
		return nil, err
	}

	t := time.Now()

	var g errgroup.Group
	results := make(chan []*models.GithubPullRequest)

	// We want at most 5 concurrent requests, but no less than 10 branches per request
	concurrency := 5
	minBranchesPerRequest := 10
	branchesPerRequest := max(len(branches)/concurrency, minBranchesPerRequest)
	for i := 0; i < len(branches); i += branchesPerRequest {
		end := i + branchesPerRequest
		if end > len(branches) {
			end = len(branches)
		}
		branchChunk := branches[i:end]

		// Launch a goroutine for each chunk of branches
		g.Go(func() error {
			prs, err := self.FetchRecentPRsAux(repoOwner, repoName, branchChunk)
			if err != nil {
				return err
			}
			results <- prs
			return nil
		})
	}

	// Close the results channel when all goroutines are done
	go func() {
		g.Wait()
		close(results)
	}()

	// Collect results from all goroutines
	var allPRs []*models.GithubPullRequest
	for prs := range results {
		allPRs = append(allPRs, prs...)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	self.Log.Warnf("Fetched PRs in %s", time.Since(t))

	return allPRs, nil
}

func (self *GitHubCommands) FetchRecentPRsAux(repoOwner string, repoName string, branches []string) ([]*models.GithubPullRequest, error) {
	queryString := fetchPullRequestsQuery(branches, repoOwner, repoName)
	escapedQueryString := strconv.Quote(queryString)

	body := fmt.Sprintf(`{"query": %s}`, escapedQueryString)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	defaultHost, _ := auth.DefaultHost()
	token, _ := auth.TokenForHost(defaultHost)
	if token == "" {
		return nil, fmt.Errorf("No token found for GitHub")
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyStr := new(bytes.Buffer)
		bodyStr.ReadFrom(resp.Body)
		return nil, fmt.Errorf("GraphQL query failed with status: %s. Body: %s", resp.Status, bodyStr.String())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	prs := []*models.GithubPullRequest{}
	for _, repoQuery := range result.Data.Repository {
		for _, edge := range repoQuery.Edges {
			node := edge.Node
			pr := &models.GithubPullRequest{
				HeadRefName: node.HeadRefName,
				Number:      node.Number,
				State:       node.State,
				Url:         node.Url,
				HeadRepositoryOwner: models.GithubRepositoryOwner{
					Login: node.HeadRepositoryOwner.Login,
				},
			}
			prs = append(prs, pr)
		}
	}

	return prs, nil
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

	remote := GetMainRemote(remotes)

	if len(remote.Config().URLs) == 0 {
		return false
	}

	url := remote.Config().URLs[0]
	return strings.Contains(url, "github.com")
}

func GetMainRemote(remotes []*gogit.Remote) *gogit.Remote {
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			return remote
		}
	}

	// need to sort remotes by name so that this is deterministic
	return lo.MinBy(remotes, func(a, b *gogit.Remote) bool {
		return a.Config().Name < b.Config().Name
	})
}

func GetSuggestedRemoteName(remotes []*models.Remote) string {
	if len(remotes) == 0 {
		return "origin"
	}

	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote.Name
		}
	}

	return remotes[0].Name
}

func (self *GitHubCommands) GetBaseRepoOwnerAndName() (string, string, error) {
	remotes, err := self.repo.Remotes()
	if err != nil {
		return "", "", err
	}

	if len(remotes) == 0 {
		return "", "", fmt.Errorf("No remotes found")
	}

	firstRemote := remotes[0]
	if len(firstRemote.Config().URLs) == 0 {
		return "", "", fmt.Errorf("No URLs found for remote")
	}

	url := firstRemote.Config().URLs[0]

	repoInfo := getRepoInfoFromURL(url)

	return repoInfo.Owner, repoInfo.Repository, nil
}
