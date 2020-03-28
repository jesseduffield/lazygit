package gui

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateResetMenu(g *gocui.Gui, v *gocui.View) error {
	red := color.New(color.FgRed)

	menuItems := []*menuItem{
		{
			displayStrings: []string{
				gui.Tr.SLocalize("discardAllChangesToAllFiles"),
				red.Sprint("reset --hard HEAD && git clean -fd"),
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
				gui.Tr.SLocalize("discardAnyUnstagedChanges"),
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
				gui.Tr.SLocalize("discardUntrackedFiles"),
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
				gui.Tr.SLocalize("softReset"),
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
				gui.Tr.SLocalize("hardReset"),
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
