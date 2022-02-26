package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// this is in its own file given that the workspace controller file is already quite long

func (self *FilesController) createResetMenu() error {
	red := style.FgRed

	nukeStr := "reset --hard HEAD && git clean -fd"
	if len(self.model.Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, self.c.Tr.LcAndResetSubmodules)
	}

	menuItems := []*types.MenuItem{
		{
			DisplayStrings: []string{
				self.c.Tr.LcDiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.NukeWorkingTree)
				if err := self.git.WorkingTree.ResetAndClean(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				self.c.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedFileChanges)
				if err := self.git.WorkingTree.DiscardAnyUnstagedFileChanges(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				self.c.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RemoveUntrackedFiles)
				if err := self.git.WorkingTree.RemoveUntrackedFiles(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				self.c.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.SoftReset)
				if err := self.git.WorkingTree.ResetSoft("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.MixedReset)
				if err := self.git.WorkingTree.ResetMixed("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				self.c.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.HardReset)
				if err := self.git.WorkingTree.ResetHard("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
	}

	return self.c.Menu(types.CreateMenuOptions{Title: "", Items: menuItems})
}
