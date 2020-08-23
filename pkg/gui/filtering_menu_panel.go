package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateFilteringMenuPanel(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	fileName := ""
	switch v.Name() {
	case "files":
		file := gui.getSelectedFile()
		if file != nil {
			fileName = file.Name
		}
	case "commitFiles":
		file := gui.getSelectedCommitFile()
		if file != nil {
			fileName = file.Name
		}
	}

	menuItems := []*menuItem{}

	if fileName != "" {
		menuItems = append(menuItems, &menuItem{
			displayString: fmt.Sprintf("%s '%s'", gui.Tr.SLocalize("filterBy"), fileName),
			onPress: func() error {
				gui.State.Modes.Filtering.Path = fileName
				return gui.Errors.ErrRestart
			},
		})
	}

	menuItems = append(menuItems, &menuItem{
		displayString: gui.Tr.SLocalize("filterPathOption"),
		onPress: func() error {
			return gui.prompt(gui.Tr.SLocalize("enterFileName"), "", func(response string) error {
				gui.State.Modes.Filtering.Path = strings.TrimSpace(response)
				return gui.Errors.ErrRestart
			})
		},
	})

	if gui.State.Modes.Filtering.Active() {
		menuItems = append(menuItems, &menuItem{
			displayString: gui.Tr.SLocalize("exitFilterMode"),
			onPress: func() error {
				gui.State.Modes.Filtering.Path = ""
				return gui.Errors.ErrRestart
			},
		})
	}

	return gui.createMenu(gui.Tr.SLocalize("FilteringMenuTitle"), menuItems, createMenuOptions{showCancel: true})
}
