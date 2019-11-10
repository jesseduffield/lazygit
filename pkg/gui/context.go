package gui

func (gui *Gui) changeContext(context string) error {
	oldContext := gui.State.Context

	if gui.State.Context == context {
		return nil
	}

	contextMap := gui.GetContextMap()

	oldBindings := contextMap[oldContext]
	for _, binding := range oldBindings {
		if err := gui.g.DeleteKeybinding(binding.ViewName, binding.Key, binding.Modifier); err != nil {
			return err
		}
	}

	bindings := contextMap[context]
	for _, binding := range bindings {
		if err := gui.g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	gui.State.Context = context
	return nil
}

func (gui *Gui) setInitialContext() error {
	contextMap := gui.GetContextMap()

	initialContext := "normal"

	bindings := contextMap[initialContext]
	for _, binding := range bindings {
		if err := gui.g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	gui.State.Context = initialContext

	return nil
}
