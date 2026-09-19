package git_commands

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type RemoteLoader struct {
	*GitCommon
}

func NewRemoteLoader(gitCommon *GitCommon) *RemoteLoader {
	return &RemoteLoader{GitCommon: gitCommon}
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
		Arg("--local", "--get-regexp", `^remote\.[^.]+\.(url|pushurl)$`).ToArgv()
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
		// key is "remote.<name>.url" or "remote.<name>.pushurl";
		// strip prefix and suffix to get the name
		rest, ok := strings.CutPrefix(key, "remote.")
		if !ok {
			continue
		}
		var remoteName string
		var isPushUrl bool
		if name, ok := strings.CutSuffix(rest, ".pushurl"); ok {
			remoteName, isPushUrl = name, true
		} else if name, ok := strings.CutSuffix(rest, ".url"); ok {
			remoteName, isPushUrl = name, false
		} else {
			continue
		}
		if _, ok := remotesByName[remoteName]; !ok {
			remotesByName[remoteName] = &models.Remote{Name: remoteName}
		}
		if isPushUrl {
			remotesByName[remoteName].PushUrls = append(remotesByName[remoteName].PushUrls, url)
		} else {
			remotesByName[remoteName].Urls = append(remotesByName[remoteName].Urls, url)
		}
	}

	return slices.Collect(maps.Values(remotesByName))
}

func (self *RemoteLoader) getRemoteBranchesByRemoteName() (map[string][]*models.RemoteBranch, error) {
	remoteBranches, err := self.getRemoteBranches()
	if err != nil {
		return nil, err
	}

	return lo.GroupBy(remoteBranches, func(branch *models.RemoteBranch) string {
		return branch.RemoteName
	}), nil
}

// Returns all remote branches, sorted the way the config asks for
func (self *RemoteLoader) getRemoteBranches() ([]*models.RemoteBranch, error) {
	sortByDate := strings.ToLower(self.UserConfig().Git.RemoteBranchSortOrder) == "date"
	sortOrder := "refname"
	if sortByDate {
		sortOrder = "-committerdate"
	}

	// Asking for the tip of a branch makes git read its commit, so only do it
	// when we are going to sort by ancestry below
	format := "%(refname)"
	if sortByDate {
		format += "%00%(objectname)%00%(committerdate:unix)"
	}

	cmdArgs := NewGitCmd("for-each-ref").
		Arg(fmt.Sprintf("--sort=%s", sortOrder)).
		Arg(fmt.Sprintf("--format=%s", format)).
		Arg("refs/remotes").
		ToArgv()

	remoteBranches := []*models.RemoteBranch{}
	tips := map[string]refTip{}
	err := self.cmd.New(cmdArgs).DontLog().RunAndProcessLines(func(line string) (bool, error) {
		fields := strings.Split(strings.TrimSpace(line), "\x00")
		refName := fields[0]

		split := strings.SplitN(refName, "/", 4)
		if len(split) != 4 {
			return false, nil
		}
		remoteName := split[2]
		name := split[3]

		if name == "HEAD" {
			return false, nil
		}

		remoteBranches = append(remoteBranches,
			&models.RemoteBranch{
				Name:       name,
				RemoteName: remoteName,
			})
		if len(fields) == 3 {
			tips[refName] = refTip{hash: fields[1], committerDate: fields[2]}
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if sortByDate {
		if err := sortRefsWithEqualDatesByAncestry(
			self.cmd, self.version, remoteBranches, (*models.RemoteBranch).FullRefName, tips,
		); err != nil {
			self.Log.Errorf("Failed to sort remote branches by ancestry: %v", err)
		}
	}

	return remoteBranches, nil
}
