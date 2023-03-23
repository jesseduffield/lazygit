package controllers

import (
	"github.com/jesseduffield/gocui"
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
		{
			Key:         opts.GetKey(opts.Config.Universal.CreateRebaseOptionsMenu),
			Handler:     self.c.Helpers().MergeAndRebase.CreateRebaseOptionsMenu,
			Description: self.c.Tr.ViewMergeRebaseOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Refresh),
			Handler:     self.refresh,
			Description: self.c.Tr.LcRefresh,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextScreenMode),
			Handler:     self.nextScreenMode,
			Description: self.c.Tr.LcNextScreenMode,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevScreenMode),
			Handler:     self.prevScreenMode,
			Description: self.c.Tr.LcPrevScreenMode,
		},
		{
			ViewName:  "",
			Key:       opts.GetKey(opts.Config.Universal.OptionMenu),
			Handler:   self.createOptionsMenu,
			OpensMenu: true,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.OptionMenuAlt1),
			Modifier: gocui.ModNone,
			// we have the description on the alt key and not the main key for legacy reasons
			// (the original main key was 'x' but we've reassigned that to other purposes)
			Description: self.c.Tr.LcOpenMenu,
			Handler:     self.createOptionsMenu,
		},
	}
}

func (self *GlobalController) Context() types.Context {
	return nil
}

func (self *GlobalController) customCommand() error {
	return (&CustomCommandAction{c: self.c}).Call()
}

func (self *GlobalController) createCustomPatchOptionsMenu() error {
	return (&CustomPatchOptionsMenuAction{c: self.c}).Call()
}

func (self *GlobalController) refresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *GlobalController) nextScreenMode() error {
	return (&ScreenModeActions{c: self.c}).Next()
}

func (self *GlobalController) prevScreenMode() error {
	return (&ScreenModeActions{c: self.c}).Prev()
}

func (self *GlobalController) createOptionsMenu() error {
	return (&OptionsMenuAction{c: self.c}).Call()
}
