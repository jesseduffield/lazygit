package gui

import (
	"fmt"
	"strings"
)

func (gui *Gui) handleCreateFilteringMenuPanel() error {
	fileName := ""
	switch gui.currentSideListContext() {
	case gui.State.Contexts.Files:
		node := gui.getSelectedFileNode()
		if node != nil {
			fileName = node.GetPath()
		}
	case gui.State.Contexts.CommitFiles:
		node := gui.getSelectedCommitFileNode()
		if node != nil {
			fileName = node.GetPath()
		}
	}

	menuItems := []*menuItem{}

	if fileName != "" {
		menuItems = append(menuItems, &menuItem{
			displayString: fmt.Sprintf("%s '%s'", gui.Tr.LcFilterBy, fileName),
			onPress: func() error {
				return gui.setFiltering(fileName)
			},
		})
	}

	menuItems = append(menuItems, &menuItem{
		displayString: gui.Tr.LcFilterPathOption,
		onPress: func() error {
			return gui.prompt(promptOpts{
				title: gui.Tr.LcEnterFileName,
				handleConfirm: func(response string) error {
					return gui.setFiltering(strings.TrimSpace(response))
				},
			})
		},
	})

	if gui.State.Modes.Filtering.Active() {
		menuItems = append(menuItems, &menuItem{
			displayString: gui.Tr.LcExitFilterMode,
			onPress:       gui.clearFiltering,
		})
	}

	return gui.createMenu(gui.Tr.FilteringMenuTitle, menuItems, createMenuOptions{showCancel: true})
}
