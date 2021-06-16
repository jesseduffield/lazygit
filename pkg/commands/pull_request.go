package commands

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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
	Git         *Git
}

// RepoInformation holds some basic information about the repo
type RepoInformation struct {
	Owner      string
	Repository string
}

// NewService builds a Service based on the host type
func NewService(typeName string, repositoryDomain string, siteDomain string) *Service {
	var service *Service

	switch typeName {
	case "github":
		service = &Service{
			Name:           repositoryDomain,
			PullRequestURL: fmt.Sprintf("https://%s%s", siteDomain, "/%s/%s/compare/%s?expand=1"),
		}
	case "bitbucket":
		service = &Service{
			Name:           repositoryDomain,
			PullRequestURL: fmt.Sprintf("https://%s%s", siteDomain, "/%s/%s/pull-requests/new?source=%s&t=1"),
		}
	case "gitlab":
		service = &Service{
			Name:           repositoryDomain,
			PullRequestURL: fmt.Sprintf("https://%s%s", siteDomain, "/%s/%s/merge_requests/new?merge_request[source_branch]=%s"),
		}
	}

	return service
}

func getServices(config config.AppConfigurer) []*Service {
	services := []*Service{
		NewService("github", "github.com", "github.com"),
		NewService("bitbucket", "bitbucket.org", "bitbucket.org"),
		NewService("gitlab", "gitlab.com", "gitlab.com"),
	}

	configServices := config.GetUserConfig().Services

	for repoDomain, typeAndDomain := range configServices {
		splitData := strings.Split(typeAndDomain, ":")
		if len(splitData) != 2 {
			// TODO log this misconfiguration
			continue
		}

		service := NewService(splitData[0], repoDomain, splitData[1])
		if service == nil {
			// TODO log this unsupported service
			continue
		}

		services = append(services, service)
	}

	return services
}

// NewPullRequest creates new instance of PullRequest
func NewPullRequest(gitCommand *Git) *PullRequest {
	return &PullRequest{
		GitServices: getServices(gitCommand.config),
		Git:         gitCommand,
	}
}

// Create opens link to new pull request in browser
func (pr *PullRequest) Create(branch *models.Branch) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(branch)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.Git.GetOS().OpenLink(pullRequestURL)
}

// CopyURL copies the pull request URL to the clipboard
func (pr *PullRequest) CopyURL(branch *models.Branch) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(branch)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.Git.GetOS().CopyToClipboard(pullRequestURL)
}

func (pr *PullRequest) getPullRequestURL(branch *models.Branch) (string, error) {
	branchExistsOnRemote := pr.Git.Remotes().RemoteBranchExists(branch)

	if !branchExistsOnRemote {
		return "", errors.New(pr.Git.tr.NoBranchOnRemote)
	}

	repoURL := pr.Git.Remotes().GetCurrentRemoteUrl()
	var gitService *Service

	for _, service := range pr.GitServices {
		if strings.Contains(repoURL, service.Name) {
			gitService = service
			break
		}
	}

	if gitService == nil {
		return "", errors.New(pr.Git.tr.UnsupportedGitService)
	}

	repoInfo := getRepoInfoFromURL(repoURL)
	pullRequestURL := fmt.Sprintf(
		gitService.PullRequestURL, repoInfo.Owner, repoInfo.Repository, branch.Name,
	)

	return pullRequestURL, nil
}

func getRepoInfoFromURL(url string) *RepoInformation {
	isHTTP := strings.HasPrefix(url, "http")

	if isHTTP {
		splits := strings.Split(url, "/")
		owner := strings.Join(splits[3:len(splits)-1], "/")
		repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

		return &RepoInformation{
			Owner:      owner,
			Repository: repo,
		}
	}

	tmpSplit := strings.Split(url, ":")
	splits := strings.Split(tmpSplit[1], "/")
	owner := strings.Join(splits[0:len(splits)-1], "/")
	repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

	return &RepoInformation{
		Owner:      owner,
		Repository: repo,
	}
}
