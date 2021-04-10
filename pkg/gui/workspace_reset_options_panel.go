package gui

import (
	"fmt"

	"github.com/fatih/color"
)

func (gui *Gui) handleCreateResetMenu() error {
	red := color.New(color.FgRed)

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
				if err := gui.GitCommand.WithSpan("Nuke working tree").ResetAndClean(); err != nil {
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
				if err := gui.GitCommand.WithSpan("Discard unstaged file changes").DiscardAnyUnstagedFileChanges(); err != nil {
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
				if err := gui.GitCommand.WithSpan("Remove untracked files").RemoveUntrackedFiles(); err != nil {
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
				if err := gui.GitCommand.WithSpan("Soft reset").ResetSoft("HEAD"); err != nil {
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
				if err := gui.GitCommand.WithSpan("Mixed reset").ResetMixed("HEAD"); err != nil {
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
				if err := gui.GitCommand.WithSpan("Hard reset").ResetHard("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
			},
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
