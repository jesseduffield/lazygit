package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Client is the entry point to this package. It returns a list of keybindings based on the config's user-defined custom commands.
// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Command_Keybindings.md for more info.
type Client struct {
	c                 *common.Common
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

	return &Client{
		c:                 c.Common,
		keybindingCreator: keybindingCreator,
		handlerCreator:    handlerCreator,
	}
}

func (self *Client) GetCustomCommandKeybindings() ([]*types.Binding, error) {
	bindings := []*types.Binding{}
	for _, customCommand := range self.c.UserConfig().CustomCommands {
		handler := self.handlerCreator.call(customCommand)
		compoundBindings, err := self.keybindingCreator.call(customCommand, handler)
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, compoundBindings...)
	}

	return bindings, nil
}
