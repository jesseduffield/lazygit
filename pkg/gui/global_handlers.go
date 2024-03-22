package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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
	return gui.handleCopySelectedSideContextItemToClipboardWithTruncation(-1)
}

func (gui *Gui) handleCopySelectedSideContextItemCommitHashToClipboard() error {
	return gui.handleCopySelectedSideContextItemToClipboardWithTruncation(
		gui.UserConfig.Git.TruncateCopiedCommitHashesTo)
}

func (gui *Gui) handleCopySelectedSideContextItemToClipboardWithTruncation(maxWidth int) error {
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

	if maxWidth > 0 {
		itemId = itemId[:utils.Min(len(itemId), maxWidth)]
	}

	gui.c.LogAction(gui.c.Tr.Actions.CopyToClipboard)
	if err := gui.os.CopyToClipboard(itemId); err != nil {
		return gui.c.Error(err)
	}

	truncatedItemId := utils.TruncateWithEllipsis(strings.Replace(itemId, "\n", " ", -1), 50)

	gui.c.Toast(fmt.Sprintf("'%s' %s", truncatedItemId, gui.c.Tr.CopiedToClipboard))

	return nil
}

func (gui *Gui) getCopySelectedSideContextItemToClipboardDisabledReason() *types.DisabledReason {
	// important to note that this assumes we've selected an item in a side context
	currentSideContext := gui.c.CurrentSideContext()
	if currentSideContext == nil {
		// This should never happen but if it does we'll just ignore the keypress
		return nil
	}

	listContext, ok := currentSideContext.(types.IListContext)
	if !ok {
		// This should never happen but if it does we'll just ignore the keypress
		return nil
	}

	startIdx, endIdx := listContext.GetList().GetSelectionRange()
	if startIdx != endIdx {
		return &types.DisabledReason{Text: gui.Tr.RangeSelectNotSupported}
	}

	return nil
}

func (gui *Gui) setCaption(caption string) {
	gui.Views.Options.FgColor = gocui.ColorWhite
	gui.Views.Options.FgColor |= gocui.AttrBold
	gui.Views.Options.SetContent(captionPrefix + " " + style.FgCyan.SetBold().Sprint(caption))
	gui.c.Render()
}

var captionPrefix = ""

func (gui *Gui) setCaptionPrefix(prefix string) {
	gui.Views.Options.FgColor = gocui.ColorWhite
	gui.Views.Options.FgColor |= gocui.AttrBold

	captionPrefix = prefix

	gui.Views.Options.SetContent(prefix)
	gui.c.Render()
}
