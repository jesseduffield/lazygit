package gui

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This file is for the management of contexts. There is a context stack such that
// for example you might start off in the commits context and then open a menu, putting
// you in the menu context. When contexts are activated/deactivated certain things need
// to happen like showing/hiding views and rendering content.

func (gui *Gui) popupViewNames() []string {
	popups := slices.Filter(gui.State.Contexts.Flatten(), func(c types.Context) bool {
		return c.GetKind() == types.PERSISTENT_POPUP || c.GetKind() == types.TEMPORARY_POPUP
	})

	return slices.Map(popups, func(c types.Context) string {
		return c.GetViewName()
	})
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
	if !c.IsFocusable() {
		return nil
	}

	gui.State.ContextManager.Lock()

	if len(gui.State.ContextManager.ContextStack) == 0 {
		gui.State.ContextManager.ContextStack = []types.Context{c}
	} else {
		// replace the last item with the given item
		gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack[0:len(gui.State.ContextManager.ContextStack)-1], c)
	}

	defer gui.State.ContextManager.Unlock()

	return gui.activateContext(c)
}

func (gui *Gui) pushContext(c types.Context, opts ...types.OnFocusOpts) error {
	// using triple dot but you should only ever pass one of these opt structs
	if len(opts) > 1 {
		return errors.New("cannot pass multiple opts to pushContext")
	}

	if !c.IsFocusable() {
		return nil
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

// pushContextWithView is to be used when you don't know which context you
// want to switch to: you only know the view that you want to switch to. It will
// look up the context currently active for that view and switch to that context
func (gui *Gui) pushContextWithView(viewName string) error {
	return gui.c.PushContext(gui.State.ViewContextMap.Get(viewName))
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

	gui.g.SetCurrentContext(string(newContext.GetKey()))

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
	if view != nil &&
		(c.GetKind() == types.TEMPORARY_POPUP ||
			c.GetKind() == types.PERSISTENT_POPUP) {
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
	if gui.State.ViewContextMap.Get(c.GetViewName()).GetKey() != c.GetKey() {
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

	originalViewContext := gui.State.ViewContextMap.Get(viewName)
	var originalViewContextKey types.ContextKey = ""
	if originalViewContext != nil {
		originalViewContextKey = originalViewContext.GetKey()
	}

	gui.setWindowContext(c)
	gui.setViewTabForContext(c)

	if viewName == "main" {
		gui.changeMainViewsContext(c)
	} else {
		gui.changeMainViewsContext(gui.State.Contexts.Normal)
	}

	gui.g.SetCurrentContext(string(c.GetKey()))
	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	desiredTitle := c.Title()
	if desiredTitle != "" {
		v.Title = desiredTitle
	}

	v.Visible = true

	// if the new context's view was previously displaying another context, render the new context
	if originalViewContextKey != c.GetKey() {
		if err := c.HandleRender(); err != nil {
			return err
		}
	}

	gui.ViewContextMapSet(viewName, c)

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

	return nil
}

func (gui *Gui) optionsMapToString(optionsMap map[string]string) string {
	options := maps.MapToSlice(optionsMap, func(key string, description string) string {
		return key + ": " + description
	})
	sort.Strings(options)
	return strings.Join(options, ", ")
}

func (gui *Gui) renderOptionsMap(optionsMap map[string]string) {
	_ = gui.renderString(gui.Views.Options, gui.optionsMapToString(optionsMap))
}

// also setting context on view for now. We'll need to pick one of these two approaches to stick with.
func (gui *Gui) ViewContextMapSet(viewName string, c types.Context) {
	gui.State.ViewContextMap.Set(viewName, c)
	view, err := gui.g.View(viewName)
	if err != nil {
		panic(err)
	}
	view.Context = string(c.GetKey())
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
		return gui.State.Contexts.LocalCommits
	} else {
		return gui.State.Contexts.Files
	}
}

// getFocusLayout returns a manager function for when view gain and lose focus
func (gui *Gui) getFocusLayout() func(g *gocui.Gui) error {
	var previousView *gocui.View
	return func(g *gocui.Gui) error {
		newView := gui.g.CurrentView()
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

func (gui *Gui) onViewFocusLost(oldView *gocui.View, newView *gocui.View) error {
	if oldView == nil {
		return nil
	}

	oldView.Highlight = false

	_ = oldView.SetOriginX(0)

	return nil
}

func (gui *Gui) TransientContexts() []types.Context {
	return slices.Filter(gui.State.Contexts.Flatten(), func(context types.Context) bool {
		return context.IsTransient()
	})
}

// changeContext is a helper function for when we want to change a 'main' context
// which currently just means a context that affects both the main and secondary views
// other views can have their context changed directly but this function helps
// keep the main and secondary views in sync
func (gui *Gui) changeMainViewsContext(c types.Context) {
	if gui.State.MainContext == c.GetKey() {
		return
	}

	switch c.GetKey() {
	case context.MAIN_NORMAL_CONTEXT_KEY, context.MAIN_PATCH_BUILDING_CONTEXT_KEY, context.MAIN_STAGING_CONTEXT_KEY, context.MAIN_MERGING_CONTEXT_KEY:
		gui.ViewContextMapSet(gui.Views.Main.Name(), c)
		gui.ViewContextMapSet(gui.Views.Secondary.Name(), c)
	default:
		panic(fmt.Sprintf("unknown context for main: %s", c.GetKey()))
	}

	gui.State.MainContext = c.GetKey()
}

func (gui *Gui) viewTabNames(viewName string) []string {
	tabContexts := gui.State.ViewTabContextMap[viewName]

	return slices.Map(tabContexts, func(tabContext context.TabContext) string {
		return tabContext.Tab
	})
}

func (gui *Gui) setViewTabForContext(c types.Context) {
	// search for the context in our map and if we find it, set the tab for the corresponding view
	tabContexts, ok := gui.State.ViewTabContextMap[c.GetViewName()]
	if !ok {
		return
	}

	for tabIndex, tabContext := range tabContexts {
		if tabContext.Context.GetKey() == c.GetKey() {
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

func (gui *Gui) rerenderView(view *gocui.View) error {
	return gui.State.ViewContextMap.Get(view.Name()).HandleRender()
}

func (gui *Gui) getSideContextSelectedItemId() string {
	currentSideContext := gui.currentSideListContext()
	if currentSideContext == nil {
		return ""
	}

	return currentSideContext.GetSelectedItemId()
}

func (gui *Gui) isContextVisible(c types.Context) bool {
	return gui.State.WindowViewNameMap[c.GetWindowName()] == c.GetViewName() &&
		gui.State.ViewContextMap.Get(c.GetViewName()).GetKey() == c.GetKey()
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
