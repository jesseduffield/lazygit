package gui

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) popupViewNames() []string {
	result := []string{}
	for _, context := range gui.allContexts() {
		if context.GetKind() == types.PERSISTENT_POPUP || context.GetKind() == types.TEMPORARY_POPUP {
			result = append(result, context.GetViewName())
		}
	}

	return result
}

func (gui *Gui) currentContextKeyIgnoringPopups() types.ContextKey {
	gui.State.ContextManager.RLock()
	defer gui.State.ContextManager.RUnlock()

	stack := gui.State.ContextManager.ContextStack

	for i := range stack {
		reversedIndex := len(stack) - 1 - i
		context := stack[reversedIndex]
		kind := stack[reversedIndex].GetKind()
		if kind != types.TEMPORARY_POPUP && kind != types.PERSISTENT_POPUP {
			return context.GetKey()
		}
	}

	return ""
}

// use replaceContext when you don't want to return to the original context upon
// hitting escape: you want to go that context's parent instead.
func (gui *Gui) replaceContext(c types.Context) error {
	gui.State.ContextManager.Lock()
	defer gui.State.ContextManager.Unlock()

	if len(gui.State.ContextManager.ContextStack) == 0 {
		gui.State.ContextManager.ContextStack = []types.Context{c}
	} else {
		// replace the last item with the given item
		gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack[0:len(gui.State.ContextManager.ContextStack)-1], c)
	}

	return gui.activateContext(c)
}

func (gui *Gui) pushContext(c types.Context, opts ...types.OnFocusOpts) error {
	// using triple dot but you should only ever pass one of these opt structs
	if len(opts) > 1 {
		return errors.New("cannot pass multiple opts to pushContext")
	}

	gui.State.ContextManager.Lock()

	// push onto stack
	// if we are switching to a side context, remove all other contexts in the stack
	if c.GetKind() == types.SIDE_CONTEXT {
		for _, stackContext := range gui.State.ContextManager.ContextStack {
			if stackContext.GetKey() != c.GetKey() {
				if err := gui.deactivateContext(stackContext); err != nil {
					gui.State.ContextManager.Unlock()
					return err
				}
			}
		}
		gui.State.ContextManager.ContextStack = []types.Context{c}
	} else if len(gui.State.ContextManager.ContextStack) == 0 || gui.currentContextWithoutLock().GetKey() != c.GetKey() {
		// Do not append if the one at the end is the same context (e.g. opening a menu from a menu)
		// In that case we'll just close the menu entirely when the user hits escape.

		// TODO: think about other exceptional cases
		gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack, c)
	}

	gui.State.ContextManager.Unlock()

	return gui.activateContext(c, opts...)
}

// asynchronous code idea: functions return an error via a channel, when done

// pushContextWithView is to be used when you don't know which context you
// want to switch to: you only know the view that you want to switch to. It will
// look up the context currently active for that view and switch to that context
func (gui *Gui) pushContextWithView(viewName string) error {
	return gui.c.PushContext(gui.State.ViewContextMap[viewName])
}

func (gui *Gui) returnFromContext() error {
	gui.State.ContextManager.Lock()

	if len(gui.State.ContextManager.ContextStack) == 1 {
		// cannot escape from bottommost context
		gui.State.ContextManager.Unlock()
		return nil
	}

	n := len(gui.State.ContextManager.ContextStack) - 1

	currentContext := gui.State.ContextManager.ContextStack[n]
	newContext := gui.State.ContextManager.ContextStack[n-1]

	gui.State.ContextManager.ContextStack = gui.State.ContextManager.ContextStack[:n]

	gui.State.ContextManager.Unlock()

	if err := gui.deactivateContext(currentContext); err != nil {
		return err
	}

	return gui.activateContext(newContext)
}

func (gui *Gui) deactivateContext(c types.Context) error {
	view, _ := gui.g.View(c.GetViewName())

	if view != nil && view.IsSearching() {
		if err := gui.onSearchEscape(); err != nil {
			return err
		}
	}

	// if we are the kind of context that is sent to back upon deactivation, we should do that
	if view != nil && (c.GetKind() == types.TEMPORARY_POPUP || c.GetKind() == types.PERSISTENT_POPUP || c.GetKey() == context.COMMIT_FILES_CONTEXT_KEY) {
		view.Visible = false
	}

	if err := c.HandleFocusLost(); err != nil {
		return err
	}

	return nil
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context) error {
	v, err := gui.g.View(c.GetViewName())
	if err != nil {
		return nil
	}

	if types.ContextKey(v.Context) != c.GetKey() {
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

func (gui *Gui) activateContext(c types.Context, opts ...types.OnFocusOpts) error {
	viewName := c.GetViewName()
	v, err := gui.g.View(viewName)
	if err != nil {
		return err
	}
	originalViewContextKey := types.ContextKey(v.Context)

	// ensure that any other window for which this view was active is now set to the default for that window.
	gui.setViewAsActiveForWindow(v)

	if viewName == "main" {
		gui.changeMainViewsContext(c.GetKey())
	} else {
		gui.changeMainViewsContext(context.MAIN_NORMAL_CONTEXT_KEY)
	}

	gui.setViewTabForContext(c)

	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		return err
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

	if err := c.HandleFocus(opts...); err != nil {
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

func (gui *Gui) currentContext() types.Context {
	gui.State.ContextManager.RLock()
	defer gui.State.ContextManager.RUnlock()

	return gui.currentContextWithoutLock()
}

func (gui *Gui) currentContextWithoutLock() types.Context {
	if len(gui.State.ContextManager.ContextStack) == 0 {
		return gui.defaultSideContext()
	}

	return gui.State.ContextManager.ContextStack[len(gui.State.ContextManager.ContextStack)-1]
}

// the status panel is not yet a list context (and may never be), so this method is not
// quite the same as currentSideContext()
func (gui *Gui) currentSideListContext() types.IListContext {
	context := gui.currentSideContext()
	listContext, ok := context.(types.IListContext)
	if !ok {
		return nil
	}
	return listContext
}

func (gui *Gui) currentSideContext() types.Context {
	gui.State.ContextManager.RLock()
	defer gui.State.ContextManager.RUnlock()

	stack := gui.State.ContextManager.ContextStack

	// on startup the stack can be empty so we'll return an empty string in that case
	if len(stack) == 0 {
		return gui.defaultSideContext()
	}

	// find the first context in the stack with the type of types.SIDE_CONTEXT
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() == types.SIDE_CONTEXT {
			return context
		}
	}

	return gui.defaultSideContext()
}

// static as opposed to popup
func (gui *Gui) currentStaticContext() types.Context {
	gui.State.ContextManager.RLock()
	defer gui.State.ContextManager.RUnlock()

	stack := gui.State.ContextManager.ContextStack

	if len(stack) == 0 {
		return gui.defaultSideContext()
	}

	// find the first context in the stack without a popup type
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() != types.TEMPORARY_POPUP && context.GetKind() != types.PERSISTENT_POPUP {
			return context
		}
	}

	return gui.defaultSideContext()
}

func (gui *Gui) defaultSideContext() types.Context {
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
		view.Highlight = view.Name() != "main" && view.Name() != "extras" && view == currentView
	}
	return nil
}

func (gui *Gui) onViewFocusLost(oldView *gocui.View, newView *gocui.View) error {
	if oldView == nil {
		return nil
	}

	_ = oldView.SetOriginX(0)

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
func (gui *Gui) changeMainViewsContext(contextKey types.ContextKey) {
	if gui.State.MainContext == contextKey {
		return
	}

	switch contextKey {
	case context.MAIN_NORMAL_CONTEXT_KEY, context.MAIN_PATCH_BUILDING_CONTEXT_KEY, context.MAIN_STAGING_CONTEXT_KEY, context.MAIN_MERGING_CONTEXT_KEY:
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
		result[i] = tabContext.Tab
	}

	return result
}

func (gui *Gui) setViewTabForContext(c types.Context) {
	// search for the context in our map and if we find it, set the tab for the corresponding view
	tabContexts, ok := gui.State.ViewTabContextMap[c.GetViewName()]
	if !ok {
		return
	}

	for tabIndex, tabContext := range tabContexts {
		for _, context := range tabContext.Contexts {
			if context.GetKey() == c.GetKey() {
				// get the view, set the tab
				v, err := gui.g.View(c.GetViewName())
				if err != nil {
					gui.c.Log.Error(err)
					return
				}
				v.TabIndex = tabIndex
				return
			}
		}
	}
}

func (gui *Gui) mustContextForContextKey(contextKey types.ContextKey) types.Context {
	context, ok := gui.contextForContextKey(contextKey)

	if !ok {
		panic(fmt.Sprintf("context not found for key %s", contextKey))
	}

	return context
}

func (gui *Gui) contextForContextKey(contextKey types.ContextKey) (types.Context, bool) {
	for _, context := range gui.allContexts() {
		if context.GetKey() == contextKey {
			return context, true
		}
	}

	return nil, false
}

func (gui *Gui) rerenderView(view *gocui.View) error {
	contextKey := types.ContextKey(view.Context)
	context := gui.mustContextForContextKey(contextKey)

	return context.HandleRender()
}

func (gui *Gui) getSideContextSelectedItemId() string {
	currentSideContext := gui.currentSideListContext()
	if currentSideContext == nil {
		return ""
	}

	return currentSideContext.GetSelectedItemId()
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

// currently unused
// func (gui *Gui) renderContextStack() string {
// 	result := ""
// 	for _, context := range gui.State.ContextManager.ContextStack {
// 		result += context.GetViewName() + "\n"
// 	}
// 	return result
// }
