package custom_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// KeybindingCreator takes a custom command along with its handler and returns a corresponding keybinding
type KeybindingCreator struct {
	c *helpers.HelperCommon
}

func NewKeybindingCreator(c *helpers.HelperCommon) *KeybindingCreator {
	return &KeybindingCreator{
		c: c,
	}
}

func (self *KeybindingCreator) call(customCommand config.CustomCommand, handler func() error) (*types.Binding, error) {
	if customCommand.Context == "" {
		return nil, formatContextNotProvidedError(customCommand)
	}

	viewName, err := self.getViewNameAndContexts(customCommand)
	if err != nil {
		return nil, err
	}

	description := customCommand.Description
	if description == "" {
		description = customCommand.Command
	}

	return &types.Binding{
		ViewName:    viewName,
		Key:         keybindings.GetKey(customCommand.Key),
		Modifier:    gocui.ModNone,
		Handler:     handler,
		Description: description,
	}, nil
}

func (self *KeybindingCreator) getViewNameAndContexts(customCommand config.CustomCommand) (string, error) {
	if customCommand.Context == "global" {
		return "", nil
	}

	ctx, ok := self.contextForContextKey(types.ContextKey(customCommand.Context))
	if !ok {
		return "", formatUnknownContextError(customCommand)
	}

	viewName := ctx.GetViewName()
	return viewName, nil
}

func (self *KeybindingCreator) contextForContextKey(contextKey types.ContextKey) (types.Context, bool) {
	for _, context := range self.c.Contexts().Flatten() {
		if context.GetKey() == contextKey {
			return context, true
		}
	}

	return nil, false
}

func formatUnknownContextError(customCommand config.CustomCommand) error {
	allContextKeyStrings := lo.Map(context.AllContextKeys, func(key types.ContextKey, _ int) string {
		return string(key)
	})

	return fmt.Errorf("Error when setting custom command keybindings: unknown context: %s. Key: %s, Command: %s.\nPermitted contexts: %s", customCommand.Context, customCommand.Key, customCommand.Command, strings.Join(allContextKeyStrings, ", "))
}

func formatContextNotProvidedError(customCommand config.CustomCommand) error {
	return fmt.Errorf("Error parsing custom command keybindings: context not provided (use context: 'global' for the global context). Key: %s, Command: %s", customCommand.Key, customCommand.Command)
}
