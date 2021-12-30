package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
)

type RebaseMode int

const (
	// this means we're neither rebasing nor merging
	REBASE_MODE_NONE RebaseMode = iota
	// this means normal rebase as opposed to interactive rebase
	REBASE_MODE_NORMAL
	REBASE_MODE_INTERACTIVE
	// REBASE_MODE_REBASING is a general state that captures both REBASE_MODE_NORMAL and REBASE_MODE_INTERACTIVE
	REBASE_MODE_REBASING
	REBASE_MODE_MERGING
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (RebaseMode, error) {
	exists, err := c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-apply"))
	if err != nil {
		return REBASE_MODE_NONE, err
	}
	if exists {
		return REBASE_MODE_NORMAL, nil
	}
	exists, err = c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-merge"))
	if exists {
		return REBASE_MODE_INTERACTIVE, err
	} else {
		return REBASE_MODE_NONE, err
	}
}

func (c *GitCommand) WorkingTreeState() RebaseMode {
	rebaseMode, _ := c.RebaseMode()
	if rebaseMode != REBASE_MODE_NONE {
		return REBASE_MODE_REBASING
	}
	merging, _ := c.IsInMergeState()
	if merging {
		return REBASE_MODE_MERGING
	}
	return REBASE_MODE_NONE
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
