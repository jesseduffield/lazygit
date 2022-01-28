package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateResetMenu() error {
	red := style.FgRed

	nukeStr := "reset --hard HEAD && git clean -fd"
	if len(gui.State.Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, gui.Tr.LcAndResetSubmodules)
	}

	menuItems := []*popup.MenuItem{
		{
			DisplayStrings: []string{
				gui.Tr.LcDiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.NukeWorkingTree)
				if err := gui.Git.WorkingTree.ResetAndClean(); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.DiscardUnstagedFileChanges)
				if err := gui.Git.WorkingTree.DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.RemoveUntrackedFiles)
				if err := gui.Git.WorkingTree.RemoveUntrackedFiles(); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.SoftReset)
				if err := gui.Git.WorkingTree.ResetSoft("HEAD"); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.MixedReset)
				if err := gui.Git.WorkingTree.ResetMixed("HEAD"); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
		{
			DisplayStrings: []string{
				gui.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			OnPress: func() error {
				gui.logAction(gui.Tr.Actions.HardReset)
				if err := gui.Git.WorkingTree.ResetHard("HEAD"); err != nil {
					return gui.PopupHandler.Error(err)
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}})
			},
		},
	}

	return gui.PopupHandler.Menu(popup.CreateMenuOptions{Title: "", Items: menuItems})
}
