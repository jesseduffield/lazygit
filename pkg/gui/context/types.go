package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/modes"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// GuiContextState is for obtaining info about the gui's state as well as
// state related to the given context
type GuiContextState interface {
	Modes() *modes.Modes

	// this is the search string when we're in filtering mode.
	Needle() string

	ScreenMode() types.WindowMaximisation

	IsFocused() bool

	BisectInfo() *git_commands.BisectInfo
}
