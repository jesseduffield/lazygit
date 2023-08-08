package hosting_service

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/fakes"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

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
			testName:  "Opens a link to new pull request on github with https remote url",
			from:      "feature/sum-operation",
			remoteUrl: "https://github.com/peter/calculator.git",
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
			testName:  "Opens a link to new pull request on github with specific target branch (different git username)",
			from:      "feature/sum-operation",
			to:        "feature/operations",
			remoteUrl: "ssh://org-12345@github.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://github.com/peter/calculator/compare/feature%2Foperations...feature%2Fsum-operation?expand=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on github with https remote url with specific target branch",
			from:      "feature/sum-operation",
			to:        "feature/operations",
			remoteUrl: "https://github.com/peter/calculator.git",
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
				assert.Equal(t, "https://gitlab.com/peter/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab in nested groups",
			from:      "feature/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with https remote url in nested groups",
			from:      "feature/ui",
			remoteUrl: "https://gitlab.com/peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fcommit-ui&merge_request[target_branch]=epic%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with specific target branch in nested groups",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "git@gitlab.com:peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fcommit-ui&merge_request[target_branch]=epic%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on gitlab with https remote url with specific target branch in nested groups",
			from:      "feature/commit-ui",
			to:        "epic/ui",
			remoteUrl: "https://gitlab.com/peter/public/calculator.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/peter/public/calculator/-/merge_requests/new?merge_request[source_branch]=feature%2Fcommit-ui&merge_request[target_branch]=epic%2Fui", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on bitbucket with a custom SSH username",
			from:      "feature/profile-page",
			remoteUrl: "john@bitbucket.org:johndoe/social_network.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://bitbucket.org/johndoe/social_network/pull-requests/new?source=feature%2Fprofile-page&t=1", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Azure DevOps (SSH)",
			from:      "feature/new",
			remoteUrl: "git@ssh.dev.azure.com:v3/myorg/myproject/myrepo",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://dev.azure.com/myorg/myproject/_git/myrepo/pullrequestcreate?sourceRef=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Azure DevOps (SSH) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "git@ssh.dev.azure.com:v3/myorg/myproject/myrepo",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://dev.azure.com/myorg/myproject/_git/myrepo/pullrequestcreate?sourceRef=feature%2Fnew&targetRef=dev", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Azure DevOps (HTTP)",
			from:      "feature/new",
			remoteUrl: "https://myorg@dev.azure.com/myorg/myproject/_git/myrepo",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://dev.azure.com/myorg/myproject/_git/myrepo/pullrequestcreate?sourceRef=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Azure DevOps (HTTP) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "https://myorg@dev.azure.com/myorg/myproject/_git/myrepo",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://dev.azure.com/myorg/myproject/_git/myrepo/pullrequestcreate?sourceRef=feature%2Fnew&targetRef=dev", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Bitbucket Server (SSH)",
			from:      "feature/new",
			remoteUrl: "ssh://git@mycompany.bitbucket.com/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a bitbucket server URL
				"mycompany.bitbucket.com": "bitbucketServer:mycompany.bitbucket.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.bitbucket.com/projects/myproject/repos/myrepo/pull-requests?create&sourceBranch=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Bitbucket Server (SSH) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "ssh://git@mycompany.bitbucket.com/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a bitbucket server URL
				"mycompany.bitbucket.com": "bitbucketServer:mycompany.bitbucket.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.bitbucket.com/projects/myproject/repos/myrepo/pull-requests?create&targetBranch=dev&sourceBranch=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Bitbucket Server (HTTP)",
			from:      "feature/new",
			remoteUrl: "https://mycompany.bitbucket.com/scm/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a bitbucket server URL
				"mycompany.bitbucket.com": "bitbucketServer:mycompany.bitbucket.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.bitbucket.com/projects/myproject/repos/myrepo/pull-requests?create&sourceBranch=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Bitbucket Server (HTTP) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "https://mycompany.bitbucket.com/scm/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a bitbucket server URL
				"mycompany.bitbucket.com": "bitbucketServer:mycompany.bitbucket.com",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.bitbucket.com/projects/myproject/repos/myrepo/pull-requests?create&targetBranch=dev&sourceBranch=feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Gitea Server (SSH)",
			from:      "feature/new",
			remoteUrl: "ssh://git@mycompany.gitea.io/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a gitea server URL
				"mycompany.gitea.io": "gitea:mycompany.gitea.io",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.gitea.io/myproject/myrepo/compare/feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Gitea Server (SSH) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "ssh://git@mycompany.gitea.io/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a gitea server URL
				"mycompany.gitea.io": "gitea:mycompany.gitea.io",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.gitea.io/myproject/myrepo/compare/dev...feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Gitea Server (HTTP)",
			from:      "feature/new",
			remoteUrl: "https://mycompany.gitea.io/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a gitea server URL
				"mycompany.gitea.io": "gitea:mycompany.gitea.io",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.gitea.io/myproject/myrepo/compare/feature%2Fnew", url)
			},
		},
		{
			testName:  "Opens a link to new pull request on Gitea Server (HTTP) with specific target",
			from:      "feature/new",
			to:        "dev",
			remoteUrl: "https://mycompany.gitea.io/myproject/myrepo.git",
			configServiceDomains: map[string]string{
				// valid configuration for a gitea server URL
				"mycompany.gitea.io": "gitea:mycompany.gitea.io",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://mycompany.gitea.io/myproject/myrepo/compare/dev...feature%2Fnew", url)
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
			testName:  "Does not log error when config service webDomain contains a port",
			from:      "feature/profile-page",
			remoteUrl: "git@my.domain.test:johndoe/social_network.git",
			configServiceDomains: map[string]string{
				"my.domain.test": "gitlab:my.domain.test:1111",
			},
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://my.domain.test:1111/johndoe/social_network/-/merge_requests/new?merge_request[source_branch]=feature%2Fprofile-page", url)
			},
		},
		{
			testName:  "Logs error when webDomain contains more than one colon",
			from:      "feature/profile-page",
			remoteUrl: "git@my.domain.test:johndoe/social_network.git",
			configServiceDomains: map[string]string{
				"my.domain.test": "gitlab:my.domain.test:1111:2222",
			},
			test: func(url string, err error) {
				assert.Error(t, err)
			},
			expectedLoggedErrors: []string{"Unexpected format for git service: 'gitlab:my.domain.test:1111:2222'. Expected something like 'github.com:github.com'"},
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
			expectedLoggedErrors: []string{"Unknown git service type: 'noservice'. Expected one of github, bitbucket, gitlab, azuredevops, bitbucketServer, gitea"},
		},
		{
			testName:  "Escapes reserved URL characters in from branch name",
			from:      "feature/someIssue#123",
			to:        "master",
			remoteUrl: "git@gitlab.com:me/public/repo-with-issues.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/me/public/repo-with-issues/-/merge_requests/new?merge_request[source_branch]=feature%2FsomeIssue%23123&merge_request[target_branch]=master", url)
			},
		},
		{
			testName:  "Escapes reserved URL characters in to branch name",
			from:      "yolo",
			to:        "archive/never-ending-feature#666",
			remoteUrl: "git@gitlab.com:me/public/repo-with-issues.git",
			test: func(url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://gitlab.com/me/public/repo-with-issues/-/merge_requests/new?merge_request[source_branch]=yolo&merge_request[target_branch]=archive%2Fnever-ending-feature%23666", url)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			tr := i18n.EnglishTranslationSet()
			log := &fakes.FakeFieldLogger{}
			hostingServiceMgr := NewHostingServiceMgr(log, &tr, s.remoteUrl, s.configServiceDomains)
			s.test(hostingServiceMgr.GetPullRequestURL(s.from, s.to))
			log.AssertErrors(t, s.expectedLoggedErrors)
		})
	}
}
