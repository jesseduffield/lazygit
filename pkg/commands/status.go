package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
)

type WorkingTreeState string

const (
	REBASE_MODE_NORMAL      WorkingTreeState = "normal"
	REBASE_MODE_INTERACTIVE                  = "interactive"
	REBASE_MODE_REBASING                     = "rebasing"
	REBASE_MODE_MERGING                      = "merging"
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *Git) RebaseMode() (WorkingTreeState, error) {
	exists, err := c.GetOS().FileExists(filepath.Join(c.dotGitDir, "rebase-apply"))
	if err != nil {
		return "", err
	}
	if exists {
		return REBASE_MODE_NORMAL, nil
	}
	exists, err = c.GetOS().FileExists(filepath.Join(c.dotGitDir, "rebase-merge"))
	if exists {
		return REBASE_MODE_INTERACTIVE, err
	} else {
		return "", err
	}
}

func (c *Git) WorkingTreeState() WorkingTreeState {
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
func (c *Git) IsInMergeState() (bool, error) {
	return c.GetOS().FileExists(filepath.Join(c.dotGitDir, "MERGE_HEAD"))
}

func (c *Git) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := c.repo.Worktree()
	return err == gogit.ErrIsBareRepository
}
