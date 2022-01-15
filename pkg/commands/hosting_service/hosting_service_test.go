package hosting_service

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoInfoFromURL(t *testing.T) {
	type scenario struct {
		serviceDefinition ServiceDefinition
		testName          string
		repoURL           string
		test              func(*RepoInformation)
	}

	scenarios := []scenario{
		{
			githubServiceDef,
			"Returns repository information for git remote url",
			"git@github.com:petersmith/super_calculator",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "petersmith")
				assert.EqualValues(t, repoInfo.Repository, "super_calculator")
			},
		},
		{
			githubServiceDef,
			"Returns repository information for git remote url, trimming trailing '.git'",
			"git@github.com:petersmith/super_calculator.git",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "petersmith")
				assert.EqualValues(t, repoInfo.Repository, "super_calculator")
			},
		},
		{
			githubServiceDef,
			"Returns repository information for ssh remote url",
			"ssh://git@github.com/petersmith/super_calculator",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "petersmith")
				assert.EqualValues(t, repoInfo.Repository, "super_calculator")
			},
		},
		{
			githubServiceDef,
			"Returns repository information for http remote url",
			"https://my_username@bitbucket.org/johndoe/social_network.git",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "johndoe")
				assert.EqualValues(t, repoInfo.Repository, "social_network")
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			result, err := s.serviceDefinition.getRepoInfoFromURL(s.repoURL)
			assert.NoError(t, err)
			s.test(result)
		})
	}
}

func TestGetPullRequestURL(t *testing.T) {
	type scenario struct {
		testName             string
		from                 string
		to                   string
		remoteUrl            string
		configServiceDomains map[string]string
		test                 func(url string, err error)
		expectedLoggedErrors []string
	}

	scenarios := []scenario{
		{
			testName:  "Opens a link to new pull request on bitbucket",
			from:      "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with http remote url",
			from:      "feature/events",
			remoteUrl: "https://my_username@bitbucket.org/johndoe/social_network.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fevents&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on github",
			from:      "feature/sum-operation",
			remoteUrl: "git@github.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature%2Fsum-operation?expand=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with specific target branch",
			from:      "feature/profile-page/avatar",
			to:        "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page%2Favatar&dest=feature%2Fprofile-page&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with http remote url with specified target branch",
			from:      "feature/remote-events",
			to:        "feature/events",
			remoteUrl: "https://my_username@bitbucket.org/johndoe/social_network.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fremote-events&dest=feature%2Fevents&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on github with specific target branch",
			from:      "feature/sum-operation",
			to:        "feature/operations",
			remoteUrl: "git@github.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature%2Foperations...feature%2Fsum-operation?expand=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab",
			from:      "feature/ui",
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab in nested groups",
			from:      "feature/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/merge_requests/new?merge_request[source_branch]=feature%2Fcommit-ui&merge_request[target_branch]=epic%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch in nested groups",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/merge_requests/new?merge_request[source_branch]=feature%2Fcommit-ui&merge_request[target_branch]=epic%2Fui", url)
			},
		},
		{
			testName:  "Throws an error if git service is unsupported",
			from:      "feature/divide-operation",
			remoteUrl: "git@something.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.EqualError(t, err, "Unsupported git service")
			},
		},
		{
			testName:  "Does not log error when config service domains are valid",
			from:      "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			configServiceDomains: map[string]string{
				// valid configuration for a custom service URL
				"git.work.com": "gitlab:code.work.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page&t=1", url)
			},
			expectedLoggedErrors: nil,
		},
		{
			testName:  "Logs error when config service domain is malformed",
			from:      "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			configServiceDomains: map[string]string{
				"noservice.work.com": "noservice.work.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page&t=1", url)
			},
			expectedLoggedErrors: []string{"Unexpected format for git service: 'noservice.work.com'. Expected something like 'github.com:github.com'"},
		},
		{
			testName:  "Logs error when config service domain uses unknown provider",
			from:      "feature/profile-page",
			remoteUrl: "git@bitbucket.org:johndoe/social_network.git",
			configServiceDomains: map[string]string{
				"invalid.work.com": "noservice:invalid.work.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page&t=1", url)
			},
			expectedLoggedErrors: []string{"Unknown git service type: 'noservice'. Expected one of github, bitbucket, gitlab"},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			tr := i18n.EnglishTranslationSet()
			log := &test.FakeFieldLogger{}
			hostingServiceMgr := NewHostingServiceMgr(log, &tr, s.remoteUrl, s.configServiceDomains)
			s.test(hostingServiceMgr.GetPullRequestURL(s.from, s.to))
			log.AssertErrors(t, s.expectedLoggedErrors)
		})
	}
}
