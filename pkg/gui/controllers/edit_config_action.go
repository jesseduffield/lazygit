package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type EditConfigAction struct {
	c *ControllerCommon
}

func (self *EditConfigAction) Call() error {
	confPaths := self.c.GetConfig().GetUserConfigPaths()
	switch len(confPaths) {
	case 0:
		return errors.New(self.c.Tr.NoConfigFileFoundErr)
	case 1:
		return self.c.Helpers().Files.EditFiles(confPaths)
	default:
		menuItems := lo.Map(confPaths, func(path string, _ int) *types.MenuItem {
			return &types.MenuItem{
				Label: path,
				OnPress: func() error {
					return self.c.Helpers().Files.EditFiles([]string{path})
				},
			}
		})

		return self.c.Menu(types.CreateMenuOptions{
			Title: self.c.Tr.SelectConfigFile,
			Items: menuItems,
		})
	}
}
