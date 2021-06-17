package commands

import (
	"path/filepath"
	"regexp"
	"strings"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RebasingMode int

const (
	REBASE_MODE_NONE RebasingMode = iota
	REBASE_MODE_NON_INTERACTIVE
	REBASE_MODE_INTERACTIVE
)

//counterfeiter:generate . IStatusMgr
type IStatusMgr interface {
	RebaseMode() RebasingMode
	IsMerging() bool
	IsRebasing() bool
	InNormalWorkingTreeState() bool
	IsBareRepo() bool
	IsHeadDetached() bool
	CurrentBranchName() (string, string, error)
}

type StatusMgr struct {
	*MgrCtx
}

func NewStatusMgr(mgrCtx *MgrCtx) *StatusMgr {
	return &StatusMgr{MgrCtx: mgrCtx}
}

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *StatusMgr) RebaseMode() RebasingMode {
	if c.gitDirFileExists("rebase-apply") {
		return REBASE_MODE_NON_INTERACTIVE
	}

	if c.gitDirFileExists("rebase-merge") {
		return REBASE_MODE_INTERACTIVE
	}

	return REBASE_MODE_NONE
}

func (c *StatusMgr) IsRebasing() bool {
	switch c.RebaseMode() {
	case REBASE_MODE_NON_INTERACTIVE, REBASE_MODE_INTERACTIVE:
		return true
	default:
		return false
	}
}

// IsMerging states whether we are still mid-merge
func (c *StatusMgr) IsMerging() bool {
	return c.gitDirFileExists("MERGE_HEAD")
}

func (c *StatusMgr) InNormalWorkingTreeState() bool {
	return !c.IsRebasing() && !c.IsMerging()
}

// arguably this belongs somewhere else. Unlike the other functions in this file,
// the result of this function will not change unless we switch repos.
func (c *StatusMgr) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := c.repo.Worktree()
	return err == gogit.ErrIsBareRepository
}

func (c *StatusMgr) gitDirFileExists(path string) bool {
	result, err := c.os.FileExists(filepath.Join(c.config.GetDotGitDir(), path))
	if err != nil {
		// swallowing error
		c.log.Error(err)
	}

	return result
}

func (c *StatusMgr) IsHeadDetached() bool {
	err := c.RunGitCmdFromStr("symbolic-ref -q HEAD")
	return err != nil
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *StatusMgr) CurrentBranchName() (string, string, error) {
	branchName, err := c.RunWithOutput(
		c.BuildGitCmdObjFromStr("symbolic-ref --short HEAD"),
	)

	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}

	output, err := c.RunWithOutput(
		c.BuildGitCmdObjFromStr("branch --contains"),
	)
	if err != nil {
		return "", "", err
	}

	for _, line := range utils.SplitLines(output) {
		re := regexp.MustCompile(CurrentBranchNameRegex)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			branchName = match[1]
			displayBranchName := match[0][2:]
			return branchName, displayBranchName, nil
		}
	}

	return "HEAD", "HEAD", nil
}
