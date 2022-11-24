package git_commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/generics/slices"
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type RemoteLoader struct {
	*common.Common
	cmd             oscommands.ICmdObjBuilder
	getGoGitRemotes func() ([]*gogit.Remote, error)
}

func NewRemoteLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
	getGoGitRemotes func() ([]*gogit.Remote, error),
) *RemoteLoader {
	return &RemoteLoader{
		Common:          common,
		cmd:             cmd,
		getGoGitRemotes: getGoGitRemotes,
	}
}

func (self *RemoteLoader) GetRemotes() ([]*models.Remote, error) {
	remoteBranchesStr, err := self.cmd.New("git branch -r").DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	goGitRemotes, err := self.getGoGitRemotes()
	if err != nil {
		return nil, err
	}

	// first step is to get our remotes from go-git
	remotes := slices.Map(goGitRemotes, func(goGitRemote *gogit.Remote) *models.Remote {
		remoteName := goGitRemote.Config().Name

		re := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s\/([\S]+)`, regexp.QuoteMeta(remoteName)))
		matches := re.FindAllStringSubmatch(remoteBranchesStr, -1)
		branches := slices.Map(matches, func(match []string) *models.RemoteBranch {
			return &models.RemoteBranch{
				Name:       match[1],
				RemoteName: remoteName,
			}
		})

		return &models.Remote{
			Name:     goGitRemote.Config().Name,
			Urls:     goGitRemote.Config().URLs,
			Branches: branches,
		}
	})

	// now lets sort our remotes by name alphabetically
	slices.SortFunc(remotes, func(a, b *models.Remote) bool {
		// we want origin at the top because we'll be most likely to want it
		if a.Name == "origin" {
			return true
		}
		if b.Name == "origin" {
			return false
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	return remotes, nil
}
