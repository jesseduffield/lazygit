package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

// isSpiceEnabled checks if git-spice integration is both available and enabled
func (gui *Gui) isSpiceEnabled() bool {
	if !gui.c.UserConfig().Git.Spice.Enabled {
		return false
	}
	return gui.git.Spice != nil && gui.git.Spice.IsAvailable()
}

// spiceStacksContextOrNil returns the SpiceStacks context if it exists, otherwise nil.
// This is useful for conditionally adding the context to lists.
func (gui *Gui) spiceStacksContextOrNil() types.Context {
	if gui.State.Contexts.SpiceStacks != nil {
		return gui.State.Contexts.SpiceStacks
	}
	return nil
}
