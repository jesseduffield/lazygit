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
	for _, customCommand := range self.c.UserConfig().CustomCommands {
		if len(customCommand.CommandMenu) > 0 {
			handler := func() error {
				return self.showCustomCommandsMenu(customCommand)
			}
			bindings = append(bindings, &types.Binding{
				ViewName:    "", // custom commands menus are global; we filter the commands inside by context
				Key:         keybindings.GetKey(customCommand.Key),
				Modifier:    gocui.ModNone,
				Handler:     handler,
				Description: getCustomCommandsMenuDescription(customCommand, self.c.Tr),
				OpensMenu:   true,
			})
		} else {
			handler := self.handlerCreator.call(customCommand)
			compoundBindings, err := self.keybindingCreator.call(customCommand, handler)
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, compoundBindings...)
		}
	}

	return bindings, nil
}

func (self *Client) showCustomCommandsMenu(customCommand config.CustomCommand) error {
	menuItems := make([]*types.MenuItem, 0, len(customCommand.CommandMenu))
	for _, subCommand := range customCommand.CommandMenu {
		if len(subCommand.CommandMenu) > 0 {
			handler := func() error {
				return self.showCustomCommandsMenu(subCommand)
			}
			menuItems = append(menuItems, &types.MenuItem{
				Label:     subCommand.GetDescription(),
				Key:       keybindings.GetKey(subCommand.Key),
				OnPress:   handler,
				OpensMenu: true,
			})
		} else {
			if subCommand.Context != "" && subCommand.Context != "global" {
				viewNames, err := self.keybindingCreator.getViewNamesAndContexts(subCommand)
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
				Label:   subCommand.GetDescription(),
				Key:     keybindings.GetKey(subCommand.Key),
				OnPress: self.handlerCreator.call(subCommand),
			})
		}
	}

	if len(menuItems) == 0 {
		menuItems = append(menuItems, &types.MenuItem{
			Label:   self.c.Tr.NoApplicableCommandsInThisContext,
			OnPress: func() error { return nil },
		})
	}

	title := getCustomCommandsMenuDescription(customCommand, self.c.Tr)
	return self.c.Menu(types.CreateMenuOptions{Title: title, Items: menuItems, HideCancel: true})
}

func getCustomCommandsMenuDescription(customCommand config.CustomCommand, tr *i18n.TranslationSet) string {
	if customCommand.Description != "" {
		return customCommand.Description
	}

	return tr.CustomCommands
}
