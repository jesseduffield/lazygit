package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// this is in its own file given that the workspace controller file is already quite long

func (self *FilesController) createResetMenu() error {
	red := style.FgRed

	nukeStr := "git reset --hard HEAD && git clean -fd"
	if len(self.c.Model().Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, self.c.Tr.AndResetSubmodules)
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{
				self.c.Tr.DiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.NukeWorkingTree)
				if err := self.c.Git().WorkingTree.ResetAndClean(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key:     'x',
			Tooltip: self.c.Tr.NukeDescription,
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedFileChanges)
				if err := self.c.Git().WorkingTree.DiscardAnyUnstagedFileChanges(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 'u',
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RemoveUntrackedFiles)
				if err := self.c.Git().WorkingTree.RemoveUntrackedFiles(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 'c',
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardStagedChanges,
				red.Sprint("stash staged and drop stash"),
			},
			Tooltip: self.c.Tr.DiscardStagedChangesDescription,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RemoveStagedFiles)
				if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
					return self.c.ErrorMsg(self.c.Tr.NoTrackedStagedFilesStash)
				}
				if err := self.c.Git().Stash.SaveStagedChanges("[lazygit] tmp stash"); err != nil {
					return self.c.Error(err)
				}
				if err := self.c.Git().Stash.DropNewest(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 'S',
		},
		{
			LabelColumns: []string{
				self.c.Tr.SoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.SoftReset)
				if err := self.c.Git().WorkingTree.ResetSoft("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 's',
		},
		{
			LabelColumns: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.MixedReset)
				if err := self.c.Git().WorkingTree.ResetMixed("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 'm',
		},
		{
			LabelColumns: []string{
				self.c.Tr.HardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.HardReset)
				if err := self.c.Git().WorkingTree.ResetHard("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
			Key: 'h',
		},
	}

	return self.c.Menu(types.CreateMenuOptions{Title: "", Items: menuItems})
}
