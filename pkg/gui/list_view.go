package gui

import "github.com/jesseduffield/gocui"

type listView struct {
	ViewName              string
	Context               string
	GetItemsLength        func() int
	GetSelectedLineIdxPtr func() *int
	OnFocus               func() error
	OnItemSelect          func() error
	OnClickSelectedItem   func() error
	Gui                   *Gui
	RendersToMainView     bool
}

func (lv *listView) handlePrevLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(-1)
}

func (lv *listView) handleNextLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(1)
}

func (lv *listView) handleLineChange(change int) error {
	if !lv.Gui.isPopupPanel(lv.ViewName) && lv.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lv.Gui.g.View(lv.ViewName)
	if err != nil {
		return err
	}

	lv.Gui.changeSelectedLine(lv.GetSelectedLineIdxPtr(), lv.GetItemsLength(), change)
	view.FocusPoint(0, *lv.GetSelectedLineIdxPtr())

	if lv.RendersToMainView {
		if err := lv.Gui.resetOrigin(lv.Gui.getMainView()); err != nil {
			return err
		}
		if err := lv.Gui.resetOrigin(lv.Gui.getSecondaryView()); err != nil {
			return err
		}
	}

	if lv.OnItemSelect != nil {
		return lv.OnItemSelect()
	}
	return nil
}

func (lv *listView) handleNextPage(g *gocui.Gui, v *gocui.View) error {
	view, err := lv.Gui.g.View(lv.ViewName)
	if err != nil {
		return nil
	}
	_, height := view.Size()
	delta := height - 1
	if delta == 0 {
		delta = 1
	}
	return lv.handleLineChange(delta)
}

func (lv *listView) handleGotoTop(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(-lv.GetItemsLength())
}

func (lv *listView) handleGotoBottom(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(lv.GetItemsLength())
}

func (lv *listView) handlePrevPage(g *gocui.Gui, v *gocui.View) error {
	view, err := lv.Gui.g.View(lv.ViewName)
	if err != nil {
		return nil
	}
	_, height := view.Size()
	delta := height - 1
	if delta == 0 {
		delta = 1
	}
	return lv.handleLineChange(-delta)
}

func (lv *listView) handleClick(g *gocui.Gui, v *gocui.View) error {
	if !lv.Gui.isPopupPanel(lv.ViewName) && lv.Gui.popupPanelFocused() {
		return nil
	}

	selectedLineIdxPtr := lv.GetSelectedLineIdxPtr()
	prevSelectedLineIdx := *selectedLineIdxPtr
	newSelectedLineIdx := v.SelectedLineIdx()

	// we need to focus the view
	if err := lv.Gui.switchFocus(nil, v); err != nil {
		return err
	}

	if newSelectedLineIdx > lv.GetItemsLength()-1 {
		return lv.OnFocus()
	}

	*selectedLineIdxPtr = newSelectedLineIdx

	prevViewName := lv.Gui.currentViewName()
	if prevSelectedLineIdx == newSelectedLineIdx && prevViewName == lv.ViewName && lv.OnClickSelectedItem != nil {
		return lv.OnClickSelectedItem()
	}
	if lv.OnItemSelect != nil {
		return lv.OnItemSelect()
	}
	return nil
}

func (lv *listView) onSearchSelect(selectedLineIdx int) error {
	*lv.GetSelectedLineIdxPtr() = selectedLineIdx
	if lv.OnItemSelect != nil {
		return lv.OnItemSelect()
	}
	return nil
}

func (gui *Gui) menuListView() *listView {
	return &listView{
		ViewName:              "menu",
		GetItemsLength:        func() int { return gui.getMenuView().LinesHeight() },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Menu.SelectedLine },
		OnFocus:               gui.handleMenuSelect,
		OnItemSelect:          gui.handleMenuSelect,
		// need to add a layer of indirection here because the callback changes during runtime
		OnClickSelectedItem: func() error { return gui.State.Panels.Menu.OnPress(gui.g, nil) },
		Gui:                 gui,
		RendersToMainView:   false,
	}
}

func (gui *Gui) filesListView() *listView {
	return &listView{
		ViewName:              "files",
		GetItemsLength:        func() int { return len(gui.State.Files) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Files.SelectedLine },
		OnFocus:               gui.focusAndSelectFile,
		OnItemSelect:          gui.focusAndSelectFile,
		OnClickSelectedItem:   gui.handleFilePress,
		Gui:                   gui,
		RendersToMainView:     false,
	}
}

func (gui *Gui) branchesListView() *listView {
	return &listView{
		ViewName:              "branches",
		Context:               "local-branches",
		GetItemsLength:        func() int { return len(gui.State.Branches) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Branches.SelectedLine },
		OnFocus:               gui.handleBranchSelect,
		OnItemSelect:          gui.handleBranchSelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) remotesListView() *listView {
	return &listView{
		ViewName:              "branches",
		Context:               "remotes",
		GetItemsLength:        func() int { return len(gui.State.Remotes) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Remotes.SelectedLine },
		OnFocus:               gui.renderRemotesWithSelection,
		OnItemSelect:          gui.handleRemoteSelect,
		OnClickSelectedItem:   gui.handleRemoteEnter,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) remoteBranchesListView() *listView {
	return &listView{
		ViewName:              "branches",
		Context:               "remote-branches",
		GetItemsLength:        func() int { return len(gui.State.RemoteBranches) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.RemoteBranches.SelectedLine },
		OnFocus:               gui.handleRemoteBranchSelect,
		OnItemSelect:          gui.handleRemoteBranchSelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) tagsListView() *listView {
	return &listView{
		ViewName:              "branches",
		Context:               "tags",
		GetItemsLength:        func() int { return len(gui.State.Tags) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Tags.SelectedLine },
		OnFocus:               gui.handleTagSelect,
		OnItemSelect:          gui.handleTagSelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) branchCommitsListView() *listView {
	return &listView{
		ViewName:              "commits",
		Context:               "branch-commits",
		GetItemsLength:        func() int { return len(gui.State.Commits) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Commits.SelectedLine },
		OnFocus:               gui.handleCommitSelect,
		OnItemSelect:          gui.handleCommitSelect,
		OnClickSelectedItem:   gui.handleSwitchToCommitFilesPanel,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) reflogCommitsListView() *listView {
	return &listView{
		ViewName:              "commits",
		Context:               "reflog-commits",
		GetItemsLength:        func() int { return len(gui.State.FilteredReflogCommits) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.ReflogCommits.SelectedLine },
		OnFocus:               gui.handleReflogCommitSelect,
		OnItemSelect:          gui.handleReflogCommitSelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) stashListView() *listView {
	return &listView{
		ViewName:              "stash",
		GetItemsLength:        func() int { return len(gui.State.StashEntries) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Stash.SelectedLine },
		OnFocus:               gui.handleStashEntrySelect,
		OnItemSelect:          gui.handleStashEntrySelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) commitFilesListView() *listView {
	return &listView{
		ViewName:              "commitFiles",
		GetItemsLength:        func() int { return len(gui.State.CommitFiles) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.CommitFiles.SelectedLine },
		OnFocus:               gui.handleCommitFileSelect,
		OnItemSelect:          gui.handleCommitFileSelect,
		Gui:                   gui,
		RendersToMainView:     true,
	}
}

func (gui *Gui) getListViews() []*listView {
	return []*listView{
		gui.menuListView(),
		gui.filesListView(),
		gui.branchesListView(),
		gui.remotesListView(),
		gui.remoteBranchesListView(),
		gui.tagsListView(),
		gui.branchCommitsListView(),
		gui.reflogCommitsListView(),
		gui.stashListView(),
		gui.commitFilesListView(),
	}
}
