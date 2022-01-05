package gui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func (gui *Gui) handleCreateResetMenu() error {
	red := style.FgRed

	nukeStr := "reset --hard HEAD && git clean -fd"
	if len(gui.State.Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, gui.Tr.LcAndResetSubmodules)
	}

	menuItems := []*menuItem{
		{
			displayStrings: []string{
				gui.Tr.LcDiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.NukeWorkingTree)
				if err := gui.GitCommand.ResetAndClean(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.DiscardUnstagedFileChanges)
				if err := gui.GitCommand.DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.RemoveUntrackedFiles)
				if err := gui.GitCommand.RemoveUntrackedFiles(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.SoftReset)
				if err := gui.GitCommand.ResetSoft("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.MixedReset)
				if err := gui.GitCommand.ResetMixed("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			onPress: func() error {
				gui.logSpan(gui.Tr.Spans.HardReset)
				if err := gui.GitCommand.ResetHard("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
