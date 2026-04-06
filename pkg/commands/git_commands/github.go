package git_commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type GitHubCommands struct {
	*GitCommon
}

func NewGitHubCommands(gitCommon *GitCommon) *GitHubCommands {
	return &GitHubCommands{
		GitCommon: gitCommon,
	}
}

// https://github.com/cli/cli/issues/2300
func (self *GitHubCommands) ConfiguredBaseRemoteName() string {
	// TODO: we only support the (common) case where the value of the config is "base", meaning that
	// the remote's URL determines the GitHub repo. Since `gh repo set-default` on the command line
	// sets the config this way, it's probably good enough in practice, but for completeness it
	// would be nice to also support the case where the config value is a full remote name (e.g.
	// "jesseduffield/lazygit").

	cmdArgs := NewGitCmd("config").
		Arg("--local", "--get-regexp", `remote\..*\.gh-resolved`).
		ToArgv()

	output, _, err := self.cmd.New(cmdArgs).DontLog().RunWithOutputs()
	if err != nil {
		return ""
	}

	regex := regexp.MustCompile(`remote\.(.+)\.gh-resolved`)
	matches := regex.FindStringSubmatch(output)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func (self *GitHubCommands) SetConfiguredBaseRemoteName(remoteName string) error {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--add", fmt.Sprintf("remote.%s.gh-resolved", remoteName), "base").
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().Run()
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
	IsDraft             bool                  `json:"isDraft"`
}

type GithubRepositoryOwner struct {
	Login string `json:"login"`
}

type graphQLRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func fetchPullRequestsQuery(branches []string, owner string, repo string) (string, map[string]string) {
	variables := make(map[string]string, len(branches)+2)
	variables["owner"] = owner
	variables["repo"] = repo
	varDecls := make([]string, 0, len(branches)+2)
	varDecls = append(varDecls, "$owner: String!", "$repo: String!")
	queries := make([]string, 0, len(branches))
	for i, branch := range branches {
		// We're making a sub-query per branch, and arbitrarily labelling each subquery
		// as a1, a2, etc.
		fieldName := fmt.Sprintf("a%d", i+1)
		varName := fmt.Sprintf("branch%d", i+1)
		variables[varName] = branch
		varDecls = append(varDecls, fmt.Sprintf("$%s: String!", varName))
		// We fetch a few PRs per branch name because multiple forks may have PRs
		// with the same head ref name. The mapping logic filters by owner later.
		queries = append(queries, fmt.Sprintf(`%s: pullRequests(first: 5, headRefName: $%s, orderBy: {field: CREATED_AT, direction: DESC}) {
      edges {
        node {
          title
          headRefName
          state
          number
          url
          isDraft
          headRepositoryOwner {
            login
          }
        }
      }
    }`, fieldName, varName))
	}

	queryString := fmt.Sprintf(`query(%s) {
  repository(owner: $owner, name: $repo) {
    %s
  }
}`, strings.Join(varDecls, ", "), strings.Join(queries, "\n"))

	return queryString, variables
}

func (self *GitHubCommands) GetAuthToken() string {
	defaultHost, _ := auth.DefaultHost()
	token, _ := auth.TokenForHost(defaultHost)
	return token
}

// FetchRecentPRs fetches recent pull requests using GraphQL.
func (self *GitHubCommands) FetchRecentPRs(branches []string, baseRemote *models.Remote, token string) ([]*models.GithubPullRequest, error) {
	repoOwner, repoName, err := self.GetBaseRepoOwnerAndName(baseRemote)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	var g errgroup.Group

	// We want at most 5 concurrent requests, but no less than 10 branches per request
	concurrency := 5
	minBranchesPerRequest := 10
	branchesPerRequest := max(len(branches)/concurrency, minBranchesPerRequest)
	numChunks := (len(branches) + branchesPerRequest - 1) / branchesPerRequest
	results := make(chan []*models.GithubPullRequest, numChunks)

	for i := 0; i < len(branches); i += branchesPerRequest {
		end := i + branchesPerRequest
		if end > len(branches) {
			end = len(branches)
		}
		branchChunk := branches[i:end]

		// Launch a goroutine for each chunk of branches
		g.Go(func() error {
			prs, err := self.fetchRecentPRsAux(repoOwner, repoName, branchChunk, token)
			if err != nil {
				return err
			}
			results <- prs
			return nil
		})
	}

	// Wait for all goroutines, then close the channel so the range loop exits
	err = g.Wait()
	close(results)
	if err != nil {
		return nil, err
	}

	// Collect results from all goroutines
	var allPRs []*models.GithubPullRequest
	for prs := range results {
		allPRs = append(allPRs, prs...)
	}

	self.Log.Infof("Fetched %d PRs in %s", len(allPRs), time.Since(t))

	return allPRs, nil
}

func (self *GitHubCommands) fetchRecentPRsAux(repoOwner string, repoName string, branches []string, token string) ([]*models.GithubPullRequest, error) {
	queryString, variables := fetchPullRequestsQuery(branches, repoOwner, repoName)

	bodyBytes, err := json.Marshal(graphQLRequest{Query: queryString, Variables: variables})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
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
		_, _ = bodyStr.ReadFrom(resp.Body)
		return nil, fmt.Errorf("GraphQL query failed with status: %s. Body: %s", resp.Status, bodyStr.String())
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(respBytes, &result)
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
				Title:       node.Title,
				State:       lo.Ternary(node.IsDraft && node.State != "CLOSED", "DRAFT", node.State),
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

	// A PR can be identified by two things: the owner e.g. 'jesseduffield' and the
	// branch name e.g. 'feature/my-feature'. The owner might be different
	// to the owner of the repo if the PR is from a fork of that repo.
	type prKey struct {
		owner      string
		branchName string
	}

	prByKey := map[prKey]models.GithubPullRequest{}

	for _, pr := range prs {
		key := prKey{owner: pr.UserName(), branchName: pr.BranchName()}
		// PRs are returned newest-first from the API, so the first one we
		// see for each key is the most recent and therefore the most relevant.
		if _, exists := prByKey[key]; !exists {
			prByKey[key] = *pr
		}
	}

	for _, branch := range branches {
		if !branch.IsTrackingRemote() {
			continue
		}

		owner, foundRemoteOwner := remotesToOwnersMap[branch.UpstreamRemote]
		if !foundRemoteOwner {
			// UpstreamRemote may be a full URL rather than a remote name;
			// try parsing the owner directly from it.
			repoInfo, err := hosting_service.GetRepoInfoFromURL(branch.UpstreamRemote)
			if err != nil {
				continue
			}
			owner = repoInfo.Owner
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

		repoInfo, err := hosting_service.GetRepoInfoFromURL(remote.Urls[0])
		if err != nil {
			continue
		}

		res[remote.Name] = repoInfo.Owner
	}
	return res
}

func (self *GitHubCommands) InGithubRepo(remotes []*models.Remote) bool {
	if len(remotes) == 0 {
		return false
	}

	remote := getMainRemote(remotes)

	if len(remote.Urls) == 0 {
		return false
	}

	url := remote.Urls[0]
	return strings.Contains(url, "github.com")
}

func getMainRemote(remotes []*models.Remote) *models.Remote {
	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote
		}
	}

	// need to sort remotes by name so that this is deterministic
	return lo.MinBy(remotes, func(a, b *models.Remote) bool {
		return a.Name < b.Name
	})
}

func (self *GitHubCommands) GetBaseRepoOwnerAndName(baseRemote *models.Remote) (string, string, error) {
	if len(baseRemote.Urls) == 0 {
		return "", "", fmt.Errorf("No URLs found for remote")
	}

	url := baseRemote.Urls[0]

	repoInfo, err := hosting_service.GetRepoInfoFromURL(url)
	if err != nil {
		return "", "", err
	}

	return repoInfo.Owner, repoInfo.Repository, nil
}
