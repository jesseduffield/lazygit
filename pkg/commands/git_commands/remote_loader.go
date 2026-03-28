package git_commands

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemoteLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewRemoteLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *RemoteLoader {
	return &RemoteLoader{
		Common: common,
		cmd:    cmd,
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

	remotes := self.getRemotesFromConfig()

	wg.Wait()

	if remoteBranchesErr != nil {
		return nil, remoteBranchesErr
	}

	for _, remote := range remotes {
		remote.Branches = remoteBranchesByRemoteName[remote.Name]
	}

	// now lets sort our remotes by name alphabetically
	slices.SortFunc(remotes, func(a, b *models.Remote) int {
		// we want origin at the top because we'll be most likely to want it
		if a.Name == "origin" {
			return -1
		}
		if b.Name == "origin" {
			return 1
		}
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return remotes, nil
}

func (self *RemoteLoader) getRemotesFromConfig() []*models.Remote {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--get-regexp", `^remote\.[^.]+\.url$`).ToArgv()
	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		// exit code 1 means no matching keys (no remotes configured)
		return nil
	}

	remotesByName := make(map[string]*models.Remote)

	for _, line := range strings.Split(output, "\n") {
		key, url, found := strings.Cut(strings.TrimSpace(line), " ")
		if !found {
			continue
		}
		// key is "remote.<name>.url"; strip prefix and suffix to get the name
		remoteName := strings.TrimSuffix(strings.TrimPrefix(key, "remote."), ".url")
		if _, ok := remotesByName[remoteName]; !ok {
			remotesByName[remoteName] = &models.Remote{Name: remoteName}
		}
		remotesByName[remoteName].Urls = append(remotesByName[remoteName].Urls, url)
	}

	return slices.Collect(maps.Values(remotesByName))
}

func (self *RemoteLoader) getRemoteBranchesByRemoteName() (map[string][]*models.RemoteBranch, error) {
	remoteBranchesByRemoteName := make(map[string][]*models.RemoteBranch)

	var sortOrder string
	switch strings.ToLower(self.UserConfig().Git.RemoteBranchSortOrder) {
	case "alphabetical":
		sortOrder = "refname"
	case "date":
		sortOrder = "-committerdate"
	default:
		sortOrder = "refname"
	}

	cmdArgs := NewGitCmd("for-each-ref").
		Arg(fmt.Sprintf("--sort=%s", sortOrder)).
		Arg("--format=%(refname)").
		Arg("refs/remotes").
		ToArgv()

	err := self.cmd.New(cmdArgs).DontLog().RunAndProcessLines(func(line string) (bool, error) {
		line = strings.TrimSpace(line)

		split := strings.SplitN(line, "/", 4)
		if len(split) != 4 {
			return false, nil
		}
		remoteName := split[2]
		name := split[3]

		if name == "HEAD" {
			return false, nil
		}

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
