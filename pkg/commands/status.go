package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
)

const (
	REBASE_MODE_NORMAL      = "normal"
	REBASE_MODE_INTERACTIVE = "interactive"
	REBASE_MODE_REBASING    = "rebasing"
	REBASE_MODE_MERGING     = "merging"
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (string, error) {
	exists, err := c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-apply"))
	if err != nil {
		return "", err
	}
	if exists {
		return REBASE_MODE_NORMAL, nil
	}
	exists, err = c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-merge"))
	if exists {
		return REBASE_MODE_INTERACTIVE, err
	} else {
		return "", err
	}
}

func (c *GitCommand) WorkingTreeState() string {
	rebaseMode, _ := c.RebaseMode()
	if rebaseMode != "" {
		return REBASE_MODE_REBASING
	}
	merging, _ := c.IsInMergeState()
	if merging {
		return REBASE_MODE_MERGING
	}
	return REBASE_MODE_NORMAL
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
