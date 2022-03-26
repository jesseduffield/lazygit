package custom_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// KeybindingCreator takes a custom command along with its handler and returns a corresponding keybinding
type KeybindingCreator struct {
	contexts *context.ContextTree
	getKey   func(string) types.Key
}

func NewKeybindingCreator(contexts *context.ContextTree, getKey func(string) types.Key) *KeybindingCreator {
	return &KeybindingCreator{
		contexts: contexts,
		getKey:   getKey,
	}
}

func (self *KeybindingCreator) call(customCommand config.CustomCommand, handler func() error) (*types.Binding, error) {
	if customCommand.Context == "" {
		return nil, formatContextNotProvidedError(customCommand)
	}

	viewName, contexts, err := self.getViewNameAndContexts(customCommand)
	if err != nil {
		return nil, err
	}

	description := customCommand.Description
	if description == "" {
		description = customCommand.Command
	}

	return &types.Binding{
		ViewName:    viewName,
		Contexts:    contexts,
		Key:         self.getKey(customCommand.Key),
		Modifier:    gocui.ModNone,
		Handler:     handler,
		Description: description,
	}, nil
}

func (self *KeybindingCreator) getViewNameAndContexts(customCommand config.CustomCommand) (string, []string, error) {
	if customCommand.Context == "global" {
		return "", nil, nil
	}

	ctx, ok := self.contextForContextKey(types.ContextKey(customCommand.Context))
	if !ok {
		return "", nil, formatUnknownContextError(customCommand)
	}

	// here we assume that a given context will always belong to the same view.
	// Currently this is a safe bet but it's by no means guaranteed in the long term
	// and we might need to make some changes in the future to support it.
	viewName := ctx.GetViewName()
	contexts := []string{customCommand.Context}
	return viewName, contexts, nil
}

func (self *KeybindingCreator) contextForContextKey(contextKey types.ContextKey) (types.Context, bool) {
	for _, context := range self.contexts.Flatten() {
		if context.GetKey() == contextKey {
			return context, true
		}
	}

	return nil, false
}

func formatUnknownContextError(customCommand config.CustomCommand) error {
	allContextKeyStrings := slices.Map(context.AllContextKeys, func(key types.ContextKey) string {
		return string(key)
	})

	return fmt.Errorf("Error when setting custom command keybindings: unknown context: %s. Key: %s, Command: %s.\nPermitted contexts: %s", customCommand.Context, customCommand.Key, customCommand.Command, strings.Join(allContextKeyStrings, ", "))
}

func formatContextNotProvidedError(customCommand config.CustomCommand) error {
	return fmt.Errorf("Error parsing custom command keybindings: context not provided (use context: 'global' for the global context). Key: %s, Command: %s", customCommand.Key, customCommand.Command)
}
