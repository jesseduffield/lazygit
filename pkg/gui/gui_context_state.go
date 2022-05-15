package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// TODO: better name
type guiContextStateFetcher struct {
	gui        *Gui
	contextKey types.ContextKey
}

var _ context.GuiContextState = &guiContextStateFetcher{}

func newGuiContextStateFetcher(
	gui *Gui,
	contextKey types.ContextKey,
) *guiContextStateFetcher {
	return &guiContextStateFetcher{
		gui:        gui,
		contextKey: contextKey,
	}
}

func (self *guiContextStateFetcher) Modes() *modes.Modes {
	return self.gui.State.Modes
}

func (self *guiContextStateFetcher) Needle() string {
	if self.Modes().Searching.SearchingInContext(self.contextKey) {
		return self.Modes().Searching.GetSearchString()
	}

	return ""
}

func (self *guiContextStateFetcher) ScreenMode() types.WindowMaximisation {
	return self.gui.State.ScreenMode
}

func (self *guiContextStateFetcher) IsFocused() bool {
	return self.gui.currentContext().GetKey() == self.contextKey
}

func (self *guiContextStateFetcher) BisectInfo() *git_commands.BisectInfo {
	return self.gui.State.Model.BisectInfo
}
