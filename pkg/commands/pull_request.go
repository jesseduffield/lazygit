package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// if you want to make a custom regex for a given service feel free to test it out
// at regoio.herokuapp.com
var defaultUrlRegexStrings = []string{
	`^(?:https?|ssh)://.*/(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
	`^git@.*:(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
}

// Service is a service that repository is on (Github, Bitbucket, ...)
type Service struct {
	Name                            string
	pullRequestURLIntoDefaultBranch func(owner string, repository string, from string) string
	pullRequestURLIntoTargetBranch  func(owner string, repository string, from string, to string) string
	URLRegexStrings                 []string
}

func NewGithubService(repositoryDomain string, siteDomain string) *Service {
	return &Service{
		Name: repositoryDomain,
		pullRequestURLIntoDefaultBranch: func(owner string, repository string, from string) string {
			return fmt.Sprintf("https://%s/%s/%s/compare/%s?expand=1", siteDomain, owner, repository, from)
		},
		pullRequestURLIntoTargetBranch: func(owner string, repository string, from string, to string) string {
			return fmt.Sprintf("https://%s/%s/%s/compare/%s...%s?expand=1", siteDomain, owner, repository, to, from)
		},
		URLRegexStrings: defaultUrlRegexStrings,
	}
}

func NewBitBucketService(repositoryDomain string, siteDomain string) *Service {
	return &Service{
		Name: repositoryDomain,
		pullRequestURLIntoDefaultBranch: func(owner string, repository string, from string) string {
			return fmt.Sprintf("https://%s/%s/%s/pull-requests/new?source=%s&t=1", siteDomain, owner, repository, from)
		},
		pullRequestURLIntoTargetBranch: func(owner string, repository string, from string, to string) string {
			return fmt.Sprintf("https://%s/%s/%s/pull-requests/new?source=%s&dest=%s&t=1", siteDomain, owner, repository, from, to)
		},
		URLRegexStrings: defaultUrlRegexStrings,
	}
}

func NewGitLabService(repositoryDomain string, siteDomain string) *Service {
	return &Service{
		Name: repositoryDomain,
		pullRequestURLIntoDefaultBranch: func(owner string, repository string, from string) string {
			return fmt.Sprintf("https://%s/%s/%s/merge_requests/new?merge_request[source_branch]=%s", siteDomain, owner, repository, from)
		},
		pullRequestURLIntoTargetBranch: func(owner string, repository string, from string, to string) string {
			return fmt.Sprintf("https://%s/%s/%s/merge_requests/new?merge_request[source_branch]=%s&merge_request[target_branch]=%s", siteDomain, owner, repository, from, to)
		},
		URLRegexStrings: defaultUrlRegexStrings,
	}
}

func (s *Service) PullRequestURL(repoURL string, from string, to string) string {
	repoInfo := s.getRepoInfoFromURL(repoURL)

	if to == "" {
		return s.pullRequestURLIntoDefaultBranch(repoInfo.Owner, repoInfo.Repository, from)
	} else {
		return s.pullRequestURLIntoTargetBranch(repoInfo.Owner, repoInfo.Repository, from, to)
	}
}

// PullRequest opens a link in browser to create new pull request
// with selected branch
type PullRequest struct {
	GitCommand *GitCommand
}

// RepoInformation holds some basic information about the repo
type RepoInformation struct {
	Owner      string
	Repository string
}

// NewPullRequest creates new instance of PullRequest
func NewPullRequest(gitCommand *GitCommand) *PullRequest {
	return &PullRequest{
		GitCommand: gitCommand,
	}
}

func (pr *PullRequest) getServices() []*Service {
	services := []*Service{
		NewGithubService("github.com", "github.com"),
		NewBitBucketService("bitbucket.org", "bitbucket.org"),
		NewGitLabService("gitlab.com", "gitlab.com"),
	}

	configServices := pr.GitCommand.Config.GetUserConfig().Services

	if len(configServices) > 0 {
		serviceFuncMap := map[string]func(repositoryDomain string, siteDomain string) *Service{
			"github":    NewGithubService,
			"bitbucket": NewBitBucketService,
			"gitlab":    NewGitLabService,
		}

		for repoDomain, typeAndDomain := range configServices {
			splitData := strings.Split(typeAndDomain, ":")
			if len(splitData) != 2 {
				pr.GitCommand.Log.Errorf("Unexpected format for git service: '%s'. Expected something like 'github.com:github.com'", typeAndDomain)
				continue
			}

			serviceFunc := serviceFuncMap[splitData[0]]
			if serviceFunc == nil {
				serviceNames := []string{}
				for serviceName := range serviceFuncMap {
					serviceNames = append(serviceNames, serviceName)
				}
				pr.GitCommand.Log.Errorf("Unknown git service type: '%s'. Expected one of %s", splitData[0], strings.Join(serviceNames, ", "))
				continue
			}

			services = append(services, serviceFunc(repoDomain, splitData[1]))
		}
	}

	return services
}

// Create opens link to new pull request in browser
func (pr *PullRequest) Create(from string, to string) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(from, to)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.GitCommand.OSCommand.OpenLink(pullRequestURL)
}

// CopyURL copies the pull request URL to the clipboard
func (pr *PullRequest) CopyURL(from string, to string) (string, error) {
	pullRequestURL, err := pr.getPullRequestURL(from, to)
	if err != nil {
		return "", err
	}

	return pullRequestURL, pr.GitCommand.OSCommand.CopyToClipboard(pullRequestURL)
}

func (pr *PullRequest) getPullRequestURL(from string, to string) (string, error) {
	branchExistsOnRemote := pr.GitCommand.CheckRemoteBranchExists(from)

	if !branchExistsOnRemote {
		return "", errors.New(pr.GitCommand.Tr.NoBranchOnRemote)
	}

	repoURL := pr.GitCommand.GetRemoteURL()
	var gitService *Service

	for _, service := range pr.getServices() {
		if strings.Contains(repoURL, service.Name) {
			gitService = service
			break
		}
	}

	if gitService == nil {
		return "", errors.New(pr.GitCommand.Tr.UnsupportedGitService)
	}

	pullRequestURL := gitService.PullRequestURL(repoURL, from, to)

	return pullRequestURL, nil
}

func (s *Service) getRepoInfoFromURL(url string) *RepoInformation {
	for _, regexStr := range s.URLRegexStrings {
		re := regexp.MustCompile(regexStr)
		matches := utils.FindNamedMatches(re, url)
		if matches != nil {
			return &RepoInformation{
				Owner:      matches["owner"],
				Repository: matches["repo"],
			}
		}
	}

	return nil
}
