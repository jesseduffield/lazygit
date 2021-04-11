package gui

func (gui *Gui) handleCreateExtrasMenuPanel() error {
	menuItems := []*menuItem{
		{
			displayString: "Toggle show/hide command log",
			onPress: func() error {
				gui.ShowExtrasWindow = !gui.ShowExtrasWindow
				return nil
			},
		},
	}

	return gui.createMenu(gui.Tr.DiffingMenuTitle, menuItems, createMenuOptions{showCancel: true})
}
