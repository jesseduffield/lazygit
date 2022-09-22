package git_commands

import "github.com/jesseduffield/lazygit/pkg/commands/hosting_service"

// a hosting service is something like github, gitlab, bitbucket etc
type HostingService struct {
	*GitCommon
}

func NewHostingServiceCommand(gitCommon *GitCommon) *HostingService {
	return &HostingService{
		GitCommon: gitCommon,
	}
}

func (self *HostingService) GetPullRequestURL(from string, to string) (string, error) {
	return self.getHostingServiceMgr(self.config.GetRemoteURL()).GetPullRequestURL(from, to)
}

func (self *HostingService) GetCommitURL(commitSha string) (string, error) {
	return self.getHostingServiceMgr(self.config.GetRemoteURL()).GetCommitURL(commitSha)
}

func (self *HostingService) GetRepoNameFromRemoteURL(remoteURL string) (string, error) {
	return self.getHostingServiceMgr(remoteURL).GetRepoName()
}

// getting this on every request rather than storing it in state in case our remoteURL changes
// from one invocation to the next. Note however that we're currently caching config
// results so we might want to invalidate the cache here if it becomes a problem.
func (self *HostingService) getHostingServiceMgr(remoteURL string) *hosting_service.HostingServiceMgr {
	configServices := self.UserConfig.Services
	return hosting_service.NewHostingServiceMgr(self.Log, self.Tr, remoteURL, configServices)
}
