package gui

func (gui *Gui) nextSideWindow() error {
	windows := gui.getCyclableWindows()
	currentWindow := gui.helpers.Window.CurrentWindow()
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
	gui.c.ResetViewOrigin(gui.Views.Main)

	context := gui.helpers.Window.GetContextForWindow(newWindow)

	return gui.c.PushContext(context)
}

func (gui *Gui) previousSideWindow() error {
	windows := gui.getCyclableWindows()
	currentWindow := gui.helpers.Window.CurrentWindow()
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
	gui.c.ResetViewOrigin(gui.Views.Main)

	context := gui.helpers.Window.GetContextForWindow(newWindow)

	return gui.c.PushContext(context)
}

func (gui *Gui) goToSideWindow(window string) func() error {
	return func() error {
		context := gui.helpers.Window.GetContextForWindow(window)

		return gui.c.PushContext(context)
	}
}

func (gui *Gui) getCyclableWindows() []string {
	return []string{"status", "files", "branches", "commits", "stash"}
}
