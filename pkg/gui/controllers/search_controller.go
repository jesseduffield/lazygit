package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SearchControllerFactory struct {
	c *ControllerCommon
}

func NewSearchControllerFactory(c *ControllerCommon) *SearchControllerFactory {
	return &SearchControllerFactory{
		c: c,
	}
}

func (self *SearchControllerFactory) Create(context types.ISearchableContext) *SearchController {
	return &SearchController{
		baseController: baseController{},
		c:              self.c,
		context:        context,
	}
}

type SearchController struct {
	baseController
	c *ControllerCommon

	context types.ISearchableContext
}

func (self *SearchController) Context() types.Context {
	return self.context
}

func (self *SearchController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     self.OpenSearchPrompt,
			Description: self.c.Tr.StartSearch,
		},
	}
}

func (self *SearchController) OpenSearchPrompt() error {
	return self.c.Helpers().Search.OpenSearchPrompt(self.context)
}
