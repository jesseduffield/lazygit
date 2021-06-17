package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

type RebasingMode int

const (
	REBASE_MODE_NONE RebasingMode = iota
	REBASE_MODE_NON_INTERACTIVE
	REBASE_MODE_INTERACTIVE
)

type IStatusMgr interface {
	RebaseMode() RebasingMode
	IsMerging() bool
	IsRebasing() bool
	InNormalWorkingTreeState() bool
	IsBareRepo() bool
	IsHeadDetached() bool
}

type StatusMgr struct {
	ICommander

	config IGitConfigMgr

	os   oscommands.IOS
	repo *gogit.Repository
	log  *logrus.Entry
}

func NewStatusMgr(
	commander ICommander,
	config IGitConfigMgr,
	os oscommands.IOS,
	repo *gogit.Repository,
	log *logrus.Entry,
) *StatusMgr {
	return &StatusMgr{
		ICommander: commander,
		config:     config,
		os:         os,
		repo:       repo,
		log:        log,
	}
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
