package git_commands

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/kofalt/go-memoize"
	"github.com/patrickmn/go-cache"
)

type RepoPathCache struct {
	cmd        oscommands.ICmdObjBuilder
	gitVersion *GitVersion
	cache      memoize.Memoizer
}

func NewRepoPathCache(cmd oscommands.ICmdObjBuilder, gitVersion *GitVersion) RepoPathCache {
	return RepoPathCache{
		cmd:        cmd,
		gitVersion: gitVersion,
		cache:      *memoize.NewMemoizer(cache.NoExpiration, cache.NoExpiration),
	}
}

func (self *RepoPathCache) GetRepoPathsForDir(dir string) (*RepoPaths, error) {
	getter := func() (interface{}, error) {
		return getRepoPaths(dir, self.cmd, self.gitVersion)
	}

	repoPaths, err, _ := self.cache.Memoize(dir, getter)
	if err != nil {
		return nil, err
	}

	return repoPaths.(*RepoPaths), nil
}

func (self *RepoPathCache) GetRepoPaths() (*RepoPaths, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return self.GetRepoPathsForDir(cwd)
}

func (self *RepoPathCache) GetGitVersion() *GitVersion {
	return self.gitVersion
}
