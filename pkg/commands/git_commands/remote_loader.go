package git_commands

import (
	"fmt"
	"strings"
	"sync"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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
	wg := sync.WaitGroup{}
	wg.Add(1)

	var remoteBranchesByRemoteName map[string][]*models.RemoteBranch
	var remoteBranchesErr error
	go utils.Safe(func() {
		defer wg.Done()

		remoteBranchesByRemoteName, remoteBranchesErr = self.getRemoteBranchesByRemoteName()
	})

	goGitRemotes, err := self.getGoGitRemotes()
	if err != nil {
		return nil, err
	}

	wg.Wait()

	if remoteBranchesErr != nil {
		return nil, remoteBranchesErr
	}

	remotes := lo.Map(goGitRemotes, func(goGitRemote *gogit.Remote, _ int) *models.Remote {
		remoteName := goGitRemote.Config().Name
		branches := remoteBranchesByRemoteName[remoteName]

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

func (self *RemoteLoader) getRemoteBranchesByRemoteName() (map[string][]*models.RemoteBranch, error) {
	remoteBranchesByRemoteName := make(map[string][]*models.RemoteBranch)

	var sortOrder string
	switch strings.ToLower(self.AppState.RemoteBranchSortOrder) {
	case "alphabetical":
		sortOrder = "refname"
	case "date":
		sortOrder = "-committerdate"
	default:
		sortOrder = "refname"
	}

	cmdArgs := NewGitCmd("for-each-ref").
		Arg(fmt.Sprintf("--sort=%s", sortOrder)).
		Arg("--format=%(refname:short)").
		Arg("refs/remotes").
		ToArgv()

	err := self.cmd.New(cmdArgs).DontLog().RunAndProcessLines(func(line string) (bool, error) {
		line = strings.TrimSpace(line)

		split := strings.SplitN(line, "/", 2)
		if len(split) != 2 {
			return false, nil
		}
		remoteName := split[0]
		name := split[1]

		_, ok := remoteBranchesByRemoteName[remoteName]
		if !ok {
			remoteBranchesByRemoteName[remoteName] = []*models.RemoteBranch{}
		}

		remoteBranchesByRemoteName[remoteName] = append(remoteBranchesByRemoteName[remoteName],
			&models.RemoteBranch{
				Name:       name,
				RemoteName: remoteName,
			})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return remoteBranchesByRemoteName, nil
}
