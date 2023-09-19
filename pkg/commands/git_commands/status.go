package git_commands

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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
	ok, err := self.IsInNormalRebase()
	if err == nil && ok {
		return enums.REBASE_MODE_NORMAL, nil
	}
	ok, err = self.IsInInteractiveRebase()
	if err == nil && ok {
		return enums.REBASE_MODE_INTERACTIVE, err
	}

	return enums.REBASE_MODE_NONE, err
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

func (self *StatusCommands) IsBareRepo() (bool, error) {
	return IsBareRepo(self.os)
}

func IsBareRepo(osCommand *oscommands.OSCommand) (bool, error) {
	res, err := osCommand.Cmd.New(
		NewGitCmd("rev-parse").Arg("--is-bare-repository").ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		return false, err
	}

	// The command returns output with a newline, so we need to strip
	return strconv.ParseBool(strings.TrimSpace(res))
}

func (self *StatusCommands) IsInNormalRebase() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply"))
}

func (self *StatusCommands) IsInInteractiveRebase() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge"))
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
