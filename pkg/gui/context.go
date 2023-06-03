package gui

import (
	"errors"
	"sync"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// This file is for the management of contexts. There is a context stack such that
// for example you might start off in the commits context and then open a menu, putting
// you in the menu context. When contexts are activated/deactivated certain things need
// to happen like showing/hiding views and rendering content.

type ContextMgr struct {
	ContextStack []types.Context
	sync.RWMutex
	gui *Gui

	allContexts *context.ContextTree
}

func NewContextMgr(
	gui *Gui,
	allContexts *context.ContextTree,
) *ContextMgr {
	return &ContextMgr{
		ContextStack: []types.Context{},
		RWMutex:      sync.RWMutex{},
		gui:          gui,
		allContexts:  allContexts,
	}
}

// use when you don't want to return to the original context upon
// hitting escape: you want to go that context's parent instead.
func (self *ContextMgr) Replace(c types.Context) error {
	if !c.IsFocusable() {
		return nil
	}

	self.Lock()

	if len(self.ContextStack) == 0 {
		self.ContextStack = []types.Context{c}
	} else {
		// replace the last item with the given item
		self.ContextStack = append(self.ContextStack[0:len(self.ContextStack)-1], c)
	}

	defer self.Unlock()

	return self.ActivateContext(c, types.OnFocusOpts{})
}

func (self *ContextMgr) Push(c types.Context, opts ...types.OnFocusOpts) error {
	if len(opts) > 1 {
		return errors.New("cannot pass multiple opts to Push")
	}

	singleOpts := types.OnFocusOpts{}
	if len(opts) > 0 {
		// using triple dot but you should only ever pass one of these opt structs
		singleOpts = opts[0]
	}

	if !c.IsFocusable() {
		return nil
	}

	contextsToDeactivate, contextToActivate := self.pushToContextStack(c)

	for _, contextToDeactivate := range contextsToDeactivate {
		if err := self.deactivateContext(contextToDeactivate, types.OnFocusLostOpts{NewContextKey: c.GetKey()}); err != nil {
			return err
		}
	}

	if contextToActivate == nil {
		return nil
	}

	return self.ActivateContext(contextToActivate, singleOpts)
}

// Adjusts the context stack based on the context that's being pushed and
// returns (contexts to deactivate, context to activate)
func (self *ContextMgr) pushToContextStack(c types.Context) ([]types.Context, types.Context) {
	contextsToDeactivate := []types.Context{}

	self.Lock()
	defer self.Unlock()

	if len(self.ContextStack) > 0 &&
		c.GetKey() == self.ContextStack[len(self.ContextStack)-1].GetKey() {
		// Context being pushed is already on top of the stack: nothing to
		// deactivate or activate
		return contextsToDeactivate, nil
	}

	if len(self.ContextStack) == 0 {
		self.ContextStack = append(self.ContextStack, c)
	} else if c.GetKind() == types.SIDE_CONTEXT {
		// if we are switching to a side context, remove all other contexts in the stack
		contextsToDeactivate = lo.Filter(self.ContextStack, func(context types.Context, _ int) bool {
			return context.GetKey() != c.GetKey()
		})
		self.ContextStack = []types.Context{c}
	} else if c.GetKind() == types.MAIN_CONTEXT {
		// if we're switching to a main context, remove all other main contexts in the stack
		contextsToKeep := []types.Context{}
		for _, stackContext := range self.ContextStack {
			if stackContext.GetKind() == types.MAIN_CONTEXT {
				contextsToDeactivate = append(contextsToDeactivate, stackContext)
			} else {
				contextsToKeep = append(contextsToKeep, stackContext)
			}
		}
		self.ContextStack = append(contextsToKeep, c)
	} else {
		topContext := self.currentContextWithoutLock()

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
				_, self.ContextStack = slices.Pop(self.ContextStack)
			}

			self.ContextStack = append(self.ContextStack, c)
		}
	}

	return contextsToDeactivate, c
}

func (self *ContextMgr) Pop() error {
	self.Lock()

	if len(self.ContextStack) == 1 {
		// cannot escape from bottommost context
		self.Unlock()
		return nil
	}

	var currentContext types.Context
	currentContext, self.ContextStack = slices.Pop(self.ContextStack)

	newContext := self.ContextStack[len(self.ContextStack)-1]

	self.Unlock()

	if err := self.deactivateContext(currentContext, types.OnFocusLostOpts{NewContextKey: newContext.GetKey()}); err != nil {
		return err
	}

	return self.ActivateContext(newContext, types.OnFocusOpts{})
}

func (self *ContextMgr) RemoveContexts(contextsToRemove []types.Context) error {
	self.Lock()

	if len(self.ContextStack) == 1 {
		self.Unlock()
		return nil
	}

	rest := lo.Filter(self.ContextStack, func(context types.Context, _ int) bool {
		for _, contextToRemove := range contextsToRemove {
			if context.GetKey() == contextToRemove.GetKey() {
				return false
			}
		}
		return true
	})
	self.ContextStack = rest
	contextToActivate := rest[len(rest)-1]
	self.Unlock()

	for _, context := range contextsToRemove {
		if err := self.deactivateContext(context, types.OnFocusLostOpts{NewContextKey: contextToActivate.GetKey()}); err != nil {
			return err
		}
	}

	// activate the item at the top of the stack
	return self.ActivateContext(contextToActivate, types.OnFocusOpts{})
}

func (self *ContextMgr) deactivateContext(c types.Context, opts types.OnFocusLostOpts) error {
	view, _ := self.gui.c.GocuiGui().View(c.GetViewName())

	if opts.NewContextKey != context.SEARCH_CONTEXT_KEY {
		if c.GetKind() == types.MAIN_CONTEXT || c.GetKind() == types.TEMPORARY_POPUP {
			self.gui.helpers.Search.CancelSearchIfSearching(c)
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

func (self *ContextMgr) ActivateContext(c types.Context, opts types.OnFocusOpts) error {
	viewName := c.GetViewName()
	v, err := self.gui.c.GocuiGui().View(viewName)
	if err != nil {
		return err
	}

	self.gui.helpers.Window.SetWindowContext(c)

	self.gui.helpers.Window.MoveToTopOfWindow(c)
	if _, err := self.gui.c.GocuiGui().SetCurrentView(viewName); err != nil {
		return err
	}

	self.gui.helpers.Search.RenderSearchStatus(c)

	desiredTitle := c.Title()
	if desiredTitle != "" {
		v.Title = desiredTitle
	}

	v.Visible = true

	self.gui.c.GocuiGui().Cursor = v.Editable

	self.gui.renderContextOptionsMap(c)

	if err := c.HandleFocus(opts); err != nil {
		return err
	}

	return nil
}

func (self *ContextMgr) Current() types.Context {
	self.RLock()
	defer self.RUnlock()

	return self.currentContextWithoutLock()
}

func (self *ContextMgr) currentContextWithoutLock() types.Context {
	if len(self.ContextStack) == 0 {
		return self.gui.defaultSideContext()
	}

	return self.ContextStack[len(self.ContextStack)-1]
}

// Note that this could return the 'status' context which is not itself a list context.
func (self *ContextMgr) CurrentSide() types.Context {
	self.RLock()
	defer self.RUnlock()

	stack := self.ContextStack

	// find the first context in the stack with the type of types.SIDE_CONTEXT
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() == types.SIDE_CONTEXT {
			return context
		}
	}

	return self.gui.defaultSideContext()
}

// static as opposed to popup
func (self *ContextMgr) CurrentStatic() types.Context {
	self.RLock()
	defer self.RUnlock()

	return self.currentStaticContextWithoutLock()
}

func (self *ContextMgr) currentStaticContextWithoutLock() types.Context {
	stack := self.ContextStack

	if len(stack) == 0 {
		return self.gui.defaultSideContext()
	}

	// find the first context in the stack without a popup type
	for i := range stack {
		context := stack[len(stack)-1-i]

		if context.GetKind() != types.TEMPORARY_POPUP && context.GetKind() != types.PERSISTENT_POPUP {
			return context
		}
	}

	return self.gui.defaultSideContext()
}

func (self *ContextMgr) ForEach(f func(types.Context)) {
	self.RLock()
	defer self.RUnlock()

	for _, context := range self.gui.State.ContextMgr.ContextStack {
		f(context)
	}
}

func (self *ContextMgr) IsCurrent(c types.Context) bool {
	return self.Current().GetKey() == c.GetKey()
}

func (self *ContextMgr) AllFilterable() []types.IFilterableContext {
	var result []types.IFilterableContext

	for _, context := range self.allContexts.Flatten() {
		if ctx, ok := context.(types.IFilterableContext); ok {
			result = append(result, ctx)
		}
	}

	return result
}

func (self *ContextMgr) AllSearchable() []types.ISearchableContext {
	var result []types.ISearchableContext

	for _, context := range self.allContexts.Flatten() {
		if ctx, ok := context.(types.ISearchableContext); ok {
			result = append(result, ctx)
		}
	}

	return result
}

// all list contexts
func (self *ContextMgr) AllList() []types.IListContext {
	var listContexts []types.IListContext

	for _, context := range self.allContexts.Flatten() {
		if listContext, ok := context.(types.IListContext); ok {
			listContexts = append(listContexts, listContext)
		}
	}

	return listContexts
}

func (self *ContextMgr) AllPatchExplorer() []types.IPatchExplorerContext {
	var listContexts []types.IPatchExplorerContext

	for _, context := range self.allContexts.Flatten() {
		if listContext, ok := context.(types.IPatchExplorerContext); ok {
			listContexts = append(listContexts, listContext)
		}
	}

	return listContexts
}
