package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FilterControllerFactory struct {
	c *ControllerCommon
}

func NewFilterControllerFactory(c *ControllerCommon) *FilterControllerFactory {
	return &FilterControllerFactory{
		c: c,
	}
}

func (self *FilterControllerFactory) Create(context types.IFilterableContext) *FilterController {
	return &FilterController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

type FilterController struct {
	baseController
	c *ControllerCommon

	context types.IFilterableContext
}

func (self *FilterController) Context() types.Context {
	return self.context
}

func (self *FilterController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     self.OpenFilterPrompt,
			Description: self.c.Tr.StartFilter,
		},
	}
}

func (self *FilterController) OpenFilterPrompt() error {
	return self.c.Helpers().Search.OpenFilterPrompt(self.context)
}
