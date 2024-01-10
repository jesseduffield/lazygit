package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
)

// this helper just wraps our hosting_service package

type IHostHelper interface {
	GetPullRequestURL(from string, to string) (string, error)
	GetCommitURL(commitSha string) (string, error)
}

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

func (self *HostHelper) GetCommitURL(commitSha string) (string, error) {
	mgr, err := self.getHostingServiceMgr()
	if err != nil {
		return "", err
	}
	return mgr.GetCommitURL(commitSha)
}

// getting this on every request rather than storing it in state in case our remoteURL changes
// from one invocation to the next.
func (self *HostHelper) getHostingServiceMgr() (*hosting_service.HostingServiceMgr, error) {
	remoteUrl, err := self.c.Git().Remote.GetRemoteURL("origin")
	if err != nil {
		return nil, err
	}
	configServices := self.c.UserConfig.Services
	return hosting_service.NewHostingServiceMgr(self.c.Log, self.c.Tr, remoteUrl, configServices), nil
}
