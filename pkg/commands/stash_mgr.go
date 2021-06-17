package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

type IStashMgr interface {
	Do(index int, method string) error
	Save(message string) error
	ShowEntryCmdObj(index int) ICmdObj
	SaveStagedChanges(message string) error
	LoadEntries(filterPath string) []*models.StashEntry
}

type StashMgr struct {
	*MgrCtx

	stashEntriesLoader *StashEntriesLoader
	worktreeMgr        IWorktreeMgr
}

func NewStashMgr(mgrCtx *MgrCtx, worktreeMgr IWorktreeMgr) *StashMgr {
	stashEntriesLoader := NewStashEntriesLoader(mgrCtx)

	return &StashMgr{
		MgrCtx:             mgrCtx,
		stashEntriesLoader: stashEntriesLoader,
		worktreeMgr:        worktreeMgr,
	}
}

func (c *StashMgr) LoadEntries(filterPath string) []*models.StashEntry {
	return c.stashEntriesLoader.Load(filterPath)
}

// StashDo modify stash
func (c *StashMgr) Do(index int, method string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("stash %s stash@{%d}", method, index))
}

// TODO: before calling this, check if there is anything to save
func (c *StashMgr) Save(message string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("stash save %s", c.Quote(message)))
}

// GetStashEntryDiff stash diff
func (c *StashMgr) ShowEntryCmdObj(index int) ICmdObj {
	return BuildGitCmdObjFromStr(
		fmt.Sprintf("stash show -p --stat --color=%s stash@{%d}", c.config.ColorArg(), index),
	)
}

// SaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (c *StashMgr) SaveStagedChanges(message string) error {
	// wrap in 'writing', which uses a mutex
	if err := c.RunGitCmdFromStr("stash --keep-index"); err != nil {
		return err
	}

	if err := c.Save(message); err != nil {
		return err
	}

	if err := c.RunGitCmdFromStr("stash apply stash@{1}"); err != nil {
		return err
	}

	err := c.os.PipeCommands(
		BuildGitCmdObjFromStr("stash show -p"),
		BuildGitCmdObjFromStr("apply -R"),
	)
	if err != nil {
		return err
	}

	if err := c.RunGitCmdFromStr("stash drop stash@{1}"); err != nil {
		return err
	}

	// if you had staged an untracked file, that will now appear as 'AD' in git status
	// meaning it's deleted in your working tree but added in your index. Given that it's
	// now safely stashed, we need to remove it.
	files := c.worktreeMgr.LoadStatusFiles(LoadStatusFilesOpts{})
	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := c.worktreeMgr.UnStageFile(file.Names(), false); err != nil {
				return err
			}
		}
	}

	return nil
}
