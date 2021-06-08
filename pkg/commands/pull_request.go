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
	PullRequestURL func(owner string, repository string, from string, to string) string
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

// NewService builds a Service based on the host type
func NewService(typeName string, repositoryDomain string, siteDomain string) *Service {
	var service *Service

	switch typeName {
	case "github":
		service = &Service{
			Name: repositoryDomain,
			PullRequestURL: func(owner string, repository string, from string, to string) string {
				if to == "" {
					return fmt.Sprintf("https://%s/%s/%s/compare/%s?expand=1",siteDomain, owner, repository, from)
				} else {
					return fmt.Sprintf("https://%s/%s/%s/compare/%s...%s?expand=1", siteDomain, owner, repository, to, from)
				}
			},
		}
	case "bitbucket":
		service = &Service{
			Name: repositoryDomain,
			PullRequestURL: func(owner string, repository string, from string, to string) string {
				if to == "" {
					return fmt.Sprintf("https://%s/%s/%s/pull-requests/new?source=%s&t=1", siteDomain, owner, repository, from)
				} else {
					return fmt.Sprintf("https://%s/%s/%s/pull-requests/new?source=%s&dest=%s&t=1", siteDomain, owner, repository, from, to)
				}
			},
		}
	case "gitlab":
		service = &Service{
			Name: repositoryDomain,
			PullRequestURL: func(owner string, repository string, from string, to string) string {
				if to == "" {
					return fmt.Sprintf("https://%s/%s/%s/merge_requests/new?merge_request[source_branch]=%s", siteDomain, owner, repository, from)
				} else {
					return fmt.Sprintf("https://%s/%s/%s/merge_requests/new?merge_request[source_branch]=%s&merge_request[target_branch]=%s", siteDomain, owner, repository, from, to)
				}
			},
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
func NewPullRequest(gitCommand *GitCommand) *PullRequest {
	return &PullRequest{
		GitServices: getServices(gitCommand.Config),
		GitCommand:  gitCommand,
	}
}

// Create opens link to new pull request in browser
func (pr *PullRequest) Create(from *models.Branch, to *models.Branch) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(from, to)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.GitCommand.OSCommand.OpenLink(pullRequestURL)
}

// CopyURL copies the pull request URL to the clipboard
func (pr *PullRequest) CopyURL(from *models.Branch, to *models.Branch) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(from, to)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.GitCommand.OSCommand.CopyToClipboard(pullRequestURL)
}

func (pr *PullRequest) getPullRequestURL(from *models.Branch, to *models.Branch) (string, error) {
	branchExistsOnRemote := pr.GitCommand.CheckRemoteBranchExists(from)

	if !branchExistsOnRemote {
		return "", errors.New(pr.GitCommand.Tr.NoBranchOnRemote)
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
		return "", errors.New(pr.GitCommand.Tr.UnsupportedGitService)
	}

	repoInfo := getRepoInfoFromURL(repoURL)
	toBranchName := ""
	if to != nil {
	        toBranchName = to.Name
	}
	pullRequestURL := gitService.PullRequestURL(repoInfo.Owner, repoInfo.Repository, from.Name, toBranchName)

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
