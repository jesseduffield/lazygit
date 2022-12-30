package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateDiffingMenuPanel() error {
	names := gui.helpers.Diff.CurrentDiffTerminals()

	menuItems := []*types.MenuItem{}
	for _, name := range names {
		name := name
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: fmt.Sprintf("%s %s", gui.c.Tr.LcDiff, name),
				OnPress: func() error {
					gui.State.Modes.Diffing.Ref = name
					// can scope this down based on current view but too lazy right now
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	menuItems = append(menuItems, []*types.MenuItem{
		{
			Label: gui.c.Tr.LcEnterRefToDiff,
			OnPress: func() error {
				return gui.c.Prompt(types.PromptOpts{
					Title:               gui.c.Tr.LcEnteRefName,
					FindSuggestionsFunc: gui.helpers.Suggestions.GetRefsSuggestionsFunc(),
					HandleConfirm: func(response string) error {
						gui.State.Modes.Diffing.Ref = strings.TrimSpace(response)
						return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
					},
				})
			},
		},
	}...)

	if gui.State.Modes.Diffing.Active() {
		menuItems = append(menuItems, []*types.MenuItem{
			{
				Label: gui.c.Tr.LcSwapDiff,
				OnPress: func() error {
					gui.State.Modes.Diffing.Reverse = !gui.State.Modes.Diffing.Reverse
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
			{
				Label: gui.c.Tr.LcExitDiffMode,
				OnPress: func() error {
					gui.State.Modes.Diffing = diffing.New()
					return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				},
			},
		}...)
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.DiffingMenuTitle, Items: menuItems})
}
