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
func (self *StatusCommands) RebaseMode() enums.RebaseMode {
	exists, err := self.os.FileExists(filepath.Join(self.dotGitDir, "rebase-apply"))
	if err != nil {
		self.Log.Error("Checking if 'rebase-apply' exists failed", err)
		return enums.REBASE_MODE_NONE
	}
	if exists {
		return enums.REBASE_MODE_NORMAL
	}
	exists, err = self.os.FileExists(filepath.Join(self.dotGitDir, "rebase-merge"))
	if err != nil {
		self.Log.Error("Checking if 'rebase-merge' exists failed", err)
	}
	if exists {
		return enums.REBASE_MODE_INTERACTIVE
	} else {
		return enums.REBASE_MODE_NONE
	}
}

func (self *StatusCommands) WorkingTreeState() enums.RebaseMode {
	rebaseMode := self.RebaseMode()
	if rebaseMode != enums.REBASE_MODE_NONE {
		return enums.REBASE_MODE_REBASING
	}
	if self.Merging() {
		return enums.REBASE_MODE_MERGING
	}
	return enums.REBASE_MODE_NONE
}

// Rebasing states whether we are still mid-rebase
func (self *StatusCommands) Rebasing() bool {
	rebaseMode := self.RebaseMode()
	return rebaseMode != enums.REBASE_MODE_NONE
}

// Merging states whether we are still mid-merge
func (self *StatusCommands) Merging() bool {
	merging, err := self.os.FileExists(filepath.Join(self.dotGitDir, "MERGE_HEAD"))
	if err != nil {
		self.Log.Error("Checking if 'MERGE_HEAD' exists failed", err)
	}
	return merging
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
