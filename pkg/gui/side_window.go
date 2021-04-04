package gui

func (gui *Gui) nextSideWindow() error {
	windows := gui.getCyclableWindows()
	currentWindow := gui.currentWindow()
	var newWindow string
	if currentWindow == "" || currentWindow == windows[len(windows)-1] {
		newWindow = windows[0]
	} else {
		for i := range windows {
			if currentWindow == windows[i] {
				newWindow = windows[i+1]
				break
			}
			if i == len(windows)-1 {
				return nil
			}
		}
	}
	if err := gui.resetOrigin(gui.Views.Main); err != nil {
		return err
	}

	viewName := gui.getViewNameForWindow(newWindow)

	return gui.pushContextWithView(viewName)
}

func (gui *Gui) previousSideWindow() error {
	windows := gui.getCyclableWindows()
	currentWindow := gui.currentWindow()
	var newWindow string
	if currentWindow == "" || currentWindow == windows[0] {
		newWindow = windows[len(windows)-1]
	} else {
		for i := range windows {
			if currentWindow == windows[i] {
				newWindow = windows[i-1]
				break
			}
			if i == len(windows)-1 {
				return nil
			}
		}
	}
	if err := gui.resetOrigin(gui.Views.Main); err != nil {
		return err
	}

	viewName := gui.getViewNameForWindow(newWindow)

	return gui.pushContextWithView(viewName)
}

func (gui *Gui) goToSideWindow(sideViewName string) func() error {
	return func() error {
		return gui.pushContextWithView(sideViewName)
	}
}
