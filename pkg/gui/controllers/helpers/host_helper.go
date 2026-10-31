package helpers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
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

// GetPullRequestURL builds a pull request URL for the given remote URL.
func (self *HostHelper) GetPullRequestURL(remoteURL string, from string, to string) (string, error) {
	return self.getHostingServiceMgr(remoteURL).GetPullRequestURL(from, to)
}

// WithPullRequestURL resolves which remote to use for hosting-service links, builds the
// pull request URL, then invokes onURL with it.
func (self *HostHelper) WithPullRequestURL(from string, to string, onURL func(url string) error) error {
	return self.withRemoteURL(func(remoteURL string) error {
		url, err := self.GetPullRequestURL(remoteURL, from, to)
		if err != nil {
			return err
		}
		return onURL(url)
	})
}

// WithCommitURL resolves which remote to use for hosting-service links, builds the
// commit URL, then invokes onURL with it.
func (self *HostHelper) WithCommitURL(commitHash string, onURL func(url string) error) error {
	return self.withRemoteURL(func(remoteURL string) error {
		url, err := self.getHostingServiceMgr(remoteURL).GetCommitURL(commitHash)
		if err != nil {
			return err
		}
		return onURL(url)
	})
}

// WithRemoteURL invokes onConfirm with a remote name and URL.
// If there is exactly one remote it is used immediately; if there are several the
// user is prompted to choose; if there are none an error is returned.
func (self *HostHelper) WithRemoteURL(onConfirm func(remoteName string, remoteURL string) error) error {
	return self.withRemote(onConfirm)
}

func (self *HostHelper) withRemoteURL(onConfirm func(remoteURL string) error) error {
	return self.withRemote(func(_ string, remoteURL string) error {
		return onConfirm(remoteURL)
	})
}

func (self *HostHelper) withRemote(onConfirm func(remoteName string, remoteURL string) error) error {
	remotes := self.c.Model().Remotes
	switch len(remotes) {
	case 0:
		return errors.New(self.c.Tr.NoRemotesFound)
	case 1:
		return self.invokeWithRemote(remotes[0], onConfirm)
	default:
		return self.promptForRemote(remotes, onConfirm)
	}
}

func (self *HostHelper) invokeWithRemote(remote *models.Remote, onConfirm func(remoteName string, remoteURL string) error) error {
	remoteURL, err := self.c.Git().Remote.GetRemoteURL(remote.Name)
	if err != nil {
		return err
	}
	return onConfirm(remote.Name, remoteURL)
}

func (self *HostHelper) promptForRemote(remotes []*models.Remote, onConfirm func(remoteName string, remoteURL string) error) error {
	menuItems := lo.Map(remotes, func(remote *models.Remote, _ int) *types.MenuItem {
		urlLabel := ""
		if len(remote.Urls) > 0 {
			urlLabel = style.FgCyan.Sprint(remote.Urls[0])
		}
		return &types.MenuItem{
			LabelColumns: []string{remote.Name, urlLabel},
			OnPress: func() error {
				return self.invokeWithRemote(remote, onConfirm)
			},
		}
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SelectRemote,
		Items: menuItems,
	})
}

func (self *HostHelper) getHostingServiceMgr(remoteURL string) *hosting_service.HostingServiceMgr {
	configServices := self.c.UserConfig().Services
	return hosting_service.NewHostingServiceMgr(self.c.Log, self.c.Tr, remoteURL, configServices)
}
