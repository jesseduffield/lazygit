package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

// given we have no fields here, arguably we shouldn't even need this factory
// struct, but we're maintaining consistency with the other files.
type VerticalScrollControllerFactory struct {
	controllerCommon     *controllerCommon
	viewBufferManagerMap *map[string]*tasks.ViewBufferManager
}

func NewVerticalScrollControllerFactory(c *controllerCommon, viewBufferManagerMap *map[string]*tasks.ViewBufferManager) *VerticalScrollControllerFactory {
	return &VerticalScrollControllerFactory{
		controllerCommon:     c,
		viewBufferManagerMap: viewBufferManagerMap,
	}
}

func (self *VerticalScrollControllerFactory) Create(context types.Context) types.IController {
	return &VerticalScrollController{
		baseController:       baseController{},
		controllerCommon:     self.controllerCommon,
		context:              context,
		viewBufferManagerMap: self.viewBufferManagerMap,
	}
}

type VerticalScrollController struct {
	baseController
	*controllerCommon

	context              types.Context
	viewBufferManagerMap *map[string]*tasks.ViewBufferManager
}

func (self *VerticalScrollController) Context() types.Context {
	return self.context
}

func (self *VerticalScrollController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{}
}

func (self *VerticalScrollController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelUp,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollUp()
			},
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseWheelDown,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.HandleScrollDown()
			},
		},
	}
}

func (self *VerticalScrollController) HandleScrollUp() error {
	self.context.GetViewTrait().ScrollUp(self.c.UserConfig.Gui.ScrollHeight)

	return nil
}

func (self *VerticalScrollController) HandleScrollDown() error {
	scrollHeight := self.c.UserConfig.Gui.ScrollHeight
	self.context.GetViewTrait().ScrollDown(scrollHeight)

	if manager, ok := (*self.viewBufferManagerMap)[self.context.GetViewName()]; ok {
		manager.ReadLines(scrollHeight)
	}

	return nil
}
