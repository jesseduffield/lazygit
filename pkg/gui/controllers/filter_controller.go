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

// A context that filters as the user types has an input field of its own, so it
// has no use for the filter prompt.
type contextThatFiltersAsYouType interface {
	FilterAsYouType() bool
}

func (self *FilterController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	if context, ok := self.context.(contextThatFiltersAsYouType); ok && context.FilterAsYouType() {
		return nil
	}

	return []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Universal.StartSearch),
			Handler:     self.OpenFilterPrompt,
			Description: self.c.Tr.StartFilter,
		},
	}
}

func (self *FilterController) OpenFilterPrompt() error {
	self.c.Helpers().Search.OpenFilterPrompt(self.context)
	return nil
}
