package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// this helper just wraps our hosting_service package

type IHostHelper interface {
	GetPullRequestURL(from string, to string) (string, error)
	GetCommitURL(commitSha string) (string, error)
}

type HostHelper struct {
	c   *types.HelperCommon
	git *commands.GitCommand
}

func NewHostHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
) *HostHelper {
	return &HostHelper{
		c:   c,
		git: git,
	}
}

func (self *HostHelper) GetPullRequestURL(from string, to string) (string, error) {
	return self.getHostingServiceMgr().GetPullRequestURL(from, to)
}

func (self *HostHelper) GetCommitURL(commitSha string) (string, error) {
	return self.getHostingServiceMgr().GetCommitURL(commitSha)
}

// getting this on every request rather than storing it in state in case our remoteURL changes
// from one invocation to the next. Note however that we're currently caching config
// results so we might want to invalidate the cache here if it becomes a problem.
func (self *HostHelper) getHostingServiceMgr() *hosting_service.HostingServiceMgr {
	remoteUrl := self.git.Config.GetRemoteURL()
	configServices := self.c.UserConfig.Services
	return hosting_service.NewHostingServiceMgr(self.c.Log, self.c.Tr, remoteUrl, configServices)
}
