package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

type ListContext struct {
	ViewName            string
	ContextKey          string
	GetItemsLength      func() int
	GetDisplayStrings   func() [][]string
	OnFocus             func() error
	OnFocusLost         func() error
	OnClickSelectedItem func() error
	OnGetOptionsMap     func() map[string]string

	// the boolean here tells us whether the item is nil. This is needed because you can't work it out on the calling end once the pointer is wrapped in an interface (unless you want to use reflection)
	SelectedItem  func() (ListItem, bool)
	GetPanelState func() IListPanelState

	Gui                        *Gui
	ResetMainViewOriginOnFocus bool
	Kind                       int
	ParentContext              Context
	// we can't know on the calling end whether a Context is actually a nil value without reflection, so we're storing this flag here to tell us. There has got to be a better way around this.
	hasParent  bool
	WindowName string
}

type ListItem interface {
	// ID is a SHA when the item is a commit, a filename when the item is a file, 'stash@{4}' when it's a stash entry, 'my_branch' when it's a branch
	ID() string

	// Description is something we would show in a message e.g. '123as14: push blah' for a commit
	Description() string
}

func (lc *ListContext) GetSelectedItem() (ListItem, bool) {
	return lc.SelectedItem()
}

func (lc *ListContext) SetWindowName(windowName string) {
	lc.WindowName = windowName
}

func (lc *ListContext) GetWindowName() string {
	windowName := lc.WindowName

	if windowName != "" {
		return windowName
	}

	// TODO: actually set this for everything so we don't default to the view name
	return lc.ViewName
}

func (lc *ListContext) SetParentContext(c Context) {
	lc.ParentContext = c
	lc.hasParent = true
}

func (lc *ListContext) GetParentContext() (Context, bool) {
	return lc.ParentContext, lc.hasParent
}

func (lc *ListContext) GetSelectedItemId() string {
	item, ok := lc.SelectedItem()

	if !ok {
		return ""
	}

	return item.ID()
}

func (lc *ListContext) GetOptionsMap() map[string]string {
	if lc.OnGetOptionsMap != nil {
		return lc.OnGetOptionsMap()
	}
	return nil
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (lc *ListContext) OnRender() error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	if lc.GetDisplayStrings != nil {
		lc.Gui.refreshSelectedLine(lc.GetPanelState(), lc.GetItemsLength())
		lc.Gui.renderDisplayStrings(view, lc.GetDisplayStrings())
	}

	return nil
}

func (lc *ListContext) GetKey() string {
	return lc.ContextKey
}

func (lc *ListContext) GetKind() int {
	return lc.Kind
}

func (lc *ListContext) GetViewName() string {
	return lc.ViewName
}

func (lc *ListContext) HandleFocusLost() error {
	if lc.OnFocusLost != nil {
		return lc.OnFocusLost()
	}

	return nil
}

func (lc *ListContext) HandleFocus() error {
	if lc.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	view.FocusPoint(0, lc.GetPanelState().GetSelectedLineIdx())

	if lc.Gui.State.Modes.Diffing.Active() {
		return lc.Gui.renderDiff()
	}

	if lc.OnFocus != nil {
		return lc.OnFocus()
	}

	return nil
}

func (lc *ListContext) HandleRender() error {
	return lc.OnRender()
}

func (lc *ListContext) handlePrevLine(g *gocui.Gui, v *gocui.View) error {
	return lc.handleLineChange(-1)
}

func (lc *ListContext) handleNextLine(g *gocui.Gui, v *gocui.View) error {
	return lc.handleLineChange(1)
}

func (lc *ListContext) handleLineChange(change int) error {
	if !lc.Gui.isPopupPanel(lc.ViewName) && lc.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return err
	}

	lc.Gui.changeSelectedLine(lc.GetPanelState(), lc.GetItemsLength(), change)
	view.FocusPoint(0, lc.GetPanelState().GetSelectedLineIdx())

	if lc.ResetMainViewOriginOnFocus {
		if err := lc.Gui.resetOrigin(lc.Gui.getMainView()); err != nil {
			return err
		}
		if err := lc.Gui.resetOrigin(lc.Gui.getSecondaryView()); err != nil {
			return err
		}
	}

	return lc.HandleFocus()
}

func (lc *ListContext) handleNextPage(g *gocui.Gui, v *gocui.View) error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}
	_, height := view.Size()
	delta := height - 1
	if delta == 0 {
		delta = 1
	}
	return lc.handleLineChange(delta)
}

func (lc *ListContext) handleGotoTop(g *gocui.Gui, v *gocui.View) error {
	return lc.handleLineChange(-lc.GetItemsLength())
}

func (lc *ListContext) handleGotoBottom(g *gocui.Gui, v *gocui.View) error {
	return lc.handleLineChange(lc.GetItemsLength())
}

func (lc *ListContext) handlePrevPage(g *gocui.Gui, v *gocui.View) error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}
	_, height := view.Size()
	delta := height - 1
	if delta == 0 {
		delta = 1
	}
	return lc.handleLineChange(-delta)
}

func (lc *ListContext) handleClick(g *gocui.Gui, v *gocui.View) error {
	if !lc.Gui.isPopupPanel(lc.ViewName) && lc.Gui.popupPanelFocused() {
		return nil
	}

	prevSelectedLineIdx := lc.GetPanelState().GetSelectedLineIdx()
	newSelectedLineIdx := v.SelectedLineIdx()

	// we need to focus the view
	if err := lc.Gui.switchContext(lc); err != nil {
		return err
	}

	if newSelectedLineIdx > lc.GetItemsLength()-1 {
		return nil
	}

	lc.GetPanelState().SetSelectedLineIdx(newSelectedLineIdx)

	prevViewName := lc.Gui.currentViewName()
	if prevSelectedLineIdx == newSelectedLineIdx && prevViewName == lc.ViewName && lc.OnClickSelectedItem != nil {
		return lc.OnClickSelectedItem()
	}
	return lc.HandleFocus()
}

func (lc *ListContext) onSearchSelect(selectedLineIdx int) error {
	lc.GetPanelState().SetSelectedLineIdx(selectedLineIdx)
	return lc.HandleFocus()
}

func (gui *Gui) menuListContext() *ListContext {
	return &ListContext{
		ViewName:                   "menu",
		ContextKey:                 "menu",
		GetItemsLength:             func() int { return gui.getMenuView().LinesHeight() },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Menu },
		OnFocus:                    gui.handleMenuSelect,
		OnClickSelectedItem:        func() error { return gui.onMenuPress() },
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,
		Kind:                       PERSISTENT_POPUP,
		OnGetOptionsMap:            gui.getMenuOptions,

		// no GetDisplayStrings field because we do a custom render on menu creation
	}
}

func (gui *Gui) filesListContext() *ListContext {
	return &ListContext{
		ViewName:   "files",
		ContextKey: FILES_CONTEXT_KEY,
		GetItemsLength: func() int {
			if !gui.State.Panels.Files.ShowTree {
				return len(gui.State.Files)
			}
			return gui.State.FilesTree.Height()
		},
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Files },
		OnFocus:                    gui.focusAndSelectFile,
		OnClickSelectedItem:        gui.handleFilePress,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			if !gui.State.Panels.Files.ShowTree {
				return presentation.GetFileListDisplayStrings(gui.State.Files, gui.State.Modes.Diffing.Ref)
			}

			newSelectedFile, newSelectedDir := gui.getSelectedDirOrFile()
			list := gui.State.FilesTree.RenderAsList(newSelectedFile, newSelectedDir)
			res := [][]string{}

			for _, item := range list {
				res = append(res, []string{item})
			}

			return res
		},
		SelectedItem: func() (ListItem, bool) {
			item, _ := gui.getSelectedDirOrFile()
			return item, item != nil
		},
	}
}

func (gui *Gui) FilesTreeListContext() *ListContext {
	return &ListContext{
		ViewName:                   "filesTree",
		ContextKey:                 FILES_TREE_CONTEXT_KEY,
		GetItemsLength:             func() int { return gui.State.FilesTree.GetY() },
		GetPanelState:              nil,
		OnFocus:                    gui.focusAndSelectFile,
		OnClickSelectedItem:        gui.handleFilePress,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: false,
		Kind:                       PERSISTENT_POPUP,
		GetDisplayStrings: func() [][]string {
			return presentation.GetFileListDisplayStrings(gui.State.Files, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item, _ := gui.getSelectedDirOrFile()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchesListContext() *ListContext {
	return &ListContext{
		ViewName:                   "branches",
		ContextKey:                 LOCAL_BRANCHES_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.Branches) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Branches },
		OnFocus:                    gui.handleBranchSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) remotesListContext() *ListContext {
	return &ListContext{
		ViewName:                   "branches",
		ContextKey:                 REMOTES_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.Remotes) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Remotes },
		OnFocus:                    gui.handleRemoteSelect,
		OnClickSelectedItem:        gui.handleRemoteEnter,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemote()
			return item, item != nil
		},
	}
}

func (gui *Gui) remoteBranchesListContext() *ListContext {
	return &ListContext{
		ViewName:                   "branches",
		ContextKey:                 REMOTE_BRANCHES_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.RemoteBranches) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.RemoteBranches },
		OnFocus:                    gui.handleRemoteBranchSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedRemoteBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) tagsListContext() *ListContext {
	return &ListContext{
		ViewName:                   "branches",
		ContextKey:                 TAGS_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.Tags) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Tags },
		OnFocus:                    gui.handleTagSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Tags, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedTag()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchCommitsListContext() *ListContext {
	return &ListContext{
		ViewName:                   "commits",
		ContextKey:                 BRANCH_COMMITS_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.Commits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Commits },
		OnFocus:                    gui.handleCommitSelect,
		OnClickSelectedItem:        gui.handleViewCommitFiles,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetCommitListDisplayStrings(gui.State.Commits, gui.State.ScreenMode != SCREEN_NORMAL, gui.cherryPickedCommitShaMap(), gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedLocalCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) reflogCommitsListContext() *ListContext {
	return &ListContext{
		ViewName:                   "commits",
		ContextKey:                 REFLOG_COMMITS_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.FilteredReflogCommits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.ReflogCommits },
		OnFocus:                    gui.handleReflogCommitSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetReflogCommitListDisplayStrings(gui.State.FilteredReflogCommits, gui.State.ScreenMode != SCREEN_NORMAL, gui.cherryPickedCommitShaMap(), gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedReflogCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) subCommitsListContext() *ListContext {
	return &ListContext{
		ViewName:                   "branches",
		ContextKey:                 SUB_COMMITS_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.SubCommits) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.SubCommits },
		OnFocus:                    gui.handleSubCommitSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			gui.Log.Warn("getting display strings for sub commits")
			return presentation.GetCommitListDisplayStrings(gui.State.SubCommits, gui.State.ScreenMode != SCREEN_NORMAL, gui.cherryPickedCommitShaMap(), gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedSubCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) stashListContext() *ListContext {
	return &ListContext{
		ViewName:                   "stash",
		ContextKey:                 STASH_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.StashEntries) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.Stash },
		OnFocus:                    gui.handleStashEntrySelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedStashEntry()
			return item, item != nil
		},
	}
}

func (gui *Gui) commitFilesListContext() *ListContext {
	return &ListContext{
		ViewName:                   "commitFiles",
		WindowName:                 "commits",
		ContextKey:                 COMMIT_FILES_CONTEXT_KEY,
		GetItemsLength:             func() int { return len(gui.State.CommitFiles) },
		GetPanelState:              func() IListPanelState { return gui.State.Panels.CommitFiles },
		OnFocus:                    gui.handleCommitFileSelect,
		Gui:                        gui,
		ResetMainViewOriginOnFocus: true,
		Kind:                       SIDE_CONTEXT,
		GetDisplayStrings: func() [][]string {
			return presentation.GetCommitFileListDisplayStrings(gui.State.CommitFiles, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (ListItem, bool) {
			item := gui.getSelectedCommitFile()
			return item, item != nil
		},
	}
}

func (gui *Gui) getListContexts() []*ListContext {
	return []*ListContext{
		gui.menuListContext(),
		gui.filesListContext(),
		gui.branchesListContext(),
		gui.remotesListContext(),
		gui.remoteBranchesListContext(),
		gui.tagsListContext(),
		gui.branchCommitsListContext(),
		gui.reflogCommitsListContext(),
		gui.subCommitsListContext(),
		gui.stashListContext(),
		gui.commitFilesListContext(),
	}
}

func (gui *Gui) getListContextKeyBindings() []*Binding {
	bindings := make([]*Binding, 0)

	for _, listContext := range gui.getListContexts() {
		bindings = append(bindings, []*Binding{
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.prevItem-alt"), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.prevItem"), Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: listContext.handlePrevLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.nextItem-alt"), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.nextItem"), Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.prevPage"), Modifier: gocui.ModNone, Handler: listContext.handlePrevPage, Description: gui.Tr.SLocalize("prevPage")},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.nextPage"), Modifier: gocui.ModNone, Handler: listContext.handleNextPage, Description: gui.Tr.SLocalize("nextPage")},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gui.getKey("universal.gotoTop"), Modifier: gocui.ModNone, Handler: listContext.handleGotoTop, Description: gui.Tr.SLocalize("gotoTop")},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: listContext.handleNextLine},
			{ViewName: listContext.ViewName, Contexts: []string{listContext.ContextKey}, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: listContext.handleClick},
		}...)

		// the commits panel needs to lazyload things so it has a couple of its own handlers
		openSearchHandler := gui.handleOpenSearch
		gotoBottomHandler := listContext.handleGotoBottom
		if listContext.ViewName == "commits" {
			openSearchHandler = gui.handleOpenSearchForCommitsPanel
			gotoBottomHandler = gui.handleGotoBottomForCommitsPanel
		}

		bindings = append(bindings, []*Binding{
			{
				ViewName:    listContext.ViewName,
				Contexts:    []string{listContext.ContextKey},
				Key:         gui.getKey("universal.startSearch"),
				Handler:     openSearchHandler,
				Description: gui.Tr.SLocalize("startSearch"),
			},
			{
				ViewName:    listContext.ViewName,
				Contexts:    []string{listContext.ContextKey},
				Key:         gui.getKey("universal.gotoBottom"),
				Handler:     gotoBottomHandler,
				Description: gui.Tr.SLocalize("gotoBottom"),
			},
		}...)
	}

	return bindings
}
