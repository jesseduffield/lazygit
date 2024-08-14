package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Client is the entry point to this package. It returns a list of keybindings based on the config's user-defined custom commands.
// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Command_Keybindings.md for more info.
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
		compoundBindings, err := self.keybindingCreator.call(customCommand, handler)
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, compoundBindings...)
	}

	return bindings, nil
}
