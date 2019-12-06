package commands

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/config"
)

// Service is a service that repository is on (Github, Bitbucket, ...)
type Service struct {
	Name           string
	PullRequestURL string
}

// PullRequest opens a link in browser to create new pull request
// with selected branch
type PullRequest struct {
	GitServices []*Service
	GitCommand  *GitCommand
}

// RepoInformation holds some basic information about the repo
type RepoInformation struct {
	Owner      string
	Repository string
}

func getServices(config config.AppConfigurer) []*Service {
	services := []*Service{
		{
			Name:           "github.com",
			PullRequestURL: "https://github.com/%s/%s/compare/%s?expand=1",
		},
		{
			Name:           "bitbucket.org",
			PullRequestURL: "https://bitbucket.org/%s/%s/pull-requests/new?source=%s&t=1",
		},
		{
			Name:           "gitlab.com",
			PullRequestURL: "https://gitlab.com/%s/%s/merge_requests/new?merge_request[source_branch]=%s",
		},
	}

	configServices := config.GetUserConfig().GetStringMapString("services")

	for name, prURL := range configServices {
		services = append(services, &Service{
			Name:           name,
			PullRequestURL: prURL,
		})
	}

	return services
}

// NewPullRequest creates new instance of PullRequest
func NewPullRequest(gitCommand *GitCommand) *PullRequest {
	return &PullRequest{
		GitServices: getServices(gitCommand.Config),
		GitCommand:  gitCommand,
	}
}

// Create opens link to new pull request in browser
func (pr *PullRequest) Create(branch *Branch) error {
	branchExistsOnRemote := pr.GitCommand.CheckRemoteBranchExists(branch)

	if !branchExistsOnRemote {
		return errors.New(pr.GitCommand.Tr.SLocalize("NoBranchOnRemote"))
	}

	repoURL := pr.GitCommand.GetRemoteURL()
	var gitService *Service

	for _, service := range pr.GitServices {
		if strings.Contains(repoURL, service.Name) {
			gitService = service
			break
		}
	}

	if gitService == nil {
		return errors.New(pr.GitCommand.Tr.SLocalize("UnsupportedGitService"))
	}

	repoInfo := getRepoInfoFromURL(repoURL)

	return pr.GitCommand.OSCommand.OpenLink(fmt.Sprintf(
		gitService.PullRequestURL, repoInfo.Owner, repoInfo.Repository, branch.Name,
	))
}

func getRepoInfoFromURL(url string) *RepoInformation {
	isHTTP := strings.HasPrefix(url, "http")

	if isHTTP {
		splits := strings.Split(url, "/")
		owner := splits[len(splits)-2]
		repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

		return &RepoInformation{
			Owner:      owner,
			Repository: repo,
		}
	}

	tmpSplit := strings.Split(url, ":")
	splits := strings.Split(tmpSplit[1], "/")
	owner := splits[0]
	repo := strings.TrimSuffix(splits[1], ".git")

	return &RepoInformation{
		Owner:      owner,
		Repository: repo,
	}
}
