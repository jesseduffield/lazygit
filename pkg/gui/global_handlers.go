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
	for _, view := range []*gocui.View{gui.Views.Branches, gui.Views.Commits} {
		if err := gui.rerenderView(view); err != nil {
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
	if gui.renderingConflicts() {
		gui.State.Panels.Merging.UserVerticalScrolling = true
	}

	gui.scrollUpView(gui.Views.Main)

	return nil
}

func (gui *Gui) scrollDownMain() error {
	if gui.renderingConflicts() {
		gui.State.Panels.Merging.UserVerticalScrolling = true
	}

	gui.scrollDownView(gui.Views.Main)

	return nil
}

func (gui *Gui) scrollLeftMain() error {
	gui.scrollLeft(gui.Views.Main)

	return nil
}

func (gui *Gui) scrollRightMain() error {
	gui.scrollRight(gui.Views.Main)

	return nil
}

func (gui *Gui) scrollLeft(view *gocui.View) {
	newOriginX := utils.Max(view.OriginX()-view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR, 0)
	_ = view.SetOriginX(newOriginX)
}

func (gui *Gui) scrollRight(view *gocui.View) {
	_ = view.SetOriginX(view.OriginX() + view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR)
}

func (gui *Gui) scrollUpSecondary() error {
	gui.scrollUpView(gui.Views.Secondary)

	return nil
}

func (gui *Gui) scrollDownSecondary() error {
	gui.scrollDownView(gui.Views.Secondary)

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
