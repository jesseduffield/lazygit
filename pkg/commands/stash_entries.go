package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

// StashDo modify stash
func (c *Git) StashDo(index int, method string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("stash %s stash@{%d}", method, index))
}

// StashSave save stash
// TODO: before calling this, check if there is anything to save
func (c *Git) StashSave(message string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("stash save %s", c.GetOS().Quote(message)))
}

// GetStashEntryDiff stash diff
func (c *Git) ShowStashEntryCmdObj(index int) ICmdObj {
	return BuildGitCmdObjFromStr(
		fmt.Sprintf("stash show -p --stat --color=%s stash@{%d}", c.colorArg(), index),
	)
}

// StashSaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (c *Git) StashSaveStagedChanges(message string) error {
	// wrap in 'writing', which uses a mutex
	if err := c.RunGitCmdFromStr("stash --keep-index"); err != nil {
		return err
	}

	if err := c.StashSave(message); err != nil {
		return err
	}

	if err := c.RunGitCmdFromStr("stash apply stash@{1}"); err != nil {
		return err
	}

	err := c.GetOS().PipeCommands(
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
	files := c.GetStatusFiles(loaders.LoadStatusFilesOpts{})
	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := c.UnStageFile(file.Names(), false); err != nil {
				return err
			}
		}
	}

	return nil
}
