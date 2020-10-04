package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateResetMenu(g *gocui.Gui, v *gocui.View) error {
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
				if err := gui.GitCommand.ResetAndClean(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			onPress: func() error {
				if err := gui.GitCommand.DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcDiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			onPress: func() error {
				if err := gui.GitCommand.RemoveUntrackedFiles(); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcSoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			onPress: func() error {
				if err := gui.GitCommand.ResetSoft("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
		{
			displayStrings: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			onPress: func() error {
				if err := gui.GitCommand.ResetSoft("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
		{
			displayStrings: []string{
				gui.Tr.LcHardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			onPress: func() error {
				if err := gui.GitCommand.ResetHard("HEAD"); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{FILES}})
			},
		},
	}

	return gui.createMenu("", menuItems, createMenuOptions{showCancel: true})
}
