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

type ServiceDefinition struct {
	provider                        string
	pullRequestURLIntoDefaultBranch string
	pullRequestURLIntoTargetBranch  string
	commitURL                       string
	regexStrings                    []string
}

func (self ServiceDefinition) getRepoInfoFromURL(url string) (*RepoInformation, error) {
	for _, regexStr := range self.regexStrings {
		re := regexp.MustCompile(regexStr)
		matches := utils.FindNamedMatches(re, url)
		if matches != nil {
			return &RepoInformation{
				Owner:      matches["owner"],
				Repository: matches["repo"],
			}, nil
		}
	}

	return nil, errors.New("Failed to parse repo information from url")
}

// a service domains pairs a service definition with the actual domain it's being served from.
// Sometimes the git service is hosted in a custom domains so although it'll use say
// the github service definition, it'll actually be served from e.g. my-custom-github.com
type ServiceDomain struct {
	gitDomain         string // the one that appears in the git remote url
	webDomain         string // the one that appears in the web url
	serviceDefinition ServiceDefinition
}

func (self ServiceDomain) getRootFromRepoURL(repoURL string) (string, error) {
	// we may want to make this more specific to the service in future e.g. if
	// some new service comes along which has a different root url structure.
	repoInfo, err := self.serviceDefinition.getRepoInfoFromURL(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s/%s/%s", self.webDomain, repoInfo.Owner, repoInfo.Repository), nil
}

// we've got less type safety using go templates but this lends itself better to
// users adding custom service definitions in their config
var GithubServiceDef = ServiceDefinition{
	provider:                        "github",
	pullRequestURLIntoDefaultBranch: "/compare/{{.From}}?expand=1",
	pullRequestURLIntoTargetBranch:  "/compare/{{.To}}...{{.From}}?expand=1",
	commitURL:                       "/commit/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var BitbucketServiceDef = ServiceDefinition{
	provider:                        "bitbucket",
	pullRequestURLIntoDefaultBranch: "/pull-requests/new?source={{.From}}&t=1",
	pullRequestURLIntoTargetBranch:  "/pull-requests/new?source={{.From}}&dest={{.To}}&t=1",
	commitURL:                       "/commits/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var GitLabServiceDef = ServiceDefinition{
	provider:                        "gitlab",
	pullRequestURLIntoDefaultBranch: "/merge_requests/new?merge_request[source_branch]={{.From}}",
	pullRequestURLIntoTargetBranch:  "/merge_requests/new?merge_request[source_branch]={{.From}}&merge_request[target_branch]={{.To}}",
	commitURL:                       "/commit/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var serviceDefinitions = []ServiceDefinition{GithubServiceDef, BitbucketServiceDef, GitLabServiceDef}
var defaultServiceDomains = []ServiceDomain{
	{
		serviceDefinition: GithubServiceDef,
		gitDomain:         "github.com",
		webDomain:         "github.com",
	},
	{
		serviceDefinition: BitbucketServiceDef,
		gitDomain:         "bitbucket.org",
		webDomain:         "bitbucket.org",
	},
	{
		serviceDefinition: GitLabServiceDef,
		gitDomain:         "gitlab.com",
		webDomain:         "gitlab.com",
	},
}

type Service struct {
	root string
	ServiceDefinition
}

func (self *Service) getPullRequestURLIntoDefaultBranch(from string) string {
	return self.resolveUrl(self.pullRequestURLIntoDefaultBranch, map[string]string{"From": from})
}

func (self *Service) getPullRequestURLIntoTargetBranch(from string, to string) string {
	return self.resolveUrl(self.pullRequestURLIntoTargetBranch, map[string]string{"From": from, "To": to})
}

func (self *Service) getCommitURL(commitSha string) string {
	return self.resolveUrl(self.commitURL, map[string]string{"CommitSha": commitSha})
}

func (self *Service) resolveUrl(templateString string, args map[string]string) string {
	return self.root + utils.ResolvePlaceholderString(templateString, args)
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

func (pr *PullRequest) getService() (*Service, error) {
	serviceDomain, err := pr.getServiceDomain()
	if err != nil {
		return nil, err
	}

	repoURL := pr.GitCommand.GetRemoteURL()

	root, err := serviceDomain.getRootFromRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	return &Service{
		root:              root,
		ServiceDefinition: serviceDomain.serviceDefinition,
	}, nil
}

func (pr *PullRequest) getServiceDomain() (*ServiceDomain, error) {
	candidateServiceDomains := pr.getCandidateServiceDomains()

	repoURL := pr.GitCommand.GetRemoteURL()

	for _, serviceDomain := range candidateServiceDomains {
		// I feel like it makes more sense to see if the repo url contains the service domain's git domain,
		// but I don't want to break anything by changing that right now.
		if strings.Contains(repoURL, serviceDomain.serviceDefinition.provider) {
			return &serviceDomain, nil
		}
	}

	return nil, errors.New(pr.GitCommand.Tr.UnsupportedGitService)
}

func (pr *PullRequest) getCandidateServiceDomains() []ServiceDomain {
	serviceDefinitionByProvider := map[string]ServiceDefinition{}
	for _, serviceDefinition := range serviceDefinitions {
		serviceDefinitionByProvider[serviceDefinition.provider] = serviceDefinition
	}

	var serviceDomains = make([]ServiceDomain, len(defaultServiceDomains))
	copy(serviceDomains, defaultServiceDomains)

	// see https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-pull-request-urls
	configServices := pr.GitCommand.Config.GetUserConfig().Services
	if len(configServices) > 0 {
		for gitDomain, typeAndDomain := range configServices {
			splitData := strings.Split(typeAndDomain, ":")
			if len(splitData) != 2 {
				pr.GitCommand.Log.Errorf("Unexpected format for git service: '%s'. Expected something like 'github.com:github.com'", typeAndDomain)
				continue
			}

			provider := splitData[0]
			webDomain := splitData[1]

			serviceDefinition, ok := serviceDefinitionByProvider[provider]
			if !ok {
				providerNames := []string{}
				for _, serviceDefinition := range serviceDefinitions {
					providerNames = append(providerNames, serviceDefinition.provider)
				}
				pr.GitCommand.Log.Errorf("Unknown git service type: '%s'. Expected one of %s", provider, strings.Join(providerNames, ", "))
				continue
			}

			serviceDomains = append(serviceDomains, ServiceDomain{
				gitDomain:         gitDomain,
				webDomain:         webDomain,
				serviceDefinition: serviceDefinition,
			})
		}
	}

	return serviceDomains
}

// CreatePullRequest opens link to new pull request in browser
func (pr *PullRequest) CreatePullRequest(from string, to string) (string, error) {
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

	gitService, err := pr.getService()
	if err != nil {
		return "", err
	}

	if to == "" {
		return gitService.getPullRequestURLIntoDefaultBranch(from), nil
	} else {
		return gitService.getPullRequestURLIntoTargetBranch(from, to), nil
	}
}

func (pr *PullRequest) getCommitURL(commitSha string) (string, error) {
	gitService, err := pr.getService()
	if err != nil {
		return "", err
	}

	pullRequestURL := gitService.getCommitURL(commitSha)

	return pullRequestURL, nil
}

func (pr *PullRequest) OpenCommitInBrowser(commitSha string) (string, error) {
	url, err := pr.getCommitURL(commitSha)
	if err != nil {
		return "", err
	}

	return url, pr.GitCommand.OSCommand.OpenLink(url)
}
