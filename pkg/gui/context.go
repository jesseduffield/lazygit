package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

const (
	SIDE_CONTEXT int = iota
	MAIN_CONTEXT
	TEMPORARY_POPUP
	PERSISTENT_POPUP
)

func GetKindWrapper(k int) func() int { return func() int { return k } }

type Context interface {
	HandleFocus() error
	HandleFocusLost() error
	GetKind() int
	GetViewName() string
	GetKey() string
}

type BasicContext struct {
	OnFocus     func() error
	OnFocusLost func() error
	Kind        int
	Key         string
	ViewName    string
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

func (c BasicContext) GetKind() int {
	return c.Kind
}

func (c BasicContext) GetKey() string {
	return c.Key
}

type SimpleContextNode struct {
	Context Context
}

type RemotesContextNode struct {
	Context  Context
	Branches SimpleContextNode
}

type CommitsContextNode struct {
	Context Context
	Files   SimpleContextNode
}

type ContextTree struct {
	Status        SimpleContextNode
	Files         SimpleContextNode
	Menu          SimpleContextNode
	Branches      SimpleContextNode
	Remotes       RemotesContextNode
	Tags          SimpleContextNode
	BranchCommits CommitsContextNode
	ReflogCommits SimpleContextNode
	Stash         SimpleContextNode
	Normal        SimpleContextNode
	Staging       SimpleContextNode
	PatchBuilding SimpleContextNode
	Merging       SimpleContextNode
	Credentials   SimpleContextNode
	Confirmation  SimpleContextNode
	CommitMessage SimpleContextNode
	Search        SimpleContextNode
}

func (gui *Gui) switchContext(c Context) error {
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

		if c.GetViewName() == "main" {
			gui.changeMainViewsContext(c.GetKey())
		}

		return gui.activateContext(c)
	})

	return nil
}

// switchContextToView is to be used when you don't know which context you
// want to switch to: you only know the view that you want to switch to. It will
// look up the context currently active for that view and switch to that context
func (gui *Gui) switchContextToView(viewName string) error {
	return gui.switchContext(gui.State.ViewContextMap[viewName])
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
	if c.GetKind() == TEMPORARY_POPUP || c.GetKind() == PERSISTENT_POPUP {
		_, _ = gui.g.SetViewOnBottom(c.GetViewName())
	}

	if err := c.HandleFocusLost(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) activateContext(c Context) error {
	gui.Log.Warn(gui.renderContextStack())

	viewName := c.GetViewName()
	_, err := gui.g.View(viewName)
	// if view no longer exists, pop again
	if err != nil {
		return gui.returnFromContext()
	}

	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	if _, err := gui.g.SetViewOnTop(viewName); err != nil {
		return err
	}

	newView := gui.g.CurrentView()
	newView.Context = c.GetKey()

	gui.g.Cursor = newView.Editable

	// TODO: move this logic to the context
	if err := gui.renderPanelOptions(); err != nil {
		return err
	}

	if err := c.HandleFocus(); err != nil {
		return err
	}

	// TODO: consider removing this and instead depending on the .Context field of views
	gui.State.ViewContextMap[c.GetViewName()] = c

	return nil
}

func (gui *Gui) renderContextStack() string {
	result := ""
	for _, context := range gui.State.ContextStack {
		result += context.GetViewName() + "\n"
	}
	return result
}

func (gui *Gui) currentContextKey() string {
	// on startup the stack can be empty so we'll return an empty string in that case
	if len(gui.State.ContextStack) == 0 {
		return ""
	}

	return gui.State.ContextStack[len(gui.State.ContextStack)-1].GetKey()
}

func (gui *Gui) contextTree() ContextTree {
	return ContextTree{
		Status: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  gui.handleStatusSelect,
				Kind:     SIDE_CONTEXT,
				ViewName: "status",
				Key:      "status",
			},
		},
		Files: SimpleContextNode{
			Context: gui.filesListContext(),
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
		BranchCommits: CommitsContextNode{
			Context: gui.branchCommitsListContext(),
			Files: SimpleContextNode{
				Context: gui.commitFilesListContext(),
			},
		},
		ReflogCommits: SimpleContextNode{
			Context: gui.reflogCommitsListContext(),
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
				Key:      "normal",
			},
		},
		Staging: SimpleContextNode{
			Context: BasicContext{
				// TODO: think about different situations where this arises
				OnFocus: func() error {
					return nil
					// return gui.refreshStagingPanel(false, -1)
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      "staging",
			},
		},
		PatchBuilding: SimpleContextNode{
			Context: BasicContext{
				// TODO: think about different situations where this arises
				OnFocus: func() error {
					return gui.refreshPatchBuildingPanel(-1)
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      "patch-building",
			},
		},
		Merging: SimpleContextNode{
			Context: BasicContext{
				// TODO: think about different situations where this arises
				OnFocus: func() error {
					return gui.refreshMergePanel()
				},
				Kind:     MAIN_CONTEXT,
				ViewName: "main",
				Key:      "merging",
			},
		},
		Credentials: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return gui.handleCredentialsViewFocused() },
				Kind:     PERSISTENT_POPUP,
				ViewName: "credentials",
				Key:      "credentials",
			},
		},
		Confirmation: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return nil },
				Kind:     TEMPORARY_POPUP,
				ViewName: "confirmation",
				Key:      "confirmation",
			},
		},
		CommitMessage: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return gui.handleCommitMessageFocused() },
				Kind:     PERSISTENT_POPUP,
				ViewName: "commitMessage",
				Key:      "commit-message", // admittedly awkward to have view names in camelCase and contexts in kebab-case
			},
		},
		Search: SimpleContextNode{
			Context: BasicContext{
				OnFocus:  func() error { return nil },
				Kind:     PERSISTENT_POPUP,
				ViewName: "search",
				Key:      "search",
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
		"commitFiles":   gui.Contexts.BranchCommits.Files.Context,
		"stash":         gui.Contexts.Stash.Context,
		"menu":          gui.Contexts.Menu.Context,
		"confirmation":  gui.Contexts.Confirmation.Context,
		"credentials":   gui.Contexts.Credentials.Context,
		"commitMessage": gui.Contexts.CommitMessage.Context,
		"main":          gui.Contexts.Normal.Context,
		"secondary":     gui.Contexts.Normal.Context,
	}
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

			if err := gui.onViewFocus(newView); err != nil {
				return err
			}
			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onViewFocusChange() error {
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

	if v.Name() == "main" {
		// if we have lost focus to a first-class panel, we need to do some cleanup
		gui.changeMainViewsContext("normal")
	}

	gui.Log.Info(v.Name() + " focus lost")
	return nil
}

func (gui *Gui) onViewFocus(newView *gocui.View) error {
	gui.setViewAsActiveForWindow(newView.Name())

	return nil
}

// changeContext is a helper function for when we want to change a 'main' context
// which currently just means a context that affects both the main and secondary views
// other views can have their context changed directly but this function helps
// keep the main and secondary views in sync
func (gui *Gui) changeMainViewsContext(context string) {
	if gui.State.MainContext == context {
		return
	}

	switch context {
	case "normal", "patch-building", "staging", "merging":
		gui.getMainView().Context = context
		gui.getSecondaryView().Context = context
	default:
		panic(fmt.Sprintf("unknown context for main: %s", context))
	}

	gui.State.MainContext = context
}
