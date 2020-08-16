package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/stack"
)

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
	}

	gui.State.MainContext = context
}

type contextManager struct {
	gui   *Gui
	stack stack.Stack
}

func (c *contextManager) push(contextKey string) {
	c.stack.Push(contextKey)
}

// push focus, pop focus.

type Context interface {
	OnFocus() error
}

type SimpleContext struct {
	Self Context
}

type RemotesContext struct {
	Self     Context
	Branches Context
}

type CommitsContext struct {
	Self  Context
	Files Context
}

type ContextTree struct {
	Status        SimpleContext
	Files         SimpleContext
	Branches      SimpleContext
	Remotes       RemotesContext
	Tags          SimpleContext
	Commits       CommitsContext
	Stash         SimpleContext
	Staging       SimpleContext
	PatchBuilding SimpleContext
	Merging       SimpleContext
	Menu          SimpleContext
	Credentials   SimpleContext
	Confirmation  SimpleContext
	CommitMessage SimpleContext
}

func (gui *Gui) createContextTree() {
	gui.State.Contexts = ContextTree{
		Files: SimpleContext{
			Self: gui.filesListView(),
		},
	}
}

// func (c *contextManager) pop() (string, bool) {
// 	value, ok := c.stack.Pop()

// 	if !ok {
// 		// bottom of the stack, let's go to the default context: the files context
// 		c.gui.switchFocus(nil, newView)
// 	}
// }
