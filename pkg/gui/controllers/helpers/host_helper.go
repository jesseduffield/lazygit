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
	return self.GetCommitURLForRemote(commitHash, "origin")
}

func (self *HostHelper) GetCommitURLForRemote(commitHash string, remoteName string) (string, error) {
	mgr, err := self.getHostingServiceMgrForRemote(remoteName)
	if err != nil {
		return "", err
	}
	return mgr.GetCommitURL(commitHash)
}

// getting this on every request rather than storing it in state in case our remoteURL changes
// from one invocation to the next.
func (self *HostHelper) getHostingServiceMgr() (*hosting_service.HostingServiceMgr, error) {
	return self.getHostingServiceMgrForRemote("origin")
}

func (self *HostHelper) getHostingServiceMgrForRemote(remoteName string) (*hosting_service.HostingServiceMgr, error) {
	remoteUrl, err := self.c.Git().Remote.GetRemoteURL(remoteName)
	if err != nil {
		return nil, err
	}
	configServices := self.c.UserConfig().Services
	return hosting_service.NewHostingServiceMgr(self.c.Log, self.c.Tr, remoteUrl, configServices), nil
}
