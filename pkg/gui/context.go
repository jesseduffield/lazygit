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

const (
	STATUS_CONTEXT_KEY              = "status"
	FILES_CONTEXT_KEY               = "files"
	LOCAL_BRANCHES_CONTEXT_KEY      = "localBranches"
	REMOTES_CONTEXT_KEY             = "remotes"
	REMOTE_BRANCHES_CONTEXT_KEY     = "remoteBranches"
	TAGS_CONTEXT_KEY                = "tags"
	BRANCH_COMMITS_CONTEXT_KEY      = "commits"
	REFLOG_COMMITS_CONTEXT_KEY      = "reflogCommits"
	SUB_COMMITS_CONTEXT_KEY         = "subCommits"
	COMMIT_FILES_CONTEXT_KEY        = "commitFiles"
	STASH_CONTEXT_KEY               = "stash"
	MAIN_NORMAL_CONTEXT_KEY         = "normal"
	MAIN_MERGING_CONTEXT_KEY        = "merging"
	MAIN_PATCH_BUILDING_CONTEXT_KEY = "patchBuilding"
	MAIN_STAGING_CONTEXT_KEY        = "staging"
	MENU_CONTEXT_KEY                = "menu"
	CREDENTIALS_CONTEXT_KEY         = "credentials"
	CONFIRMATION_CONTEXT_KEY        = "confirmation"
	SEARCH_CONTEXT_KEY              = "search"
	COMMIT_MESSAGE_CONTEXT_KEY      = "commitMessage"
	SUBMODULES_CONTEXT_KEY          = "submodules"
	SUGGESTIONS_CONTEXT_KEY         = "suggestions"
)

var allContextKeys = []string{
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

type SimpleContextNode struct {
	Context Context
}

type RemotesContextNode struct {
	Context  Context
	Branches SimpleContextNode
}

type ContextTree struct {
	Status        SimpleContextNode
	Files         SimpleContextNode
	Submodules    SimpleContextNode
	Menu          SimpleContextNode
	Branches      SimpleContextNode
	Remotes       RemotesContextNode
	Tags          SimpleContextNode
	BranchCommits SimpleContextNode
	CommitFiles   SimpleContextNode
	ReflogCommits SimpleContextNode
	SubCommits    SimpleContextNode
	Stash         SimpleContextNode
	Normal        SimpleContextNode
	Staging       SimpleContextNode
	PatchBuilding SimpleContextNode
	Merging       SimpleContextNode
	Credentials   SimpleContextNode
	Confirmation  SimpleContextNode
	CommitMessage SimpleContextNode
	Search        SimpleContextNode
	Suggestions   SimpleContextNode
}

func (gui *Gui) allContexts() []Context {
	return []Context{
		gui.Contexts.Status.Context,
		gui.Contexts.Files.Context,
		gui.Contexts.Submodules.Context,
		gui.Contexts.Branches.Context,
		gui.Contexts.Remotes.Context,
		gui.Contexts.Remotes.Branches.Context,
		gui.Contexts.Tags.Context,
		gui.Contexts.BranchCommits.Context,
		gui.Contexts.CommitFiles.Context,
		gui.Contexts.ReflogCommits.Context,
		gui.Contexts.Stash.Context,
		gui.Contexts.Menu.Context,
		gui.Contexts.Confirmation.Context,
		gui.Contexts.Credentials.Context,
		gui.Contexts.CommitMessage.Context,
		gui.Contexts.Normal.Context,
		gui.Contexts.Staging.Context,
		gui.Contexts.Merging.Context,
		gui.Contexts.PatchBuilding.Context,
		gui.Contexts.SubCommits.Context,
		gui.Contexts.Suggestions.Context,
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
	GetKey() string
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
	Key             string
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

func (c BasicContext) GetKey() string {
	return c.Key
}

func (gui *Gui) contextTree() ContextTree {
	return ContextTree{
		Status: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  gui.handleStatusSelect,
				Kind:     SIDE_CONTEXT,
				ViewName: "status",
				Key:      STATUS_CONTEXT_KEY,
			},
		},
		Files: SimpleContextNode{
			Context: gui.filesListContext(),
		},
		Submodules: SimpleContextNode{
			Context: gui.submodulesListContext(),
		},
		Menu: SimpleContextNode{
			Context: gui.menuListContext(),
		},
		Remotes: RemotesContextNode{
			Context: gui.remotesListContext(),
			Branches: SimpleContextNode{
				Context: gui.remoteBranchesListContext(),
			},
		},
		BranchCommits: SimpleContextNode{
			Context: gui.branchCommitsListContext(),
		},
		CommitFiles: SimpleContextNode{
			Context: gui.commitFilesListContext(),
		},
		ReflogCommits: SimpleContextNode{
			Context: gui.reflogCommitsListContext(),
		},
		SubCommits: SimpleContextNode{
			Context: gui.subCommitsListContext(),
		},
		Branches: SimpleContextNode{
			Context: gui.branchesListContext(),
		},
		Tags: SimpleContextNode{
			Context: gui.tagsListContext(),
		},
		Stash: SimpleContextNode{
			Context: gui.stashListContext(),
		},
		Normal: SimpleContextNode{
			Context: BasicContext{
				OnFocus: func() error {
					return nil // TODO: should we do something here? We should allow for scrolling the panel
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      MAIN_NORMAL_CONTEXT_KEY,
			},
		},
		Staging: SimpleContextNode{
			Context: BasicContext{
				OnFocus: func() error {
					return nil
					// TODO: centralise the code here
					// return gui.refreshStagingPanel(false, -1)
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      MAIN_STAGING_CONTEXT_KEY,
			},
		},
		PatchBuilding: SimpleContextNode{
			Context: BasicContext{
				OnFocus: func() error {
					return nil
					// TODO: centralise the code here
					// return gui.refreshPatchBuildingPanel(-1)
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      MAIN_PATCH_BUILDING_CONTEXT_KEY,
			},
		},
		Merging: SimpleContextNode{
			Context: BasicContext{
				OnFocus:         gui.refreshMergePanelWithLock,
				Kind:            MAIN_CONTEXT,
				ViewName:        "main",
				Key:             MAIN_MERGING_CONTEXT_KEY,
				OnGetOptionsMap: gui.getMergingOptions,
			},
		},
		Credentials: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  gui.handleCredentialsViewFocused,
				Kind:     PERSISTENT_POPUP,
				ViewName: "credentials",
				Key:      CREDENTIALS_CONTEXT_KEY,
			},
		},
		Confirmation: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return nil },
				Kind:     TEMPORARY_POPUP,
				ViewName: "confirmation",
				Key:      CONFIRMATION_CONTEXT_KEY,
			},
		},
		Suggestions: SimpleContextNode{
			Context: gui.suggestionsListContext(),
		},
		CommitMessage: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  gui.handleCommitMessageFocused,
				Kind:     PERSISTENT_POPUP,
				ViewName: "commitMessage",
				Key:      COMMIT_MESSAGE_CONTEXT_KEY,
			},
		},
		Search: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return nil },
				Kind:     PERSISTENT_POPUP,
				ViewName: "search",
				Key:      SEARCH_CONTEXT_KEY,
			},
		},
	}
}

func (gui *Gui) initialViewContextMap() map[string]Context {
	return map[string]Context{
		"status":        gui.Contexts.Status.Context,
		"files":         gui.Contexts.Files.Context,
		"branches":      gui.Contexts.Branches.Context,
		"commits":       gui.Contexts.BranchCommits.Context,
		"commitFiles":   gui.Contexts.CommitFiles.Context,
		"stash":         gui.Contexts.Stash.Context,
		"menu":          gui.Contexts.Menu.Context,
		"confirmation":  gui.Contexts.Confirmation.Context,
		"credentials":   gui.Contexts.Credentials.Context,
		"commitMessage": gui.Contexts.CommitMessage.Context,
		"main":          gui.Contexts.Normal.Context,
		"secondary":     gui.Contexts.Normal.Context,
	}
}

func (gui *Gui) viewTabContextMap() map[string][]tabContext {
	return map[string][]tabContext{
		"branches": {
			{
				tab:      "Local Branches",
				contexts: []Context{gui.Contexts.Branches.Context},
			},
			{
				tab: "Remotes",
				contexts: []Context{
					gui.Contexts.Remotes.Context,
					gui.Contexts.Remotes.Branches.Context,
				},
			},
			{
				tab:      "Tags",
				contexts: []Context{gui.Contexts.Tags.Context},
			},
		},
		"commits": {
			{
				tab:      "Commits",
				contexts: []Context{gui.Contexts.BranchCommits.Context},
			},
			{
				tab: "Reflog",
				contexts: []Context{
					gui.Contexts.ReflogCommits.Context,
				},
			},
		},
		"files": {
			{
				tab:      "Files",
				contexts: []Context{gui.Contexts.Files.Context},
			},
			{
				tab: "Submodules",
				contexts: []Context{
					gui.Contexts.Submodules.Context,
				},
			},
		},
	}
}

func (gui *Gui) currentContextKeyIgnoringPopups() string {
	stack := gui.State.ContextStack

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
		if len(gui.State.ContextStack) == 0 {
			gui.State.ContextStack = []Context{c}
		} else {
			// replace the last item with the given item
			gui.State.ContextStack = append(gui.State.ContextStack[0:len(gui.State.ContextStack)-1], c)
		}

		return gui.activateContext(c)
	})

	return nil
}

func (gui *Gui) pushContext(c Context) error {
	gui.g.Update(func(*gocui.Gui) error {
		// push onto stack
		// if we are switching to a side context, remove all other contexts in the stack
		if c.GetKind() == SIDE_CONTEXT {
			for _, stackContext := range gui.State.ContextStack {
				if stackContext.GetKey() != c.GetKey() {
					if err := gui.deactivateContext(stackContext); err != nil {
						return err
					}
				}
			}
			gui.State.ContextStack = []Context{c}
		} else {
			// TODO: think about other exceptional cases
			gui.State.ContextStack = append(gui.State.ContextStack, c)
		}

		return gui.activateContext(c)
	})

	return nil
}

// pushContextWithView is to be used when you don't know which context you
// want to switch to: you only know the view that you want to switch to. It will
// look up the context currently active for that view and switch to that context
func (gui *Gui) pushContextWithView(viewName string) error {
	return gui.pushContext(gui.State.ViewContextMap[viewName])
}

func (gui *Gui) returnFromContext() error {
	gui.g.Update(func(*gocui.Gui) error {
		// TODO: add mutexes

		if len(gui.State.ContextStack) == 1 {
			// cannot escape from bottommost context
			return nil
		}

		n := len(gui.State.ContextStack) - 1

		currentContext := gui.State.ContextStack[n]
		newContext := gui.State.ContextStack[n-1]

		gui.State.ContextStack = gui.State.ContextStack[:n]

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
		_, _ = gui.g.SetViewOnBottom(c.GetViewName())
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

	if v.Context != c.GetKey() {
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
	if err != nil {
		return gui.returnFromContext()
	}
	originalViewContextKey := v.Context

	// ensure that any other window for which this view was active is now set to the default for that window.
	gui.setViewAsActiveForWindow(viewName)

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

	if _, err := gui.g.SetViewOnTop(viewName); err != nil {
		// if view no longer exists, pop again
		return gui.returnFromContext()
	}

	// if the new context's view was previously displaying another context, render the new context
	if originalViewContextKey != c.GetKey() {
		if err := c.HandleRender(); err != nil {
			return err
		}
	}

	v.Context = c.GetKey()

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

// currently unused
// func (gui *Gui) renderContextStack() string {
// 	result := ""
// 	for _, context := range gui.State.ContextStack {
// 		result += context.GetKey() + "\n"
// 	}
// 	return result
// }

func (gui *Gui) currentContext() Context {
	if len(gui.State.ContextStack) == 0 {
		return gui.defaultSideContext()
	}

	return gui.State.ContextStack[len(gui.State.ContextStack)-1]
}

func (gui *Gui) currentSideContext() *ListContext {
	stack := gui.State.ContextStack

	// on startup the stack can be empty so we'll return an empty string in that case
	if len(stack) == 0 {
		return nil
	}

	// find the first context in the stack with the type of SIDE_CONTEXT
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() == SIDE_CONTEXT {
			// not all side contexts are list contexts (e.g. the status panel)
			listContext, ok := context.(*ListContext)
			if !ok {
				return nil
			}
			return listContext
		}
	}

	return nil
}

func (gui *Gui) defaultSideContext() Context {
	return gui.Contexts.Files.Context
}

func (gui *Gui) setInitialViewContexts() {
	// arguably we should only have our ViewContextMap and we should do away with
	// contexts on views, or vice versa
	for viewName, context := range gui.State.ViewContextMap {
		// see if the view exists. If it does, set the context on it
		view, err := gui.g.View(viewName)
		if err != nil {
			continue
		}

		view.Context = context.GetKey()
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

func (gui *Gui) onViewFocusLost(v *gocui.View, newView *gocui.View) error {
	if v == nil {
		return nil
	}

	if v.IsSearching() && newView.Name() != "search" {
		if err := gui.onSearchEscape(); err != nil {
			return err
		}
	}

	if v.Name() == "commitFiles" && newView.Name() != "main" && newView.Name() != "secondary" {
		gui.resetWindowForView("commitFiles")
		if err := gui.deactivateContext(gui.Contexts.CommitFiles.Context); err != nil {
			return err
		}
	}

	return nil
}

// changeContext is a helper function for when we want to change a 'main' context
// which currently just means a context that affects both the main and secondary views
// other views can have their context changed directly but this function helps
// keep the main and secondary views in sync
func (gui *Gui) changeMainViewsContext(contextKey string) {
	if gui.State.MainContext == contextKey {
		return
	}

	switch contextKey {
	case MAIN_NORMAL_CONTEXT_KEY, MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY, MAIN_MERGING_CONTEXT_KEY:
		gui.getMainView().Context = contextKey
		gui.getSecondaryView().Context = contextKey
	default:
		panic(fmt.Sprintf("unknown context for main: %s", contextKey))
	}

	gui.State.MainContext = contextKey
}

func (gui *Gui) viewTabNames(viewName string) []string {
	tabContexts := gui.ViewTabContextMap[viewName]

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
	tabContexts, ok := gui.ViewTabContextMap[c.GetViewName()]
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

func (gui *Gui) mustContextForContextKey(contextKey string) Context {
	context, ok := gui.contextForContextKey(contextKey)

	if !ok {
		panic(fmt.Sprintf("context not found for key %s", contextKey))
	}

	return context
}

func (gui *Gui) contextForContextKey(contextKey string) (Context, bool) {
	for _, context := range gui.allContexts() {
		if context.GetKey() == contextKey {
			return context, true
		}
	}

	return nil, false
}

func (gui *Gui) rerenderView(viewName string) error {
	v, err := gui.g.View(viewName)
	if err != nil {
		return nil
	}

	contextKey := v.Context
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
	currentSideContext := gui.currentSideContext()
	if currentSideContext == nil {
		return ""
	}

	item, ok := currentSideContext.GetSelectedItem()

	if ok {
		return item.ID()
	}

	return ""
}
