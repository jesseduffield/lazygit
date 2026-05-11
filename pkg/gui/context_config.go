package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) contextTree() *context.ContextTree {
	contextCommon := &context.ContextCommon{
		IGuiCommon: gui.c.IGuiCommon,
		Common:     gui.c.Common,
	}
	return context.NewContextTree(contextCommon)
}

func (gui *Gui) defaultSideContext() types.Context {
	if gui.State.Modes.Filtering.Active() {
		return gui.State.Contexts.LocalCommits
	}

	return gui.State.Contexts.Files
}
