package gui

import "github.com/jesseduffield/gocui"

type listView struct {
	viewName          string
	context           string
	getItemsLength    func() int
	getSelectedLine   func() *int
	handleItemSelect  func(g *gocui.Gui, v *gocui.View) error
	gui               *Gui
	rendersToMainView bool
}

func (lv *listView) handlePrevLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(-1)
}

func (lv *listView) handleNextLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(1)
}

func (lv *listView) handleLineChange(change int) error {
	if !lv.gui.isPopupPanel(lv.viewName) && lv.gui.popupPanelFocused() {
		return nil
	}

	lv.gui.changeSelectedLine(lv.getSelectedLine(), lv.getItemsLength(), change)

	if lv.rendersToMainView {
		if err := lv.gui.resetOrigin(lv.gui.getMainView()); err != nil {
			return err
		}
	}
	view, err := lv.gui.g.View(lv.viewName)
	if err != nil {
		return err
	}
	return lv.handleItemSelect(lv.gui.g, view)
}

func (gui *Gui) getListViews() []*listView {
	return []*listView{
		{
			viewName:          "menu",
			getItemsLength:    func() int { return gui.getMenuView().LinesHeight() },
			getSelectedLine:   func() *int { return &gui.State.Panels.Menu.SelectedLine },
			handleItemSelect:  gui.handleMenuSelect,
			gui:               gui,
			rendersToMainView: false,
		},
		{
			viewName:          "files",
			getItemsLength:    func() int { return len(gui.State.Files) },
			getSelectedLine:   func() *int { return &gui.State.Panels.Files.SelectedLine },
			handleItemSelect:  gui.handleFileSelect,
			gui:               gui,
			rendersToMainView: true,
		},
		{
			viewName:          "branches",
			context:           "local-branches",
			getItemsLength:    func() int { return len(gui.State.Branches) },
			getSelectedLine:   func() *int { return &gui.State.Panels.Branches.SelectedLine },
			handleItemSelect:  gui.handleBranchSelect,
			gui:               gui,
			rendersToMainView: true,
		},
		{
			viewName:          "commits",
			getItemsLength:    func() int { return len(gui.State.Commits) },
			getSelectedLine:   func() *int { return &gui.State.Panels.Commits.SelectedLine },
			handleItemSelect:  gui.handleCommitSelect,
			gui:               gui,
			rendersToMainView: true,
		},
		{
			viewName:          "stash",
			getItemsLength:    func() int { return len(gui.State.StashEntries) },
			getSelectedLine:   func() *int { return &gui.State.Panels.Stash.SelectedLine },
			handleItemSelect:  gui.handleStashEntrySelect,
			gui:               gui,
			rendersToMainView: true,
		},
		{
			viewName:          "commitFiles",
			getItemsLength:    func() int { return len(gui.State.CommitFiles) },
			getSelectedLine:   func() *int { return &gui.State.Panels.CommitFiles.SelectedLine },
			handleItemSelect:  gui.handleCommitFileSelect,
			gui:               gui,
			rendersToMainView: true,
		},
	}
}
