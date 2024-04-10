package custom_commands

import (
	"github.com/lobes/lazytask/pkg/config"
	"github.com/lobes/lazytask/pkg/gui/controllers/helpers"
	"github.com/lobes/lazytask/pkg/gui/types"
)

// Client is the entry point to this package. It returns a list of keybindings based on the config's user-defined custom commands.
// See https://github.com/lobes/lazytask/blob/master/docs/Custom_Command_Keybindings.md for more info.
type Client struct {
	customCommands    []config.CustomCommand
	handlerCreator    *HandlerCreator
	keybindingCreator *KeybindingCreator
}

func NewClient(
	c *helpers.HelperCommon,
	helpers *helpers.Helpers,
) *Client {
	sessionStateLoader := NewSessionStateLoader(c, helpers.Refs)
	handlerCreator := NewHandlerCreator(
		c,
		sessionStateLoader,
		helpers.Suggestions,
		helpers.MergeAndRebase,
	)
	keybindingCreator := NewKeybindingCreator(c)
	customCommands := c.UserConfig.CustomCommands

	return &Client{
		customCommands:    customCommands,
		keybindingCreator: keybindingCreator,
		handlerCreator:    handlerCreator,
	}
}

func (self *Client) GetCustomCommandKeybindings() ([]*types.Binding, error) {
	bindings := []*types.Binding{}
	for _, customCommand := range self.customCommands {
		handler := self.handlerCreator.call(customCommand)
		binding, err := self.keybindingCreator.call(customCommand, handler)
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, binding)
	}

	return bindings, nil
}
