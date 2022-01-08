package git_commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type StatusCommands struct {
	*common.Common
	osCommand *oscommands.OSCommand
	repo      *gogit.Repository
	dotGitDir string
}

func NewStatusCommands(
	common *common.Common,
	osCommand *oscommands.OSCommand,
	repo *gogit.Repository,
	dotGitDir string,
) *StatusCommands {
	return &StatusCommands{
		Common:    common,
		osCommand: osCommand,
		repo:      repo,
		dotGitDir: dotGitDir,
	}
}

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (self *StatusCommands) RebaseMode() (enums.RebaseMode, error) {
	exists, err := self.osCommand.FileExists(filepath.Join(self.dotGitDir, "rebase-apply"))
	if err != nil {
		return enums.REBASE_MODE_NONE, err
	}
	if exists {
		return enums.REBASE_MODE_NORMAL, nil
	}
	exists, err = self.osCommand.FileExists(filepath.Join(self.dotGitDir, "rebase-merge"))
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
	return self.osCommand.FileExists(filepath.Join(self.dotGitDir, "MERGE_HEAD"))
}

func (self *StatusCommands) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := self.repo.Worktree()
	return err == gogit.ErrIsBareRepository
}
