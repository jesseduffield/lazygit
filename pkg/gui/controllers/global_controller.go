package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type GlobalController struct {
	baseController
	c *ControllerCommon
}

func NewGlobalController(
	common *ControllerCommon,
) *GlobalController {
	return &GlobalController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *GlobalController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.ExecuteCustomCommand),
			Handler:     self.customCommand,
			Description: self.c.Tr.LcExecuteCustomCommand,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.CreatePatchOptionsMenu),
			Handler:     self.createCustomPatchOptionsMenu,
			Description: self.c.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
	}
}

func (self *GlobalController) customCommand() error {
	return (&CustomCommandAction{c: self.c}).Call()
}

func (self *GlobalController) createCustomPatchOptionsMenu() error {
	return (&CustomPatchOptionsMenuAction{c: self.c}).Call()
}

func (self *GlobalController) Context() types.Context {
	return nil
}
