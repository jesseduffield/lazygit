package commands

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

//counterfeiter:generate . IRemotesMgr
type IRemotesMgr interface {
	Add(name string, url string) error
	Remove(name string) error
	Rename(oldRemoteName string, newRemoteName string) error
	UpdateUrl(remoteName string, updatedUrl string) error
	RemoteBranchExists(branch *models.Branch) bool
	GetCurrentRemoteUrl() string
	GetRemotes() ([]*models.Remote, error)
}

type RemotesMgr struct {
	ICommander

	config IGitConfigMgr
	repo   *gogit.Repository
}

func NewRemotesMgr(
	commander ICommander,
	config IGitConfigMgr,
	repo *gogit.Repository,
) *RemotesMgr {
	return &RemotesMgr{
		ICommander: commander,
		config:     config,
		repo:       repo,
	}
}

func (c *RemotesMgr) Add(name string, url string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("remote add %s %s", name, url))
}

func (c *RemotesMgr) Remove(name string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("remote remove %s", name))
}

func (c *RemotesMgr) Rename(oldRemoteName string, newRemoteName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("remote rename %s %s", oldRemoteName, newRemoteName))
}

func (c *RemotesMgr) UpdateUrl(remoteName string, updatedUrl string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("remote set-url %s %s", remoteName, updatedUrl))
}

// CheckRemoteBranchExists Returns remote branch
func (c *RemotesMgr) RemoteBranchExists(branch *models.Branch) bool {
	_, err := c.RunWithOutput(
		BuildGitCmdObjFromStr(
			fmt.Sprintf("show-ref --verify -- refs/remotes/origin/%s",
				branch.Name),
		),
	)

	return err == nil
}

// GetRemoteURL returns current repo remote url
func (c *RemotesMgr) GetCurrentRemoteUrl() string {
	return c.config.GetConfigValue("remote.origin.url")
}

func (c *RemotesMgr) GetRemotes() ([]*models.Remote, error) {
	// get remote branches
	remoteBranchesStr, err := c.RunWithOutput(
		BuildGitCmdObjFromStr("branch -r"),
	)
	if err != nil {
		return nil, err
	}

	goGitRemotes, err := c.repo.Remotes()
	if err != nil {
		return nil, err
	}

	// first step is to get our remotes from go-git
	remotes := make([]*models.Remote, len(goGitRemotes))
	for i, goGitRemote := range goGitRemotes {
		remoteName := goGitRemote.Config().Name

		re := regexp.MustCompile(fmt.Sprintf(`%s\/([\S]+)`, remoteName))
		matches := re.FindAllStringSubmatch(remoteBranchesStr, -1)
		branches := make([]*models.RemoteBranch, len(matches))
		for j, match := range matches {
			branches[j] = &models.RemoteBranch{
				Name:       match[1],
				RemoteName: remoteName,
			}
		}

		remotes[i] = &models.Remote{
			Name:     goGitRemote.Config().Name,
			Urls:     goGitRemote.Config().URLs,
			Branches: branches,
		}
	}

	// now lets sort our remotes by name alphabetically
	sort.Slice(remotes, func(i, j int) bool {
		// we want origin at the top because we'll be most likely to want it
		if remotes[i].Name == "origin" {
			return true
		}
		if remotes[j].Name == "origin" {
			return false
		}
		return strings.ToLower(remotes[i].Name) < strings.ToLower(remotes[j].Name)
	})

	return remotes, nil
}
