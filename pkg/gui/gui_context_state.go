package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// TODO: better name
type guiContextStateFetcher[T any] struct {
	gui        *Gui
	contextKey types.ContextKey
	getItems   func() []T
}

var _ context.GuiContextState[any] = &guiContextStateFetcher[any]{}

func newGuiContextStateFetcher[T any](
	gui *Gui,
	contextKey types.ContextKey,
	getItems func() []T,
) *guiContextStateFetcher[T] {
	return &guiContextStateFetcher[T]{
		gui:        gui,
		contextKey: contextKey,
		getItems:   getItems,
	}
}

func (self *guiContextStateFetcher[T]) Modes() *modes.Modes {
	return self.gui.State.Modes
}

func (self *guiContextStateFetcher[T]) Items() []T {
	return self.getItems()
}

func (self *guiContextStateFetcher[T]) Needle() string {
	if self.Modes().Searching.SearchingInContext(self.contextKey) {
		return self.Modes().Searching.GetSearchString()
	}

	return ""
}

func (self *guiContextStateFetcher[T]) ScreenMode() types.WindowMaximisation {
	return self.gui.State.ScreenMode
}

func (self *guiContextStateFetcher[T]) IsFocused() bool {
	return self.gui.currentContext().GetKey() == self.contextKey
}
