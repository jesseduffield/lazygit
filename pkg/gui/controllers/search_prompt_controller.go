package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SearchPromptController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &SearchPromptController{}

func NewSearchPromptController(
	common *ControllerCommon,
) *SearchPromptController {
	return &SearchPromptController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *SearchPromptController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:      opts.GetKey(opts.Config.Universal.Confirm),
			Modifier: gocui.ModNone,
			Handler:  self.confirm,
		},
		{
			Key:      opts.GetKey(opts.Config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  self.cancel,
		},
	}
}

func (self *SearchPromptController) Context() types.Context {
	return self.context()
}

func (self *SearchPromptController) context() types.Context {
	return self.c.Contexts().Search
}

func (self *SearchPromptController) confirm() error {
	return self.c.Helpers().Search.Confirm()
}

func (self *SearchPromptController) cancel() error {
	return self.c.Helpers().Search.CancelPrompt()
}
