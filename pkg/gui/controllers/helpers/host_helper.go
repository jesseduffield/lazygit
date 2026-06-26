package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
)

// this helper just wraps our hosting_service package

type HostHelper struct {
	c *HelperCommon
}

func NewHostHelper(
	c *HelperCommon,
) *HostHelper {
	return &HostHelper{
		c: c,
	}
}

func (self *HostHelper) GetPullRequestURL(from string, to string) (string, error) {
	mgr, err := self.getHostingServiceMgr()
	if err != nil {
		return "", err
	}
	return mgr.GetPullRequestURL(from, to)
}

func (self *HostHelper) GetCommitURL(commitHash string) (string, error) {
	mgr, err := self.getHostingServiceMgr()
	if err != nil {
		return "", err
	}
	return mgr.GetCommitURL(commitHash)
}

// HasGithubRemote reports whether the repo has any GitHub remote. It's cheap
// (it only parses the remote URLs already in the model) and, unlike
// GithubBaseRemote, doesn't shell out to git or look up an auth token, so it's
// safe to call on the UI thread.
func (self *HostHelper) HasGithubRemote() bool {
	return len(getGithubRemotes(self.c)) > 0
}

// GithubBaseRemote resolves the GitHub remote that this repository's pull
// requests and deployments are made against, along with an auth token for its
// host. It returns ok == false when there is no GitHub remote or no token is
// available for it.
func (self *HostHelper) GithubBaseRemote() (hosting_service.ServiceInfo, string, bool) {
	githubRemotes := getAuthenticatedGithubRemotes(getGithubRemotes(self.c), self.c.Git().GitHub.GetAuthToken)
	baseRemote := getGithubBaseRemote(githubRemotes, self.c.Git().GitHub.ConfiguredBaseRemoteName())
	if baseRemote == nil {
		return hosting_service.ServiceInfo{}, "", false
	}
	return baseRemote.serviceInfo, baseRemote.authToken, true
}

// getting this on every request rather than storing it in state in case our remoteURL changes
// from one invocation to the next.
func (self *HostHelper) getHostingServiceMgr() (*hosting_service.HostingServiceMgr, error) {
	remoteUrl, err := self.c.Git().Remote.GetRemoteURL("origin")
	if err != nil {
		return nil, err
	}
	configServices := self.c.UserConfig().Services
	return hosting_service.NewHostingServiceMgr(self.c.Log, self.c.Tr, remoteUrl, configServices), nil
}
