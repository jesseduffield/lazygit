package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateResetMenu() error {
	red := style.FgRed

	nukeStr := "reset --hard HEAD && git clean -fd"
	if len(gui.State.Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, gui.c.Tr.LcAndResetSubmodules)
	}

	menuItems := []*types.MenuItem{
		{
			DisplayStrings: []string{
				gui.c.Tr.LcDiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.NukeWorkingTree)
				if err := gui.git.WorkingTree.ResetAndClean(); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.c.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.DiscardUnstagedFileChanges)
				if err := gui.git.WorkingTree.DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.c.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.RemoveUntrackedFiles)
				if err := gui.git.WorkingTree.RemoveUntrackedFiles(); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.c.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.SoftReset)
				if err := gui.git.WorkingTree.ResetSoft("HEAD"); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.MixedReset)
				if err := gui.git.WorkingTree.ResetMixed("HEAD"); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.c.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			OnPress: func() error {
				gui.c.LogAction(gui.c.Tr.Actions.HardReset)
				if err := gui.git.WorkingTree.ResetHard("HEAD"); err != nil {
					return gui.c.Error(err)
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: "", Items: menuItems})
}
