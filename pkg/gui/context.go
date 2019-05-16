package gui

func (gui *Gui) titleMap() map[string]string {
	return map[string]string{
		"commits":  gui.Tr.SLocalize("DiffTitle"),
		"branches": gui.Tr.SLocalize("LogTitle"),
		"files":    gui.Tr.SLocalize("DiffTitle"),
		"status":   "",
		"stash":    gui.Tr.SLocalize("DiffTitle"),
	}
}

func (gui *Gui) contextTitleMap() map[string]map[string]string {
	return map[string]map[string]string{
		"main": {
			"staging": gui.Tr.SLocalize("StagingMainTitle"),
			"merging": gui.Tr.SLocalize("MergingMainTitle"),
			"normal":  "",
		},
	}
}

func (gui *Gui) setMainTitle() error {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return nil
	}
	currentViewName := currentView.Name()
	var newTitle string
	if context, ok := gui.State.Contexts[currentViewName]; ok {
		newTitle = gui.contextTitleMap()[currentViewName][context]
	} else if title, ok := gui.titleMap()[currentViewName]; ok {
		newTitle = title
	} else {
		return nil
	}
	gui.getMainView().Title = newTitle
	return nil
}

func (gui *Gui) changeContext(viewName, context string) error {
	if gui.State.Contexts[viewName] == context {
		return nil
	}

	contextMap := gui.GetContextMap()

	gui.g.DeleteKeybindings(viewName)

	bindings := contextMap[viewName][context]
	for _, binding := range bindings {
		if err := gui.g.SetKeybinding(viewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	gui.State.Contexts[viewName] = context
	return gui.setMainTitle()
}

func (gui *Gui) setInitialContexts() error {
	contextMap := gui.GetContextMap()

	initialContexts := map[string]string{
		"main": "normal",
	}

	for viewName, context := range initialContexts {
		bindings := contextMap[viewName][context]
		for _, binding := range bindings {
			if err := gui.g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
				return err
			}
		}
	}

	gui.State.Contexts = initialContexts

	return nil
}
