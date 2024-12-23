package hosting_service

// if you want to make a custom regex for a given service feel free to test it out
// at regoio.herokuapp.com
var defaultUrlRegexStrings = []string{
	`^(?:https?|ssh)://[^/]+/(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
	`^.*?@.*:(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
}
var defaultRepoURLTemplate = "https://{{.webDomain}}/{{.owner}}/{{.repo}}"

// we've got less type safety using go templates but this lends itself better to
// users adding custom service definitions in their config
var githubServiceDef = ServiceDefinition{
	provider:                        "github",
	pullRequestURLIntoDefaultBranch: "/compare/{{.From}}?expand=1",
	pullRequestURLIntoTargetBranch:  "/compare/{{.To}}...{{.From}}?expand=1",
	commitURL:                       "/commit/{{.CommitHash}}",
	regexStrings:                    defaultUrlRegexStrings,
	repoURLTemplate:                 defaultRepoURLTemplate,
}

var bitbucketServiceDef = ServiceDefinition{
	provider:                        "bitbucket",
	pullRequestURLIntoDefaultBranch: "/pull-requests/new?source={{.From}}&t=1",
	pullRequestURLIntoTargetBranch:  "/pull-requests/new?source={{.From}}&dest={{.To}}&t=1",
	commitURL:                       "/commits/{{.CommitHash}}",
	regexStrings: []string{
		`^(?:https?|ssh)://.*/(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
		`^.*@.*:(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
	},
	repoURLTemplate: defaultRepoURLTemplate,
}

var gitLabServiceDef = ServiceDefinition{
	provider:                        "gitlab",
	pullRequestURLIntoDefaultBranch: "/-/merge_requests/new?merge_request[source_branch]={{.From}}",
	pullRequestURLIntoTargetBranch:  "/-/merge_requests/new?merge_request[source_branch]={{.From}}&merge_request[target_branch]={{.To}}",
	commitURL:                       "/-/commit/{{.CommitHash}}",
	regexStrings:                    defaultUrlRegexStrings,
	repoURLTemplate:                 defaultRepoURLTemplate,
}

var azdoServiceDef = ServiceDefinition{
	provider:                        "azuredevops",
	pullRequestURLIntoDefaultBranch: "/pullrequestcreate?sourceRef={{.From}}",
	pullRequestURLIntoTargetBranch:  "/pullrequestcreate?sourceRef={{.From}}&targetRef={{.To}}",
	commitURL:                       "/commit/{{.CommitHash}}",
	regexStrings: []string{
		`^git@ssh.dev.azure.com.*/(?P<org>.*)/(?P<project>.*)/(?P<repo>.*?)(?:\.git)?$`,
		`^https://.*@dev.azure.com/(?P<org>.*?)/(?P<project>.*?)/_git/(?P<repo>.*?)(?:\.git)?$`,
		`^https://.*/(?P<org>.*?)/(?P<project>.*?)/_git/(?P<repo>.*?)(?:\.git)?$`,
	},
	repoURLTemplate: "https://{{.webDomain}}/{{.org}}/{{.project}}/_git/{{.repo}}",
}

var bitbucketServerServiceDef = ServiceDefinition{
	provider:                        "bitbucketServer",
	pullRequestURLIntoDefaultBranch: "/pull-requests?create&sourceBranch={{.From}}",
	pullRequestURLIntoTargetBranch:  "/pull-requests?create&targetBranch={{.To}}&sourceBranch={{.From}}",
	commitURL:                       "/commits/{{.CommitHash}}",
	regexStrings: []string{
		`^ssh://git@.*/(?P<project>.*)/(?P<repo>.*?)(?:\.git)?$`,
		`^https://.*/scm/(?P<project>.*)/(?P<repo>.*?)(?:\.git)?$`,
	},
	repoURLTemplate: "https://{{.webDomain}}/projects/{{.project}}/repos/{{.repo}}",
}

var giteaServiceDef = ServiceDefinition{
	provider:                        "gitea",
	pullRequestURLIntoDefaultBranch: "/compare/{{.From}}",
	pullRequestURLIntoTargetBranch:  "/compare/{{.To}}...{{.From}}",
	commitURL:                       "/commit/{{.CommitHash}}",
	regexStrings:                    defaultUrlRegexStrings,
	repoURLTemplate:                 defaultRepoURLTemplate,
}

var serviceDefinitions = []ServiceDefinition{
	githubServiceDef,
	bitbucketServiceDef,
	gitLabServiceDef,
	azdoServiceDef,
	bitbucketServerServiceDef,
	giteaServiceDef,
}

var defaultServiceDomains = []ServiceDomain{
	{
		serviceDefinition: githubServiceDef,
		gitDomain:         "github.com",
		webDomain:         "github.com",
	},
	{
		serviceDefinition: bitbucketServiceDef,
		gitDomain:         "bitbucket.org",
		webDomain:         "bitbucket.org",
	},
	{
		serviceDefinition: gitLabServiceDef,
		gitDomain:         "gitlab.com",
		webDomain:         "gitlab.com",
	},
	{
		serviceDefinition: azdoServiceDef,
		gitDomain:         "dev.azure.com",
		webDomain:         "dev.azure.com",
	},
	{
		serviceDefinition: giteaServiceDef,
		gitDomain:         "try.gitea.io",
		webDomain:         "try.gitea.io",
	},
}
