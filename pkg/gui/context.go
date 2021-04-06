package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

type ContextKind int

const (
	SIDE_CONTEXT ContextKind = iota
	MAIN_CONTEXT
	TEMPORARY_POPUP
	PERSISTENT_POPUP
)

type ContextKey string

const (
	STATUS_CONTEXT_KEY              ContextKey = "status"
	FILES_CONTEXT_KEY               ContextKey = "files"
	LOCAL_BRANCHES_CONTEXT_KEY      ContextKey = "localBranches"
	REMOTES_CONTEXT_KEY             ContextKey = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY     ContextKey = "remoteBranches"
	TAGS_CONTEXT_KEY                ContextKey = "tags"
	BRANCH_COMMITS_CONTEXT_KEY      ContextKey = "commits"
	REFLOG_COMMITS_CONTEXT_KEY      ContextKey = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY         ContextKey = "subCommits"
	COMMIT_FILES_CONTEXT_KEY        ContextKey = "commitFiles"
	STASH_CONTEXT_KEY               ContextKey = "stash"
	MAIN_NORMAL_CONTEXT_KEY         ContextKey = "normal"
	MAIN_MERGING_CONTEXT_KEY        ContextKey = "merging"
	MAIN_PATCH_BUILDING_CONTEXT_KEY ContextKey = "patchBuilding"
	MAIN_STAGING_CONTEXT_KEY        ContextKey = "staging"
	MENU_CONTEXT_KEY                ContextKey = "menu"
	CREDENTIALS_CONTEXT_KEY         ContextKey = "credentials"
	CONFIRMATION_CONTEXT_KEY        ContextKey = "confirmation"
	SEARCH_CONTEXT_KEY              ContextKey = "search"
	COMMIT_MESSAGE_CONTEXT_KEY      ContextKey = "commitMessage"
	SUBMODULES_CONTEXT_KEY          ContextKey = "submodules"
	SUGGESTIONS_CONTEXT_KEY         ContextKey = "suggestions"
)

var allContextKeys = []ContextKey{
	STATUS_CONTEXT_KEY,
	FILES_CONTEXT_KEY,
	LOCAL_BRANCHES_CONTEXT_KEY,
	REMOTES_CONTEXT_KEY,
	REMOTE_BRANCHES_CONTEXT_KEY,
	TAGS_CONTEXT_KEY,
	BRANCH_COMMITS_CONTEXT_KEY,
	REFLOG_COMMITS_CONTEXT_KEY,
	SUB_COMMITS_CONTEXT_KEY,
	COMMIT_FILES_CONTEXT_KEY,
	STASH_CONTEXT_KEY,
	MAIN_NORMAL_CONTEXT_KEY,
	MAIN_MERGING_CONTEXT_KEY,
	MAIN_PATCH_BUILDING_CONTEXT_KEY,
	MAIN_STAGING_CONTEXT_KEY,
	MENU_CONTEXT_KEY,
	CREDENTIALS_CONTEXT_KEY,
	CONFIRMATION_CONTEXT_KEY,
	SEARCH_CONTEXT_KEY,
	COMMIT_MESSAGE_CONTEXT_KEY,
	SUBMODULES_CONTEXT_KEY,
	SUGGESTIONS_CONTEXT_KEY,
}

type ContextTree struct {
	Status         Context
	Files          *ListContext
	Submodules     *ListContext
	Menu           *ListContext
	Branches       *ListContext
	Remotes        *ListContext
	RemoteBranches *ListContext
	Tags           *ListContext
	BranchCommits  *ListContext
	CommitFiles    *ListContext
	ReflogCommits  *ListContext
	SubCommits     *ListContext
	Stash          *ListContext
	Suggestions    *ListContext
	Normal         Context
	Staging        Context
	PatchBuilding  Context
	Merging        Context
	Credentials    Context
	Confirmation   Context
	CommitMessage  Context
	Search         Context
}

func (gui *Gui) allContexts() []Context {
	return []Context{
		gui.State.Contexts.Status,
		gui.State.Contexts.Files,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.Menu,
		gui.State.Contexts.Confirmation,
		gui.State.Contexts.Credentials,
		gui.State.Contexts.CommitMessage,
		gui.State.Contexts.Normal,
		gui.State.Contexts.Staging,
		gui.State.Contexts.Merging,
		gui.State.Contexts.PatchBuilding,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Suggestions,
	}
}

type Context interface {
	HandleFocus() error
	HandleFocusLost() error
	HandleRender() error
	GetKind() ContextKind
	GetViewName() string
	GetWindowName() string
	SetWindowName(string)
	GetKey() ContextKey
	SetParentContext(Context)

	// we return a bool here to tell us whether or not the returned value just wraps a nil
	GetParentContext() (Context, bool)
	GetOptionsMap() map[string]string
}

type BasicContext struct {
	OnFocus         func() error
	OnFocusLost     func() error
	OnRender        func() error
	OnGetOptionsMap func() map[string]string
	Kind            ContextKind
	Key             ContextKey
	ViewName        string
}

func (c BasicContext) GetOptionsMap() map[string]string {
	if c.OnGetOptionsMap != nil {
		return c.OnGetOptionsMap()
	}
	return nil
}

func (c BasicContext) SetWindowName(windowName string) {
	panic("can't set window name on basic context")
}

func (c BasicContext) GetWindowName() string {
	// TODO: fix this up
	return c.GetViewName()
}

func (c BasicContext) SetParentContext(Context) {
	panic("can't set parent context on basic context")
}

func (c BasicContext) GetParentContext() (Context, bool) {
	return nil, false
}

func (c BasicContext) HandleRender() error {
	if c.OnRender != nil {
		return c.OnRender()
	}
	return nil
}

func (c BasicContext) GetViewName() string {
	return c.ViewName
}

func (c BasicContext) HandleFocus() error {
	return c.OnFocus()
}

func (c BasicContext) HandleFocusLost() error {
	if c.OnFocusLost != nil {
		return c.OnFocusLost()
	}
	return nil
}

func (c BasicContext) GetKind() ContextKind {
	return c.Kind
}

func (c BasicContext) GetKey() ContextKey {
	return c.Key
}

func (gui *Gui) contextTree() ContextTree {
	return ContextTree{
		Status: BasicContext{
			OnFocus:  gui.handleStatusSelect,
			Kind:     SIDE_CONTEXT,
			ViewName: "status",
			Key:      STATUS_CONTEXT_KEY,
		},
		Files:          gui.filesListContext(),
		Submodules:     gui.submodulesListContext(),
		Menu:           gui.menuListContext(),
		Remotes:        gui.remotesListContext(),
		RemoteBranches: gui.remoteBranchesListContext(),
		BranchCommits:  gui.branchCommitsListContext(),
		CommitFiles:    gui.commitFilesListContext(),
		ReflogCommits:  gui.reflogCommitsListContext(),
		SubCommits:     gui.subCommitsListContext(),
		Branches:       gui.branchesListContext(),
		Tags:           gui.tagsListContext(),
		Stash:          gui.stashListContext(),
		Normal: BasicContext{
			OnFocus: func() error {
				return nil // TODO: should we do something here? We should allow for scrolling the panel
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_NORMAL_CONTEXT_KEY,
		},
		Staging: BasicContext{
			OnFocus: func() error {
				return nil
				// TODO: centralise the code here
				// return gui.refreshStagingPanel(false, -1)
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_STAGING_CONTEXT_KEY,
		},
		PatchBuilding: BasicContext{
			OnFocus: func() error {
				return nil
				// TODO: centralise the code here
				// return gui.refreshPatchBuildingPanel(-1)
			},
			Kind:     MAIN_CONTEXT,
			ViewName: "main",
			Key:      MAIN_PATCH_BUILDING_CONTEXT_KEY,
		},
		Merging: BasicContext{
			OnFocus:         gui.refreshMergePanelWithLock,
			Kind:            MAIN_CONTEXT,
			ViewName:        "main",
			Key:             MAIN_MERGING_CONTEXT_KEY,
			OnGetOptionsMap: gui.getMergingOptions,
		},
		Credentials: BasicContext{
			OnFocus:  gui.handleCredentialsViewFocused,
			Kind:     PERSISTENT_POPUP,
			ViewName: "credentials",
			Key:      CREDENTIALS_CONTEXT_KEY,
		},
		Confirmation: BasicContext{
			OnFocus:  func() error { return nil },
			Kind:     TEMPORARY_POPUP,
			ViewName: "confirmation",
			Key:      CONFIRMATION_CONTEXT_KEY,
		},
		Suggestions: gui.suggestionsListContext(),
		CommitMessage: BasicContext{
			OnFocus:  gui.handleCommitMessageFocused,
			Kind:     PERSISTENT_POPUP,
			ViewName: "commitMessage",
			Key:      COMMIT_MESSAGE_CONTEXT_KEY,
		},
		Search: BasicContext{
			OnFocus:  func() error { return nil },
			Kind:     PERSISTENT_POPUP,
			ViewName: "search",
			Key:      SEARCH_CONTEXT_KEY,
		},
	}
}

func (tree ContextTree) initialViewContextMap() map[string]Context {
	return map[string]Context{
		"status":        tree.Status,
		"files":         tree.Files,
		"branches":      tree.Branches,
		"commits":       tree.BranchCommits,
		"commitFiles":   tree.CommitFiles,
		"stash":         tree.Stash,
		"menu":          tree.Menu,
		"confirmation":  tree.Confirmation,
		"credentials":   tree.Credentials,
		"commitMessage": tree.CommitMessage,
		"main":          tree.Normal,
		"secondary":     tree.Normal,
	}
}

func (gui *Gui) popupViewNames() []string {
	result := []string{}
	for _, context := range gui.allContexts() {
		if context.GetKind() == PERSISTENT_POPUP || context.GetKind() == TEMPORARY_POPUP {
			result = append(result, context.GetViewName())
		}
	}

	return result
}

func (tree ContextTree) initialViewTabContextMap() map[string][]tabContext {
	return map[string][]tabContext{
		"branches": {
			{
				tab:      "Local Branches",
				contexts: []Context{tree.Branches},
			},
			{
				tab: "Remotes",
				contexts: []Context{
					tree.Remotes,
					tree.RemoteBranches,
				},
			},
			{
				tab:      "Tags",
				contexts: []Context{tree.Tags},
			},
		},
		"commits": {
			{
				tab:      "Commits",
				contexts: []Context{tree.BranchCommits},
			},
			{
				tab: "Reflog",
				contexts: []Context{
					tree.ReflogCommits,
				},
			},
		},
		"files": {
			{
				tab:      "Files",
				contexts: []Context{tree.Files},
			},
			{
				tab: "Submodules",
				contexts: []Context{
					tree.Submodules,
				},
			},
		},
	}
}

func (gui *Gui) currentContextKeyIgnoringPopups() ContextKey {
	gui.State.ContextManager.Lock()
	defer gui.State.ContextManager.Unlock()

	stack := gui.State.ContextManager.ContextStack

	for i := range stack {
		reversedIndex := len(stack) - 1 - i
		context := stack[reversedIndex]
		kind := stack[reversedIndex].GetKind()
		if kind != TEMPORARY_POPUP && kind != PERSISTENT_POPUP {
			return context.GetKey()
		}
	}

	return ""
}

// use replaceContext when you don't want to return to the original context upon
// hitting escape: you want to go that context's parent instead.
func (gui *Gui) replaceContext(c Context) error {
	gui.g.Update(func(*gocui.Gui) error {
		gui.State.ContextManager.Lock()
		defer gui.State.ContextManager.Unlock()

		if len(gui.State.ContextManager.ContextStack) == 0 {
			gui.State.ContextManager.ContextStack = []Context{c}
		} else {
			// replace the last item with the given item
			gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack[0:len(gui.State.ContextManager.ContextStack)-1], c)
		}

		return gui.activateContext(c)
	})

	return nil
}

func (gui *Gui) pushContext(c Context) error {
	gui.g.Update(func(*gocui.Gui) error {
		gui.State.ContextManager.Lock()
		defer gui.State.ContextManager.Unlock()

		return gui.pushContextDirect(c)
	})

	return nil
}

func (gui *Gui) pushContextDirect(c Context) error {
	// push onto stack
	// if we are switching to a side context, remove all other contexts in the stack
	if c.GetKind() == SIDE_CONTEXT {
		for _, stackContext := range gui.State.ContextManager.ContextStack {
			if stackContext.GetKey() != c.GetKey() {
				if err := gui.deactivateContext(stackContext); err != nil {
					return err
				}
			}
		}
		gui.State.ContextManager.ContextStack = []Context{c}
	} else if len(gui.State.ContextManager.ContextStack) == 0 || gui.State.ContextManager.ContextStack[len(gui.State.ContextManager.ContextStack)-1].GetKey() != c.GetKey() {
		// Do not append if the one at the end is the same context (e.g. opening a menu from a menu)
		// In that case we'll just close the menu entirely when the user hits escape.

		// TODO: think about other exceptional cases
		gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack, c)
	}

	return gui.activateContext(c)
}

// asynchronous code idea: functions return an error via a channel, when done

// pushContextWithView is to be used when you don't know which context you
// want to switch to: you only know the view that you want to switch to. It will
// look up the context currently active for that view and switch to that context
func (gui *Gui) pushContextWithView(viewName string) error {
	return gui.pushContext(gui.State.ViewContextMap[viewName])
}

func (gui *Gui) returnFromContext() error {
	gui.g.Update(func(*gocui.Gui) error {
		gui.State.ContextManager.Lock()
		defer gui.State.ContextManager.Unlock()

		if len(gui.State.ContextManager.ContextStack) == 1 {
			// cannot escape from bottommost context
			return nil
		}

		n := len(gui.State.ContextManager.ContextStack) - 1

		currentContext := gui.State.ContextManager.ContextStack[n]
		newContext := gui.State.ContextManager.ContextStack[n-1]

		gui.State.ContextManager.ContextStack = gui.State.ContextManager.ContextStack[:n]

		if err := gui.deactivateContext(currentContext); err != nil {
			return err
		}

		return gui.activateContext(newContext)
	})

	return nil
}

func (gui *Gui) deactivateContext(c Context) error {
	// if we are the kind of context that is sent to back upon deactivation, we should do that
	if c.GetKind() == TEMPORARY_POPUP || c.GetKind() == PERSISTENT_POPUP || c.GetKey() == COMMIT_FILES_CONTEXT_KEY {
		view, err := gui.g.View(c.GetViewName())
		if err == nil {
			view.Visible = false
		}
	}

	if err := c.HandleFocusLost(); err != nil {
		return err
	}

	return nil
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c Context) error {
	v, err := gui.g.View(c.GetViewName())
	if err != nil {
		return nil
	}

	if ContextKey(v.Context) != c.GetKey() {
		return nil
	}

	if err := c.HandleRender(); err != nil {
		return err
	}

	if gui.currentViewName() == c.GetViewName() {
		if err := c.HandleFocus(); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) activateContext(c Context) error {
	viewName := c.GetViewName()
	v, err := gui.g.View(viewName)
	// if view no longer exists, pop again
	// (note: this should never happen, unless we call this code before our views are initialised)
	if err != nil {
		return gui.returnFromContext()
	}
	originalViewContextKey := ContextKey(v.Context)

	// ensure that any other window for which this view was active is now set to the default for that window.
	gui.setViewAsActiveForWindow(v)

	if viewName == "main" {
		gui.changeMainViewsContext(c.GetKey())
	} else {
		gui.changeMainViewsContext(MAIN_NORMAL_CONTEXT_KEY)
	}

	gui.setViewTabForContext(c)

	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		// if view no longer exists, pop again
		return gui.returnFromContext()
	}

	v.Visible = true

	// if the new context's view was previously displaying another context, render the new context
	if originalViewContextKey != c.GetKey() {
		if err := c.HandleRender(); err != nil {
			return err
		}
	}

	v.Context = string(c.GetKey())

	gui.g.Cursor = v.Editable

	// render the options available for the current context at the bottom of the screen
	optionsMap := c.GetOptionsMap()
	if optionsMap == nil {
		optionsMap = gui.globalOptionsMap()
	}
	gui.renderOptionsMap(optionsMap)

	if err := c.HandleFocus(); err != nil {
		return err
	}

	// TODO: consider removing this and instead depending on the .Context field of views
	gui.State.ViewContextMap[c.GetViewName()] = c

	return nil
}

// // currently unused
// func (gui *Gui) renderContextStack() string {
// 	result := ""
// 	for _, context := range gui.State.ContextManager.ContextStack {
// 		result += string(context.GetKey()) + "\n"
// 	}
// 	return result
// }

func (gui *Gui) currentContext() Context {
	gui.State.ContextManager.Lock()
	defer gui.State.ContextManager.Unlock()

	if len(gui.State.ContextManager.ContextStack) == 0 {
		return gui.defaultSideContext()
	}

	return gui.State.ContextManager.ContextStack[len(gui.State.ContextManager.ContextStack)-1]
}

// the status panel is not yet a list context (and may never be), so this method is not
// quite the same as currentSideContext()
func (gui *Gui) currentSideListContext() *ListContext {
	context := gui.currentSideContext()
	listContext, ok := context.(*ListContext)
	if !ok {
		return nil
	}
	return listContext
}

func (gui *Gui) currentSideContext() Context {
	gui.State.ContextManager.Lock()
	defer gui.State.ContextManager.Unlock()

	stack := gui.State.ContextManager.ContextStack

	// on startup the stack can be empty so we'll return an empty string in that case
	if len(stack) == 0 {
		return gui.defaultSideContext()
	}

	// find the first context in the stack with the type of SIDE_CONTEXT
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() == SIDE_CONTEXT {
			return context
		}
	}

	return gui.defaultSideContext()
}

func (gui *Gui) defaultSideContext() Context {
	if gui.State.Modes.Filtering.Active() {
		return gui.State.Contexts.BranchCommits
	} else {
		return gui.State.Contexts.Files
	}
}

// remove the need to do this: always use a mapping
func (gui *Gui) setInitialViewContexts() {
	// arguably we should only have our ViewContextMap and we should do away with
	// contexts on views, or vice versa
	for viewName, context := range gui.State.ViewContextMap {
		// see if the view exists. If it does, set the context on it
		view, err := gui.g.View(viewName)
		if err != nil {
			continue
		}

		view.Context = string(context.GetKey())
	}
}

// getFocusLayout returns a manager function for when view gain and lose focus
func (gui *Gui) getFocusLayout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := gui.g.CurrentView()
		if err := gui.onViewFocusChange(); err != nil {
			return err
		}
		// for now we don't consider losing focus to a popup panel as actually losing focus
		if newView != previousView && !gui.isPopupPanel(newView.Name()) {
			if err := gui.onViewFocusLost(previousView, newView); err != nil {
				return err
			}

			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onViewFocusChange() error {
	gui.g.Mutexes.ViewsMutex.Lock()
	defer gui.g.Mutexes.ViewsMutex.Unlock()

	currentView := gui.g.CurrentView()
	for _, view := range gui.g.Views() {
		view.Highlight = view.Name() != "main" && view == currentView
	}
	return nil
}

func (gui *Gui) onViewFocusLost(oldView *gocui.View, newView *gocui.View) error {
	if oldView == nil {
		return nil
	}

	if oldView.IsSearching() && newView != gui.Views.Search {
		if err := gui.onSearchEscape(); err != nil {
			return err
		}
	}

	if oldView == gui.Views.CommitFiles && newView != gui.Views.Main && newView != gui.Views.Secondary && newView != gui.Views.Search {
		gui.resetWindowForView(gui.Views.CommitFiles)
		if err := gui.deactivateContext(gui.State.Contexts.CommitFiles); err != nil {
			return err
		}
	}

	return nil
}

// changeContext is a helper function for when we want to change a 'main' context
// which currently just means a context that affects both the main and secondary views
// other views can have their context changed directly but this function helps
// keep the main and secondary views in sync
func (gui *Gui) changeMainViewsContext(contextKey ContextKey) {
	if gui.State.MainContext == contextKey {
		return
	}

	switch contextKey {
	case MAIN_NORMAL_CONTEXT_KEY, MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY, MAIN_MERGING_CONTEXT_KEY:
		gui.Views.Main.Context = string(contextKey)
		gui.Views.Secondary.Context = string(contextKey)
	default:
		panic(fmt.Sprintf("unknown context for main: %s", contextKey))
	}

	gui.State.MainContext = contextKey
}

func (gui *Gui) viewTabNames(viewName string) []string {
	tabContexts := gui.State.ViewTabContextMap[viewName]

	if len(tabContexts) == 0 {
		return nil
	}

	result := make([]string, len(tabContexts))
	for i, tabContext := range tabContexts {
		result[i] = tabContext.tab
	}

	return result
}

func (gui *Gui) setViewTabForContext(c Context) {
	// search for the context in our map and if we find it, set the tab for the corresponding view
	tabContexts, ok := gui.State.ViewTabContextMap[c.GetViewName()]
	if !ok {
		return
	}

	for tabIndex, tabContext := range tabContexts {
		for _, context := range tabContext.contexts {
			if context.GetKey() == c.GetKey() {
				// get the view, set the tab
				v, err := gui.g.View(c.GetViewName())
				if err != nil {
					gui.Log.Error(err)
					return
				}
				v.TabIndex = tabIndex
				return
			}
		}
	}
}

type tabContext struct {
	tab      string
	contexts []Context
}

func (gui *Gui) mustContextForContextKey(contextKey ContextKey) Context {
	context, ok := gui.contextForContextKey(contextKey)

	if !ok {
		panic(fmt.Sprintf("context not found for key %s", contextKey))
	}

	return context
}

func (gui *Gui) contextForContextKey(contextKey ContextKey) (Context, bool) {
	for _, context := range gui.allContexts() {
		if context.GetKey() == contextKey {
			return context, true
		}
	}

	return nil, false
}

func (gui *Gui) rerenderView(view *gocui.View) error {
	contextKey := ContextKey(view.Context)
	context := gui.mustContextForContextKey(contextKey)

	return context.HandleRender()
}

// currently unused
// func (gui *Gui) getCurrentSideView() *gocui.View {
// 	currentSideContext := gui.currentSideContext()
// 	if currentSideContext == nil {
// 		return nil
// 	}

// 	view, _ := gui.g.View(currentSideContext.GetViewName())

// 	return view
// }

func (gui *Gui) getSideContextSelectedItemId() string {
	currentSideContext := gui.currentSideListContext()
	if currentSideContext == nil {
		return ""
	}

	item, ok := currentSideContext.GetSelectedItem()

	if ok {
		return item.ID()
	}

	return ""
}
