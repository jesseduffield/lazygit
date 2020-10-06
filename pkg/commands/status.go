package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (string, error) {
	exists, err := c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-apply"))
	if err != nil {
		return "", err
	}
	if exists {
		return "normal", nil
	}
	exists, err = c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-merge"))
	if exists {
		return "interactive", err
	} else {
		return "", err
	}
}

func (c *GitCommand) WorkingTreeState() string {
	rebaseMode, _ := c.RebaseMode()
	if rebaseMode != "" {
		return "rebasing"
	}
	merging, _ := c.IsInMergeState()
	if merging {
		return "merging"
	}
	return "normal"
}

// IsInMergeState states whether we are still mid-merge
func (c *GitCommand) IsInMergeState() (bool, error) {
	return c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "MERGE_HEAD"))
}

func (c *GitCommand) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := c.Repo.Worktree()
	return err == gogit.ErrIsBareRepository
}
