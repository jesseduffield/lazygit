package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
		node := gui.State.Contexts.CommitFiles.GetSelected()
		if node != nil {
			fileName = node.GetPath()
		}
	}

	menuItems := []*types.MenuItem{}

	if fileName != "" {
		menuItems = append(menuItems, &types.MenuItem{
			Label: fmt.Sprintf("%s '%s'", gui.c.Tr.LcFilterBy, fileName),
			OnPress: func() error {
				return gui.setFiltering(fileName)
			},
		})
	}

	menuItems = append(menuItems, &types.MenuItem{
		Label: gui.c.Tr.LcFilterPathOption,
		OnPress: func() error {
			return gui.c.Prompt(types.PromptOpts{
				FindSuggestionsFunc: gui.helpers.Suggestions.GetFilePathSuggestionsFunc(),
				Title:               gui.c.Tr.EnterFileName,
				HandleConfirm: func(response string) error {
					return gui.setFiltering(strings.TrimSpace(response))
				},
			})
		},
	})

	if gui.State.Modes.Filtering.Active() {
		menuItems = append(menuItems, &types.MenuItem{
			Label:   gui.c.Tr.LcExitFilterMode,
			OnPress: gui.clearFiltering,
		})
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.FilteringMenuTitle, Items: menuItems})
}
