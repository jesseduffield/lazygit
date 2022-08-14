package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const HORIZONTAL_SCROLL_FACTOR = 3

// these views need to be re-rendered when the screen mode changes. The commits view,
// for example, will show authorship information in half and full screen mode.
func (gui *Gui) rerenderViewsWithScreenModeDependentContent() error {
	// for now we re-render all list views.
	for _, context := range gui.getListContexts() {
		if err := gui.rerenderView(context.GetView()); err != nil {
			return err
		}
	}

	return nil
}

func nextIntInCycle(sl []WindowMaximisation, current WindowMaximisation) WindowMaximisation {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

func prevIntInCycle(sl []WindowMaximisation, current WindowMaximisation) WindowMaximisation {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}

func (gui *Gui) nextScreenMode() error {
	gui.State.ScreenMode = nextIntInCycle([]WindowMaximisation{SCREEN_NORMAL, SCREEN_HALF, SCREEN_FULL}, gui.State.ScreenMode)

	return gui.rerenderViewsWithScreenModeDependentContent()
}

func (gui *Gui) prevScreenMode() error {
	gui.State.ScreenMode = prevIntInCycle([]WindowMaximisation{SCREEN_NORMAL, SCREEN_HALF, SCREEN_FULL}, gui.State.ScreenMode)

	return gui.rerenderViewsWithScreenModeDependentContent()
}

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
	viewName := gui.getViewNameForWindow("main")
	view, _ := gui.g.View(viewName)
	return view
}

func (gui *Gui) secondaryView() *gocui.View {
	viewName := gui.getViewNameForWindow("secondary")
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

func (gui *Gui) handleRefresh() error {
	return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) backgroundFetch() (err error) {
	err = gui.git.Sync.Fetch(git_commands.FetchOptions{Background: true})

	_ = gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.ASYNC})

	return err
}

func (gui *Gui) handleCopySelectedSideContextItemToClipboard() error {
	// important to note that this assumes we've selected an item in a side context
	itemId := gui.getSideContextSelectedItemId()

	if itemId == "" {
		return nil
	}

	gui.c.LogAction(gui.c.Tr.Actions.CopyToClipboard)
	if err := gui.os.CopyToClipboard(itemId); err != nil {
		return gui.c.Error(err)
	}

	truncatedItemId := utils.TruncateWithEllipsis(strings.Replace(itemId, "\n", " ", -1), 50)

	gui.c.Toast(fmt.Sprintf("'%s' %s", truncatedItemId, gui.c.Tr.LcCopiedToClipboard))

	return nil
}
