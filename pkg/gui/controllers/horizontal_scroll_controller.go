package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type HorizontalScrollControllerFactory struct {
	c *ControllerCommon
}

func NewHorizontalScrollControllerFactory(c *ControllerCommon) *HorizontalScrollControllerFactory {
	return &HorizontalScrollControllerFactory{
		c: c,
	}
}

func (self *HorizontalScrollControllerFactory) Create(context types.Context) types.IController {
	return &HorizontalScrollController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

// HorizontalScrollController scrolls a view sideways with the horizontal mouse wheel.
type HorizontalScrollController struct {
	baseController
	c *ControllerCommon

	context types.Context
}

func (self *HorizontalScrollController) Context() types.Context {
	return self.context
}

func (self *HorizontalScrollController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelLeft,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollLeft()
			},
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelRight,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollRight()
			},
		},
	}
}

func (self *HorizontalScrollController) HandleScrollLeft() error {
	self.context.GetViewTrait().ScrollLeftBy(self.scrollAmount())
	afterScrollingSideways(self.context)

	return nil
}

func (self *HorizontalScrollController) HandleScrollRight() error {
	self.context.GetViewTrait().ScrollRightBy(self.scrollAmount())
	afterScrollingSideways(self.context)

	return nil
}

// scrollAmount is the number of columns that one event of the wheel scrolls. It is
// twice the number of lines that the vertical wheel scrolls, because a column is
// about half as wide as a line is tall.
func (self *HorizontalScrollController) scrollAmount() int {
	return 2 * self.c.UserConfig().Gui.ScrollHeight
}
