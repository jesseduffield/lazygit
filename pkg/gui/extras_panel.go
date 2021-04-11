package gui

func (gui *Gui) handleCreateExtrasMenuPanel() error {
	menuItems := []*menuItem{
		{
			displayString: gui.Tr.ToggleShowCommandLog,
			onPress: func() error {
				gui.ShowExtrasWindow = !gui.ShowExtrasWindow
				return nil
			},
		},
		{
			displayString: gui.Tr.FocusCommandLog,
			onPress: func() error {
				gui.ShowExtrasWindow = true
				gui.State.Contexts.CommandLog.SetParentContext(gui.currentSideContext())
				return gui.pushContext(gui.State.Contexts.CommandLog)
			},
		},
	}

	return gui.createMenu(gui.Tr.DiffingMenuTitle, menuItems, createMenuOptions{showCancel: true})
}
