package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleCreateScopingMenuPanel(g *gocui.Gui, v *gocui.View) error {
	fileName := ""
	switch v.Name() {
	case "files":
		file, err := gui.getSelectedFile(gui.g)
		if err == nil {
			fileName = file.Name
		}
	case "commitFiles":
		file := gui.getSelectedCommitFile(gui.g)
		if file != nil {
			fileName = file.Name
		}
	}

	menuItems := []*menuItem{}

	if fileName != "" {
		menuItems = append(menuItems, &menuItem{
			displayString: fmt.Sprintf("%s '%s'", gui.Tr.SLocalize("filterBy"), fileName),
			onPress: func() error {
				gui.State.FilterPath = fileName
				return gui.Errors.ErrRestart
			},
		})
	}

	menuItems = append(menuItems, &menuItem{
		displayString: gui.Tr.SLocalize("filterPathOption"),
		onPress: func() error {
			return gui.createPromptPanel(gui.g, v, gui.Tr.SLocalize("enterFileName"), "", func(g *gocui.Gui, promptView *gocui.View) error {
				gui.State.FilterPath = strings.TrimSpace(promptView.Buffer())
				return gui.Errors.ErrRestart
			})
		},
	})

	if gui.inFilterMode() {
		menuItems = append(menuItems, &menuItem{
			displayString: gui.Tr.SLocalize("exitFilterMode"),
			onPress: func() error {
				gui.State.FilterPath = ""
				return gui.Errors.ErrRestart
			},
		})
	}

	return gui.createMenu(gui.Tr.SLocalize("FilteringMenuTitle"), menuItems, createMenuOptions{showCancel: true})
}
