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

	context := gui.getContextForWindow(newWindow)

	return gui.c.PushContext(context)
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

	context := gui.getContextForWindow(newWindow)

	return gui.c.PushContext(context)
}

func (gui *Gui) goToSideWindow(window string) func() error {
	return func() error {
		context := gui.getContextForWindow(window)

		return gui.c.PushContext(context)
	}
}
