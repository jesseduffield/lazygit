package gui

import (
	"fmt"
	"math"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// these views need to be re-rendered when the screen mode changes. The commits view,
// for example, will show authorship information in half and full screen mode.
func (gui *Gui) rerenderViewsWithScreenModeDependentContent() error {
	for _, viewName := range []string{"branches", "commits"} {
		if err := gui.rerenderView(viewName); err != nil {
			return err
		}
	}

	return nil
}

// TODO: GENERICS
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

// TODO: GENERICS
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

func (gui *Gui) scrollUpView(viewName string) error {
	mainView, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}
	ox, oy := mainView.Origin()
	newOy := int(math.Max(0, float64(oy-gui.Config.GetUserConfig().Gui.ScrollHeight)))
	return mainView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownView(viewName string) error {
	mainView, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}
	ox, oy := mainView.Origin()
	y := oy
	canScrollPastBottom := gui.Config.GetUserConfig().Gui.ScrollPastBottom
	if !canScrollPastBottom {
		_, sy := mainView.Size()
		y += sy
	}
	scrollHeight := gui.Config.GetUserConfig().Gui.ScrollHeight
	scrollableLines := mainView.ViewLinesHeight() - y
	if scrollableLines > 0 {
		// margin is about how many lines must still appear if you scroll
		// all the way down. In practice every file ends in a newline so it will really
		// just show a single line
		margin := 1
		if canScrollPastBottom {
			margin = 2
		}
		if scrollableLines-margin < scrollHeight {
			scrollHeight = scrollableLines - margin
		}
		if err := mainView.SetOrigin(ox, oy+scrollHeight); err != nil {
			return err
		}
	}
	if manager, ok := gui.viewBufferManagerMap[viewName]; ok {
		manager.ReadLines(scrollHeight)
	}
	return nil
}

func (gui *Gui) scrollUpMain() error {
	if gui.canScrollMergePanel() {
		gui.State.Panels.Merging.UserScrolling = true
	}

	return gui.scrollUpView("main")
}

func (gui *Gui) scrollDownMain() error {
	if gui.canScrollMergePanel() {
		gui.State.Panels.Merging.UserScrolling = true
	}

	return gui.scrollDownView("main")
}

func (gui *Gui) scrollUpSecondary() error {
	return gui.scrollUpView("secondary")
}

func (gui *Gui) scrollDownSecondary() error {
	return gui.scrollDownView("secondary")
}

func (gui *Gui) scrollUpConfirmationPanel() error {
	view := gui.getConfirmationView()
	if view != nil || view.Editable {
		return nil
	}

	return gui.scrollUpView("confirmation")
}

func (gui *Gui) scrollDownConfirmationPanel() error {
	view := gui.getConfirmationView()
	if view != nil || view.Editable {
		return nil
	}

	return gui.scrollDownView("confirmation")
}

func (gui *Gui) handleRefresh() error {
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleMouseDownMain() error {
	if gui.popupPanelFocused() {
		return nil
	}

	view := gui.getMainView()

	switch gui.g.CurrentView().Name() {
	case "files":
		// set filename, set primary/secondary selected, set line number, then switch context
		// I'll need to know it was changed though.
		// Could I pass something along to the context change?
		return gui.enterFile(false, view.SelectedLineIdx())
	case "commitFiles":
		return gui.enterCommitFile(view.SelectedLineIdx())
	}

	return nil
}

func (gui *Gui) handleMouseDownSecondary() error {
	if gui.popupPanelFocused() {
		return nil
	}

	view := gui.getSecondaryView()

	switch gui.g.CurrentView().Name() {
	case "files":
		return gui.enterFile(true, view.SelectedLineIdx())
	}

	return nil
}

func (gui *Gui) handleInfoClick() error {
	if !gui.g.Mouse {
		return nil
	}

	view := gui.getInformationView()

	cx, _ := view.Cursor()
	width, _ := view.Size()

	for _, mode := range gui.modeStatuses() {
		if mode.isActive() {
			if width-cx > len(gui.Tr.ResetInParentheses) {
				return nil
			}
			return mode.reset()
		}
	}

	// if we're not in an active mode we show the donate button
	if cx <= len(gui.Tr.Donate)+len(INFO_SECTION_PADDING) {
		return gui.OSCommand.OpenLink("https://github.com/sponsors/jesseduffield")
	}
	return nil
}

func (gui *Gui) fetch(canPromptForCredentials bool) (err error) {
	gui.Mutexes.FetchMutex.Lock()
	defer gui.Mutexes.FetchMutex.Unlock()

	fetchOpts := commands.FetchOptions{}
	if canPromptForCredentials {
		fetchOpts.PromptUserForCredential = gui.promptUserForCredential
	}

	err = gui.GitCommand.Fetch(fetchOpts)

	if canPromptForCredentials && err != nil && strings.Contains(err.Error(), "exit status 128") {
		_ = gui.createErrorPanel(gui.Tr.PassUnameWrong)
	}

	_ = gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{BRANCHES, COMMITS, REMOTES, TAGS}, mode: ASYNC})

	return err
}

func (gui *Gui) handleCopySelectedSideContextItemToClipboard() error {
	// important to note that this assumes we've selected an item in a side context
	itemId := gui.getSideContextSelectedItemId()

	if itemId == "" {
		return nil
	}

	if err := gui.OSCommand.CopyToClipboard(itemId); err != nil {
		return gui.surfaceError(err)
	}

	truncatedItemId := utils.TruncateWithEllipsis(strings.Replace(itemId, "\n", " ", -1), 50)

	gui.raiseToast(fmt.Sprintf("'%s' %s", truncatedItemId, gui.Tr.LcCopiedToClipboard))

	return nil
}
