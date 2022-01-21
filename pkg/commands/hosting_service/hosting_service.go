package hosting_service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

// This package is for handling logic specific to a git hosting service like github, gitlab, bitbucket, etc.
// Different git hosting services have different URL formats for when you want to open a PR or view a commit,
// and this package's responsibility is to determine which service you're using based on the remote URL,
// and then which URL you need for whatever use case you have.

type HostingServiceMgr struct {
	log       logrus.FieldLogger
	tr        *i18n.TranslationSet
	remoteURL string // e.g. https://github.com/jesseduffield/lazygit

	// see https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-pull-request-urls
	configServiceDomains map[string]string
}

// NewHostingServiceMgr creates new instance of PullRequest
func NewHostingServiceMgr(log logrus.FieldLogger, tr *i18n.TranslationSet, remoteURL string, configServiceDomains map[string]string) *HostingServiceMgr {
	return &HostingServiceMgr{
		log:                  log,
		tr:                   tr,
		remoteURL:            remoteURL,
		configServiceDomains: configServiceDomains,
	}
}

func (self *HostingServiceMgr) GetPullRequestURL(from string, to string) (string, error) {
	gitService, err := self.getService()
	if err != nil {
		return "", err
	}

	if to == "" {
		return gitService.getPullRequestURLIntoDefaultBranch(from), nil
	} else {
		return gitService.getPullRequestURLIntoTargetBranch(from, to), nil
	}
}

func (self *HostingServiceMgr) GetCommitURL(commitSha string) (string, error) {
	gitService, err := self.getService()
	if err != nil {
		return "", err
	}

	pullRequestURL := gitService.getCommitURL(commitSha)

	return pullRequestURL, nil
}

func (self *HostingServiceMgr) getService() (*Service, error) {
	serviceDomain, err := self.getServiceDomain(self.remoteURL)
	if err != nil {
		return nil, err
	}

	root, err := serviceDomain.getRootFromRemoteURL(self.remoteURL)
	if err != nil {
		return nil, err
	}

	return &Service{
		root:              root,
		ServiceDefinition: serviceDomain.serviceDefinition,
	}, nil
}

func (self *HostingServiceMgr) getServiceDomain(repoURL string) (*ServiceDomain, error) {
	candidateServiceDomains := self.getCandidateServiceDomains()

	for _, serviceDomain := range candidateServiceDomains {
		if strings.Contains(repoURL, serviceDomain.gitDomain) {
			return &serviceDomain, nil
		}
	}

	return nil, errors.New(self.tr.UnsupportedGitService)
}

func (self *HostingServiceMgr) getCandidateServiceDomains() []ServiceDomain {
	serviceDefinitionByProvider := map[string]ServiceDefinition{}
	for _, serviceDefinition := range serviceDefinitions {
		serviceDefinitionByProvider[serviceDefinition.provider] = serviceDefinition
	}

	var serviceDomains = make([]ServiceDomain, len(defaultServiceDomains))
	copy(serviceDomains, defaultServiceDomains)

	if len(self.configServiceDomains) > 0 {
		for gitDomain, typeAndDomain := range self.configServiceDomains {
			splitData := strings.Split(typeAndDomain, ":")
			if len(splitData) != 2 {
				self.log.Errorf("Unexpected format for git service: '%s'. Expected something like 'github.com:github.com'", typeAndDomain)
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
				self.log.Errorf("Unknown git service type: '%s'. Expected one of %s", provider, strings.Join(providerNames, ", "))
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

// a service domains pairs a service definition with the actual domain it's being served from.
// Sometimes the git service is hosted in a custom domains so although it'll use say
// the github service definition, it'll actually be served from e.g. my-custom-github.com
type ServiceDomain struct {
	gitDomain         string // the one that appears in the git remote url
	webDomain         string // the one that appears in the web url
	serviceDefinition ServiceDefinition
}

func (self ServiceDomain) getRootFromRemoteURL(repoURL string) (string, error) {
	// we may want to make this more specific to the service in future e.g. if
	// some new service comes along which has a different root url structure.
	repoInfo, err := self.serviceDefinition.getRepoInfoFromURL(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s/%s/%s", self.webDomain, repoInfo.Owner, repoInfo.Repository), nil
}

// RepoInformation holds some basic information about the repo
type RepoInformation struct {
	Owner      string
	Repository string
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
