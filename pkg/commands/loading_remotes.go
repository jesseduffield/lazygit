package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (c *GitCommand) GetBranchesFromDir(dirPath string) ([]*Branch, error) {
	branches := []*Branch{}
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// it's possible that go-git is referencing a remote we don't have locally
			// in which case we'll just swallow this error
			c.Log.Warn(err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// it's a file: we need to get the path and work out the branch name from that
		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		branches = append(branches, &Branch{
			Name: strings.TrimPrefix(path, dirPath)[1:], // stripping prefix slash
			Hash: strings.TrimSpace(string(fileContents)),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *GitCommand) GetRemotes() ([]*Remote, error) {
	goGitRemotes, err := c.Repo.Remotes()
	if err != nil {
		return nil, err
	}

	// first step is to get our remotes from go-git
	remotes := make([]*Remote, len(goGitRemotes))
	for i, goGitRemote := range goGitRemotes {
		name := goGitRemote.Config().Name

		// can't seem to easily get the branches of the remotes from go-git so we'll
		// traverse the directory recursively
		branches, err := c.GetBranchesFromDir(filepath.Join(".git", "refs", "remotes", name))
		if err != nil {
			return nil, err
		}

		remotes[i] = &Remote{
			Name:     goGitRemote.Config().Name,
			Urls:     goGitRemote.Config().URLs,
			Branches: branches,
		}
	}

	return remotes, nil
}
