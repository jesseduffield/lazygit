package hosting_service

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"golang.org/x/exp/slices"
)

// This package is for handling logic specific to a git hosting service like github, gitlab, bitbucket, gitea, etc.
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
		return gitService.getPullRequestURLIntoDefaultBranch(url.QueryEscape(from)), nil
	} else {
		return gitService.getPullRequestURLIntoTargetBranch(url.QueryEscape(from), url.QueryEscape(to)), nil
	}
}

func (self *HostingServiceMgr) GetCommitURL(commitHash string) (string, error) {
	gitService, err := self.getService()
	if err != nil {
		return "", err
	}

	pullRequestURL := gitService.getCommitURL(commitHash)

	return pullRequestURL, nil
}

func (self *HostingServiceMgr) getService() (*Service, error) {
	serviceDomain, err := self.getServiceDomain(self.remoteURL)
	if err != nil {
		return nil, err
	}

	repoURL, err := serviceDomain.serviceDefinition.getRepoURLFromRemoteURL(self.remoteURL, serviceDomain.webDomain)
	if err != nil {
		return nil, err
	}

	return &Service{
		repoURL:           repoURL,
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

	serviceDomains := slices.Clone(defaultServiceDomains)

	for gitDomain, typeAndDomain := range self.configServiceDomains {
		provider, webDomain, success := strings.Cut(typeAndDomain, ":")

		// we allow for one ':' for specifying the TCP port
		if !success || strings.Count(webDomain, ":") > 1 {
			self.log.Errorf("Unexpected format for git service: '%s'. Expected something like 'github.com:github.com'", typeAndDomain)
			continue
		}

		serviceDefinition, ok := serviceDefinitionByProvider[provider]
		if !ok {
			providerNames := lo.Map(serviceDefinitions, func(serviceDefinition ServiceDefinition, _ int) string {
				return serviceDefinition.provider
			})

			self.log.Errorf("Unknown git service type: '%s'. Expected one of %s", provider, strings.Join(providerNames, ", "))
			continue
		}

		serviceDomains = append(serviceDomains, ServiceDomain{
			gitDomain:         gitDomain,
			webDomain:         webDomain,
			serviceDefinition: serviceDefinition,
		})
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

type ServiceDefinition struct {
	provider                        string
	pullRequestURLIntoDefaultBranch string
	pullRequestURLIntoTargetBranch  string
	commitURL                       string
	regexStrings                    []string

	// can expect 'webdomain' to be passed in. Otherwise, you get to pick what we match in the regex
	repoURLTemplate string
}

func (self ServiceDefinition) getRepoURLFromRemoteURL(url string, webDomain string) (string, error) {
	for _, regexStr := range self.regexStrings {
		re := regexp.MustCompile(regexStr)
		input := utils.FindNamedMatches(re, url)
		if input != nil {
			input["webDomain"] = webDomain
			return utils.ResolvePlaceholderString(self.repoURLTemplate, input), nil
		}
	}

	return "", errors.New("Failed to parse repo information from url")
}

type Service struct {
	repoURL string
	ServiceDefinition
}

func (self *Service) getPullRequestURLIntoDefaultBranch(from string) string {
	return self.resolveUrl(self.pullRequestURLIntoDefaultBranch, map[string]string{"From": from})
}

func (self *Service) getPullRequestURLIntoTargetBranch(from string, to string) string {
	return self.resolveUrl(self.pullRequestURLIntoTargetBranch, map[string]string{"From": from, "To": to})
}

func (self *Service) getCommitURL(commitHash string) string {
	return self.resolveUrl(self.commitURL, map[string]string{"CommitHash": commitHash})
}

func (self *Service) resolveUrl(templateString string, args map[string]string) string {
	return self.repoURL + utils.ResolvePlaceholderString(templateString, args)
}
