package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const HORIZONTAL_SCROLL_FACTOR = 3

func (gui *Gui) scrollUpView(view *gocui.View) {
	view.ScrollUp(gui.c.UserConfig.Gui.ScrollHeight)
}

func (gui *Gui) scrollDownView(view *gocui.View) {
	scrollHeight := gui.c.UserConfig.Gui.ScrollHeight
	view.ScrollDown(scrollHeight)

	if manager, ok := gui.viewBufferManagerMap[view.Name()]; ok {
		manager.ReadLines(scrollHeight)
	}
}

func (gui *Gui) scrollUpMain() error {
	var view *gocui.View
	if gui.c.CurrentContext().GetWindowName() == "secondary" {
		view = gui.secondaryView()
	} else {
		view = gui.mainView()
	}

	if view.Name() == "mergeConflicts" {
		// although we have this same logic in the controller, this method can be invoked
		// via the global scroll up/down keybindings, as opposed to just the mouse wheel keybinding.
		// It would be nice to have a concept of a global keybinding that runs on the top context in a
		// window but that might be overkill for this one use case.
		gui.State.Contexts.MergeConflicts.SetUserScrolling(true)
	}

	gui.scrollUpView(view)

	return nil
}

func (gui *Gui) scrollDownMain() error {
	var view *gocui.View
	if gui.c.CurrentContext().GetWindowName() == "secondary" {
		view = gui.secondaryView()
	} else {
		view = gui.mainView()
	}

	if view.Name() == "mergeConflicts" {
		gui.State.Contexts.MergeConflicts.SetUserScrolling(true)
	}

	gui.scrollDownView(view)

	return nil
}

func (gui *Gui) mainView() *gocui.View {
	viewName := gui.helpers.Window.GetViewNameForWindow("main")
	view, _ := gui.g.View(viewName)
	return view
}

func (gui *Gui) secondaryView() *gocui.View {
	viewName := gui.helpers.Window.GetViewNameForWindow("secondary")
	view, _ := gui.g.View(viewName)
	return view
}

func (gui *Gui) scrollUpSecondary() error {
	gui.scrollUpView(gui.secondaryView())

	return nil
}

func (gui *Gui) scrollDownSecondary() error {
	secondaryView := gui.secondaryView()

	gui.scrollDownView(secondaryView)

	return nil
}

func (gui *Gui) scrollUpConfirmationPanel() error {
	if gui.Views.Confirmation.Editable {
		return nil
	}

	gui.scrollUpView(gui.Views.Confirmation)

	return nil
}

func (gui *Gui) scrollDownConfirmationPanel() error {
	if gui.Views.Confirmation.Editable {
		return nil
	}

	gui.scrollDownView(gui.Views.Confirmation)

	return nil
}

func (gui *Gui) handleCopySelectedSideContextItemToClipboard() error {
	// important to note that this assumes we've selected an item in a side context
	currentSideContext := gui.c.CurrentSideContext()
	if currentSideContext == nil {
		return nil
	}

	listContext, ok := currentSideContext.(types.IListContext)
	if !ok {
		return nil
	}

	itemId := listContext.GetSelectedItemId()

	if itemId == "" {
		return nil
	}

	gui.c.LogAction(gui.c.Tr.Actions.CopyToClipboard)
	if err := gui.os.CopyToClipboard(itemId); err != nil {
		return gui.c.Error(err)
	}

	truncatedItemId := utils.TruncateWithEllipsis(strings.Replace(itemId, "\n", " ", -1), 50)

	gui.c.Toast(fmt.Sprintf("'%s' %s", truncatedItemId, gui.c.Tr.CopiedToClipboard))

	return nil
}
