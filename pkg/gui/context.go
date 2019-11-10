package gui

func (gui *Gui) changeContext(viewName, context string) error {
	if gui.State.Contexts[viewName] == context {
		return nil
	}

	contextMap := gui.GetContextMap()

	gui.g.DeleteKeybindings(viewName)

	bindings := contextMap[viewName][context]
	for _, binding := range bindings {
		if err := gui.g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	gui.State.Contexts[viewName] = context
	return nil
}

func (gui *Gui) setInitialContexts() error {
	contextMap := gui.GetContextMap()

	initialContexts := map[string]string{
		"main":      "normal",
		"secondary": "normal",
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
