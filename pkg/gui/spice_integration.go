package gui

// isSpiceEnabled checks if git-spice integration is both available and enabled
func (gui *Gui) isSpiceEnabled() bool {
	if !gui.c.UserConfig().Git.Spice.Enabled {
		return false
	}
	return gui.git.Spice != nil && gui.git.Spice.IsAvailable()
}
