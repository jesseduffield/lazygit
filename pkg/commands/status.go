package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (enums.RebaseMode, error) {
	exists, err := c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-apply"))
	if err != nil {
		return enums.REBASE_MODE_NONE, err
	}
	if exists {
		return enums.REBASE_MODE_NORMAL, nil
	}
	exists, err = c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-merge"))
	if exists {
		return enums.REBASE_MODE_INTERACTIVE, err
	} else {
		return enums.REBASE_MODE_NONE, err
	}
}

func (c *GitCommand) WorkingTreeState() enums.RebaseMode {
	rebaseMode, _ := c.RebaseMode()
	if rebaseMode != enums.REBASE_MODE_NONE {
		return enums.REBASE_MODE_REBASING
	}
	merging, _ := c.IsInMergeState()
	if merging {
		return enums.REBASE_MODE_MERGING
	}
	return enums.REBASE_MODE_NONE
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
