package git_commands

import (
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

func (self *StatusCommands) IsBareRepo() (bool, error) {
	return IsBareRepo(self.os)
}

func IsBareRepo(osCommand *oscommands.OSCommand) (bool, error) {
	res, err := osCommand.Cmd.New("git rev-parse --is-bare-repository").DontLog().RunWithOutput()
	if err != nil {
		return false, err
	}

	// The command returns output with a newline, so we need to strip
	return strconv.ParseBool(strings.TrimSpace(res))
}

// IsInMergeState states whether we are still mid-merge
func (self *StatusCommands) IsInMergeState() (bool, error) {
	return self.os.FileExists(filepath.Join(self.dotGitDir, "MERGE_HEAD"))
}
