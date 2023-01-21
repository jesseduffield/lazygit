package gui

import (
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
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

	return gui.activateContext(c, types.OnFocusOpts{})
}

func (gui *Gui) pushContext(c types.Context, opts types.OnFocusOpts) error {
	if !c.IsFocusable() {
		return nil
	}

	contextsToDeactivate, contextToActivate := gui.pushToContextStack(c)

	for _, contextToDeactivate := range contextsToDeactivate {
		if err := gui.deactivateContext(contextToDeactivate, types.OnFocusLostOpts{NewContextKey: c.GetKey()}); err != nil {
			return err
		}
	}

	if contextToActivate == nil {
		return nil
	}

	return gui.activateContext(contextToActivate, opts)
}

// Adjusts the context stack based on the context that's being pushed and
// returns (contexts to deactivate, context to activate)
func (gui *Gui) pushToContextStack(c types.Context) ([]types.Context, types.Context) {
	contextsToDeactivate := []types.Context{}

	gui.State.ContextManager.Lock()
	defer gui.State.ContextManager.Unlock()

	if len(gui.State.ContextManager.ContextStack) > 0 &&
		c == gui.State.ContextManager.ContextStack[len(gui.State.ContextManager.ContextStack)-1] {
		// Context being pushed is already on top of the stack: nothing to
		// deactivate or activate
		return contextsToDeactivate, nil
	}

	if len(gui.State.ContextManager.ContextStack) == 0 {
		gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack, c)
	} else if c.GetKind() == types.SIDE_CONTEXT {
		// if we are switching to a side context, remove all other contexts in the stack
		contextsToDeactivate = gui.State.ContextManager.ContextStack
		gui.State.ContextManager.ContextStack = []types.Context{c}
	} else if c.GetKind() == types.MAIN_CONTEXT {
		// if we're switching to a main context, remove all other main contexts in the stack
		contextsToKeep := []types.Context{}
		for _, stackContext := range gui.State.ContextManager.ContextStack {
			if stackContext.GetKind() == types.MAIN_CONTEXT {
				contextsToDeactivate = append(contextsToDeactivate, stackContext)
			} else {
				contextsToKeep = append(contextsToKeep, stackContext)
			}
		}
		gui.State.ContextManager.ContextStack = append(contextsToKeep, c)
	} else {
		topContext := gui.currentContextWithoutLock()

		// if we're pushing the same context on, we do nothing.
		if topContext.GetKey() != c.GetKey() {
			// if top one is a temporary popup, we remove it. Ideally you'd be able to
			// escape back to previous temporary popups, but because we're currently reusing
			// views for this, you might not be able to get back to where you previously were.
			// The exception is when going to the search context e.g. for searching a menu.
			if (topContext.GetKind() == types.TEMPORARY_POPUP && c.GetKey() != context.SEARCH_CONTEXT_KEY) ||
				// we only ever want one main context on the stack at a time.
				(topContext.GetKind() == types.MAIN_CONTEXT && c.GetKind() == types.MAIN_CONTEXT) {

				contextsToDeactivate = append(contextsToDeactivate, topContext)
				_, gui.State.ContextManager.ContextStack = slices.Pop(gui.State.ContextManager.ContextStack)
			}

			gui.State.ContextManager.ContextStack = append(gui.State.ContextManager.ContextStack, c)
		}
	}

	return contextsToDeactivate, c
}

func (gui *Gui) popContext() error {
	gui.State.ContextManager.Lock()

	if len(gui.State.ContextManager.ContextStack) == 1 {
		// cannot escape from bottommost context
		gui.State.ContextManager.Unlock()
		return nil
	}

	var currentContext types.Context
	currentContext, gui.State.ContextManager.ContextStack = slices.Pop(gui.State.ContextManager.ContextStack)

	newContext := gui.State.ContextManager.ContextStack[len(gui.State.ContextManager.ContextStack)-1]

	gui.State.ContextManager.Unlock()

	if err := gui.deactivateContext(currentContext, types.OnFocusLostOpts{NewContextKey: newContext.GetKey()}); err != nil {
		return err
	}

	return gui.activateContext(newContext, types.OnFocusOpts{})
}

func (gui *Gui) removeContexts(contextsToRemove []types.Context) error {
	gui.State.ContextManager.Lock()

	if len(gui.State.ContextManager.ContextStack) == 1 {
		gui.State.ContextManager.Unlock()
		return nil
	}

	rest := lo.Filter(gui.State.ContextManager.ContextStack, func(context types.Context, _ int) bool {
		for _, contextToRemove := range contextsToRemove {
			if context.GetKey() == contextToRemove.GetKey() {
				return false
			}
		}
		return true
	})
	gui.State.ContextManager.ContextStack = rest
	contextToActivate := rest[len(rest)-1]
	gui.State.ContextManager.Unlock()

	for _, context := range contextsToRemove {
		if err := gui.deactivateContext(context, types.OnFocusLostOpts{NewContextKey: contextToActivate.GetKey()}); err != nil {
			return err
		}
	}

	// activate the item at the top of the stack
	return gui.activateContext(contextToActivate, types.OnFocusOpts{})
}

func (gui *Gui) deactivateContext(c types.Context, opts types.OnFocusLostOpts) error {
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

	if err := c.HandleFocusLost(opts); err != nil {
		return err
	}

	return nil
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context) error {
	if err := c.HandleRender(); err != nil {
		return err
	}

	if gui.currentViewName() == c.GetViewName() {
		if err := c.HandleFocus(types.OnFocusOpts{}); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) activateContext(c types.Context, opts types.OnFocusOpts) error {
	viewName := c.GetViewName()
	v, err := gui.g.View(viewName)
	if err != nil {
		return err
	}

	gui.setWindowContext(c)

	gui.moveToTopOfWindow(c)
	if _, err := gui.g.SetCurrentView(viewName); err != nil {
		return err
	}

	desiredTitle := c.Title()
	if desiredTitle != "" {
		v.Title = desiredTitle
	}

	v.Visible = true

	gui.g.Cursor = v.Editable

	// render the options available for the current context at the bottom of the screen
	optionsMap := c.GetOptionsMap()
	if optionsMap == nil {
		optionsMap = gui.globalOptionsMap()
	}
	gui.renderOptionsMap(optionsMap)

	if err := c.HandleFocus(opts); err != nil {
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

	return gui.currentStaticContextWithoutLock()
}

func (gui *Gui) currentStaticContextWithoutLock() types.Context {
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
			if err := gui.onViewFocusLost(previousView); err != nil {
				return err
			}

			previousView = newView
		}
		return nil
	}
}

func (gui *Gui) onViewFocusLost(oldView *gocui.View) error {
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

func (gui *Gui) rerenderView(view *gocui.View) error {
	context, ok := gui.contextForView(view.Name())
	if !ok {
		gui.Log.Errorf("no context found for view %s", view.Name())
		return nil
	}

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
