package git_commands

import (
	"os"
	"path/filepath"
	"strings"

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

func (self *StatusCommands) WorkingTreeState() enums.WorkingTreeState {
	isInRebase, _ := self.IsInRebase()
	if isInRebase {
		return enums.WORKING_TREE_STATE_REBASING
	}
	merging, _ := self.IsInMergeState()
	if merging {
		return enums.WORKING_TREE_STATE_MERGING
	}
	return enums.WORKING_TREE_STATE_NONE
}

func (self *StatusCommands) IsBareRepo() bool {
	return self.repoPaths.isBareRepo
}

func (self *StatusCommands) IsInRebase() (bool, error) {
	exists, err := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge"))
	if err == nil && exists {
		return true, nil
	}
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply"))
}

// IsInMergeState states whether we are still mid-merge
func (self *StatusCommands) IsInMergeState() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "MERGE_HEAD"))
}

// Full ref (e.g. "refs/heads/mybranch") of the branch that is currently
// being rebased, or empty string when we're not in a rebase
func (self *StatusCommands) BranchBeingRebased() string {
	for _, dir := range []string{"rebase-merge", "rebase-apply"} {
		if bytesContent, err := os.ReadFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), dir, "head-name")); err == nil {
			return strings.TrimSpace(string(bytesContent))
		}
	}
	return ""
}
