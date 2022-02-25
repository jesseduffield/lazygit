package hosting_service

// if you want to make a custom regex for a given service feel free to test it out
// at regoio.herokuapp.com
var defaultUrlRegexStrings = []string{
	`^(?:https?|ssh)://.*/(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
	`^git@.*:(?P<owner>.*)/(?P<repo>.*?)(?:\.git)?$`,
}

// we've got less type safety using go templates but this lends itself better to
// users adding custom service definitions in their config
var githubServiceDef = ServiceDefinition{
	provider:                        "github",
	pullRequestURLIntoDefaultBranch: "/compare/{{.From}}?expand=1",
	pullRequestURLIntoTargetBranch:  "/compare/{{.To}}...{{.From}}?expand=1",
	commitURL:                       "/commit/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var bitbucketServiceDef = ServiceDefinition{
	provider:                        "bitbucket",
	pullRequestURLIntoDefaultBranch: "/pull-requests/new?source={{.From}}&t=1",
	pullRequestURLIntoTargetBranch:  "/pull-requests/new?source={{.From}}&dest={{.To}}&t=1",
	commitURL:                       "/commits/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var gitLabServiceDef = ServiceDefinition{
	provider:                        "gitlab",
	pullRequestURLIntoDefaultBranch: "/merge_requests/new?merge_request[source_branch]={{.From}}",
	pullRequestURLIntoTargetBranch:  "/merge_requests/new?merge_request[source_branch]={{.From}}&merge_request[target_branch]={{.To}}",
	commitURL:                       "/commit/{{.CommitSha}}",
	regexStrings:                    defaultUrlRegexStrings,
}

var serviceDefinitions = []ServiceDefinition{githubServiceDef, bitbucketServiceDef, gitLabServiceDef}

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
}
