package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuController struct {
	baseController
	*ListControllerTrait[*types.MenuItem]
	c *ControllerCommon
}

var _ types.IController = &MenuController{}

func NewMenuController(
	c *ControllerCommon,
) *MenuController {
	return &MenuController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*types.MenuItem](
			c,
			c.Contexts().Menu,
			c.Contexts().Menu.GetSelected,
			c.Contexts().Menu.GetSelectedItems,
		),
		c: c,
	}
}

// NOTE: if you add a new keybinding here, you'll also need to add it to
// `reservedKeys` in `pkg/gui/context/menu_context.go`
func (self *MenuController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.press),
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Confirm),
			Handler:           self.withItem(self.press),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Execute,
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.Return),
			Handler:         self.close,
			Description:     self.c.Tr.Close,
			DisplayOnScreen: true,
		},
	}

	return bindings
}

func (self *MenuController) GetOnClick() func() error {
	return self.withItemGraceful(self.press)
}

func (self *MenuController) GetOnFocus() func(types.OnFocusOpts) error {
	return func(types.OnFocusOpts) error {
		selectedMenuItem := self.context().GetSelected()
		if selectedMenuItem != nil {
			self.c.Views().Tooltip.SetContent(self.c.Helpers().Confirmation.TooltipForMenuItem(selectedMenuItem))
		}
		return nil
	}
}

func (self *MenuController) press(selectedItem *types.MenuItem) error {
	return self.context().OnMenuPress(selectedItem)
}

func (self *MenuController) close() error {
	if self.context().IsFiltering() {
		self.c.Helpers().Search.Cancel()
		return nil
	}

	return self.c.PopContext()
}

func (self *MenuController) context() *context.MenuContext {
	return self.c.Contexts().Menu
}
