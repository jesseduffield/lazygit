package git_commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
)

type StatusCommands struct {
	*GitCommon
}

func NewStatusCommands(
	gitCommon *GitCommon,
) *StatusCommands {
	return &StatusCommands{
		GitCommon: gitCommon,
	}
}

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (self *StatusCommands) RebaseMode() (enums.RebaseMode, error) {
	exists, err := self.os.FileExists(filepath.Join(self.dotGitDir, "rebase-apply"))
	if err != nil {
		return enums.REBASE_MODE_NONE, err
	}
	if exists {
		return enums.REBASE_MODE_NORMAL, nil
	}
	exists, err = self.os.FileExists(filepath.Join(self.dotGitDir, "rebase-merge"))
	if exists {
		return enums.REBASE_MODE_INTERACTIVE, err
	} else {
		return enums.REBASE_MODE_NONE, err
	}
}

func (self *StatusCommands) WorkingTreeState() enums.RebaseMode {
	rebaseMode, _ := self.RebaseMode()
	if rebaseMode != enums.REBASE_MODE_NONE {
		return enums.REBASE_MODE_REBASING
	}
	merging, _ := self.IsInMergeState()
	if merging {
		return enums.REBASE_MODE_MERGING
	}
	return enums.REBASE_MODE_NONE
}

// IsInMergeState states whether we are still mid-merge
func (self *StatusCommands) IsInMergeState() (bool, error) {
	return self.os.FileExists(filepath.Join(self.dotGitDir, "MERGE_HEAD"))
}

func (self *StatusCommands) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := self.repo.Worktree()
	return err == gogit.ErrIsBareRepository
}
