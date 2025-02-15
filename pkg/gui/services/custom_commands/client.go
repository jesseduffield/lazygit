package custom_commands

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/samber/lo"
)

// Client is the entry point to this package. It returns a list of keybindings based on the config's user-defined custom commands.
// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Command_Keybindings.md for more info.
type Client struct {
	c                 *helpers.HelperCommon
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
		c:                 c,
		keybindingCreator: keybindingCreator,
		handlerCreator:    handlerCreator,
	}
}

func (self *Client) GetCustomCommandKeybindings() ([]*types.Binding, error) {
	bindings := []*types.Binding{}
	for _, customCommandsMenu := range self.c.UserConfig().CustomCommandsMenus {
		handler := func() error {
			return self.showCustomCommandsMenu(customCommandsMenu)
		}
		bindings = append(bindings, &types.Binding{
			ViewName:    "", // custom commands menus are global; we filter the commands inside by context
			Key:         keybindings.GetKey(customCommandsMenu.Key),
			Modifier:    gocui.ModNone,
			Handler:     handler,
			Description: getCustomCommandsMenuDescription(customCommandsMenu, self.c.Tr),
		})
	}

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

func (self *Client) showCustomCommandsMenu(customCommandsMenu config.CustomCommandsMenu) error {
	menuItems := make([]*types.MenuItem, 0, len(customCommandsMenu.Commands))
	for _, command := range customCommandsMenu.Commands {
		if command.Context != "" && command.Context != "global" {
			viewNames, err := self.keybindingCreator.getViewNamesAndContexts(command)
			if err != nil {
				return err
			}

			currentView := self.c.GocuiGui().CurrentView()
			enabled := currentView != nil && lo.Contains(viewNames, currentView.Name())
			if !enabled {
				continue
			}
		}

		menuItems = append(menuItems, &types.MenuItem{
			Label:   command.GetDescription(),
			Key:     keybindings.GetKey(command.Key),
			OnPress: self.handlerCreator.call(command),
		})
	}

	if len(menuItems) == 0 {
		menuItems = append(menuItems, &types.MenuItem{
			Label:   self.c.Tr.NoApplicableCommandsInThisContext,
			OnPress: func() error { return nil },
		})
	}

	title := getCustomCommandsMenuDescription(customCommandsMenu, self.c.Tr)
	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems, HideCancel: true})
}

func getCustomCommandsMenuDescription(customCommandsMenu config.CustomCommandsMenu, tr *i18n.TranslationSet) string {
	if customCommandsMenu.Description != "" {
		return customCommandsMenu.Description
	}

	return tr.CustomCommands
}
