package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// this helper wraps our hosting_service package, and answers what we know about the
// repo's pull requests on it

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

// PullRequestForBranch returns the pull request the given local branch is the head of,
// and false where it has none. That covers knowing of no pull requests at all: the repo
// may not be on GitHub, or the answer may not have arrived yet.
func (self *HostHelper) PullRequestForBranch(branchName string) (*models.GithubPullRequest, bool) {
	pr, ok := self.c.Model().PullRequestsMap[branchName]
	return pr, ok
}

// NoPullRequestDisabledReason disables a command that acts on a branch's pull request
// while that branch has none.
func (self *HostHelper) NoPullRequestDisabledReason(branchName string) *types.DisabledReason {
	if _, ok := self.PullRequestForBranch(branchName); !ok {
		return &types.DisabledReason{Text: self.c.Tr.NoPullRequestForBranch, ShowErrorInPanel: true}
	}

	return nil
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
