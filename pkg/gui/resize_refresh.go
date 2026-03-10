package gui

import (
	gui_context "github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// refreshMainContentOnResize recalculates the main diff-like content when a
// resize changes formatting-sensitive output (e.g. external diff pagers).
func (gui *Gui) refreshMainContentOnResize() {
	current := gui.c.Context().CurrentStatic()

	switch current.GetKind() {
	case types.SIDE_CONTEXT:
		current.HandleRenderToMain()

	case types.MAIN_CONTEXT:
		switch current.GetKey() {
		case gui_context.NORMAL_MAIN_CONTEXT_KEY, gui_context.NORMAL_SECONDARY_CONTEXT_KEY:
			// If focus is in a normal main context, rerender via the side context
			// because that's where the diff-render task is defined.
			gui.c.Context().CurrentSide().HandleRenderToMain()

		case gui_context.STAGING_MAIN_CONTEXT_KEY,
			gui_context.STAGING_SECONDARY_CONTEXT_KEY,
			gui_context.PATCH_BUILDING_MAIN_CONTEXT_KEY,
			gui_context.MERGE_CONFLICTS_CONTEXT_KEY:
			// These contexts refresh their own content from GetOnFocus handlers.
			current.HandleFocus(types.OnFocusOpts{})
		}
	}
}
