package gui

import "github.com/jesseduffield/gocui"

type ListView struct {
	ViewName              string
	Context               string
	GetItemsLength        func() int
	GetSelectedLineIdxPtr func() *int
	OnFocus               func() error
	OnFocusLost           func() error
	OnItemSelect          func() error
	OnClickSelectedItem   func() error
	Gui                   *Gui
	RendersToMainView     bool
	Kind                  int
	Key                   string
}

func (lv *ListView) GetKey() string {
	return lv.Key
}

func (lv *ListView) GetKind() int {
	return lv.Kind
}

func (lv *ListView) GetViewName() string {
	return lv.ViewName
}

func (lv *ListView) HandleFocusLost() error {
	if lv.OnFocusLost != nil {
		return lv.OnFocusLost()
	}

	return nil
}

func (lv *ListView) HandleFocus() error {
	return lv.OnFocus()
}

func (lv *ListView) handlePrevLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(-1)
}

func (lv *ListView) handleNextLine(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(1)
}

func (lv *ListView) handleLineChange(change int) error {
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

func (lv *ListView) handleNextPage(g *gocui.Gui, v *gocui.View) error {
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

func (lv *ListView) handleGotoTop(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(-lv.GetItemsLength())
}

func (lv *ListView) handleGotoBottom(g *gocui.Gui, v *gocui.View) error {
	return lv.handleLineChange(lv.GetItemsLength())
}

func (lv *ListView) handlePrevPage(g *gocui.Gui, v *gocui.View) error {
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

func (lv *ListView) handleClick(g *gocui.Gui, v *gocui.View) error {
	if !lv.Gui.isPopupPanel(lv.ViewName) && lv.Gui.popupPanelFocused() {
		return nil
	}

	selectedLineIdxPtr := lv.GetSelectedLineIdxPtr()
	prevSelectedLineIdx := *selectedLineIdxPtr
	newSelectedLineIdx := v.SelectedLineIdx()

	// we need to focus the view
	if err := lv.Gui.switchContext(lv); err != nil {
		return err
	}

	if newSelectedLineIdx > lv.GetItemsLength()-1 {
		return nil
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

func (lv *ListView) onSearchSelect(selectedLineIdx int) error {
	*lv.GetSelectedLineIdxPtr() = selectedLineIdx
	if lv.OnItemSelect != nil {
		return lv.OnItemSelect()
	}
	return nil
}

func (gui *Gui) menuListView() *ListView {
	return &ListView{
		ViewName:              "menu",
		GetItemsLength:        func() int { return gui.getMenuView().LinesHeight() },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Menu.SelectedLine },
		OnFocus:               gui.handleMenuSelect,
		OnItemSelect:          gui.handleMenuSelect,
		// need to add a layer of indirection here because the callback changes during runtime
		OnClickSelectedItem: func() error { return gui.State.Panels.Menu.OnPress(gui.g, nil) },
		Gui:                 gui,
		RendersToMainView:   false,
		Kind:                PERSISTENT_POPUP,
		Key:                 "menu",
	}
}

func (gui *Gui) filesListView() *ListView {
	return &ListView{
		ViewName:              "files",
		GetItemsLength:        func() int { return len(gui.State.Files) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Files.SelectedLine },
		OnFocus:               gui.focusAndSelectFile,
		OnItemSelect:          gui.focusAndSelectFile,
		OnClickSelectedItem:   gui.handleFilePress,
		Gui:                   gui,
		RendersToMainView:     false,
		Kind:                  SIDE_CONTEXT,
		Key:                   "files",
	}
}

func (gui *Gui) branchesListView() *ListView {
	return &ListView{
		ViewName:              "branches",
		Context:               "local-branches",
		GetItemsLength:        func() int { return len(gui.State.Branches) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Branches.SelectedLine },
		OnFocus:               gui.handleBranchSelect,
		OnItemSelect:          gui.handleBranchSelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
		Key:                   "menu",
	}
}

func (gui *Gui) remotesListView() *ListView {
	return &ListView{
		ViewName:              "branches",
		Context:               "remotes",
		GetItemsLength:        func() int { return len(gui.State.Remotes) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Remotes.SelectedLine },
		OnFocus:               gui.renderRemotesWithSelection,
		OnItemSelect:          gui.handleRemoteSelect,
		OnClickSelectedItem:   gui.handleRemoteEnter,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) remoteBranchesListView() *ListView {
	return &ListView{
		ViewName:              "branches",
		Context:               "remote-branches",
		GetItemsLength:        func() int { return len(gui.State.RemoteBranches) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.RemoteBranches.SelectedLine },
		OnFocus:               gui.handleRemoteBranchSelect,
		OnItemSelect:          gui.handleRemoteBranchSelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) tagsListView() *ListView {
	return &ListView{
		ViewName:              "branches",
		Context:               "tags",
		GetItemsLength:        func() int { return len(gui.State.Tags) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Tags.SelectedLine },
		OnFocus:               gui.handleTagSelect,
		OnItemSelect:          gui.handleTagSelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) branchCommitsListView() *ListView {
	return &ListView{
		ViewName:              "commits",
		Context:               "branch-commits",
		GetItemsLength:        func() int { return len(gui.State.Commits) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Commits.SelectedLine },
		OnFocus:               gui.handleCommitSelect,
		OnItemSelect:          gui.handleCommitSelect,
		OnClickSelectedItem:   gui.handleSwitchToCommitFilesPanel,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) reflogCommitsListView() *ListView {
	return &ListView{
		ViewName:              "commits",
		Context:               "reflog-commits",
		GetItemsLength:        func() int { return len(gui.State.FilteredReflogCommits) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.ReflogCommits.SelectedLine },
		OnFocus:               gui.handleReflogCommitSelect,
		OnItemSelect:          gui.handleReflogCommitSelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) stashListView() *ListView {
	return &ListView{
		ViewName:              "stash",
		GetItemsLength:        func() int { return len(gui.State.StashEntries) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.Stash.SelectedLine },
		OnFocus:               gui.handleStashEntrySelect,
		OnItemSelect:          gui.handleStashEntrySelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) commitFilesListView() *ListView {
	return &ListView{
		ViewName:              "commitFiles",
		GetItemsLength:        func() int { return len(gui.State.CommitFiles) },
		GetSelectedLineIdxPtr: func() *int { return &gui.State.Panels.CommitFiles.SelectedLine },
		OnFocus:               gui.handleCommitFileSelect,
		OnItemSelect:          gui.handleCommitFileSelect,
		Gui:                   gui,
		RendersToMainView:     true,
		Kind:                  SIDE_CONTEXT,
	}
}

func (gui *Gui) getListViews() []*ListView {
	return []*ListView{
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
