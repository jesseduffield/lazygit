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
	c *ControllerCommon,
) *GlobalController {
	return &GlobalController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *GlobalController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.ExecuteCustomCommand),
			Handler:     self.customCommand,
			Description: self.c.Tr.ExecuteCustomCommand,
			Tooltip:     self.c.Tr.ExecuteCustomCommandTooltip,
			OpensMenu:   true,
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
			Tooltip:     self.c.Tr.ViewMergeRebaseOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Refresh),
			Handler:     self.refresh,
			Description: self.c.Tr.Refresh,
			Tooltip:     self.c.Tr.RefreshTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextScreenMode),
			Handler:     self.nextScreenMode,
			Description: self.c.Tr.NextScreenMode,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevScreenMode),
			Handler:     self.prevScreenMode,
			Description: self.c.Tr.PrevScreenMode,
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
			Description:      self.c.Tr.OpenKeybindingsMenu,
			Handler:          self.createOptionsMenu,
			ShortDescription: self.c.Tr.Keybindings,
			DisplayOnScreen:  true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.FilteringMenu),
			Handler:     self.createFilteringMenu,
			Description: self.c.Tr.OpenFilteringMenu,
			Tooltip:     self.c.Tr.OpenFilteringMenuTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenu),
			Handler:     self.createDiffingMenu,
			Description: self.c.Tr.ViewDiffingOptions,
			Tooltip:     self.c.Tr.ViewDiffingOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenuAlt),
			Handler:     self.createDiffingMenu,
			Description: self.c.Tr.ViewDiffingOptions,
			Tooltip:     self.c.Tr.ViewDiffingOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Quit),
			Modifier:    gocui.ModNone,
			Description: self.c.Tr.Quit,
			Handler:     self.quit,
		},
		{
			Key:      opts.GetKey(opts.Config.Universal.QuitAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.quit,
		},
		{
			Key:      opts.GetKey(opts.Config.Universal.QuitWithoutChangingDirectory),
			Modifier: gocui.ModNone,
			Handler:  self.quitWithoutChangingDirectory,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.Return),
			Modifier:        gocui.ModNone,
			Handler:         self.escape,
			Description:     self.c.Tr.Cancel,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.ToggleWhitespaceInDiffView),
			Handler:     self.toggleWhitespace,
			Description: self.c.Tr.ToggleWhitespaceInDiffView,
			Tooltip:     self.c.Tr.ToggleWhitespaceInDiffViewTooltip,
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

func (self *GlobalController) createFilteringMenu() error {
	return (&FilteringMenuAction{c: self.c}).Call()
}

func (self *GlobalController) createDiffingMenu() error {
	return (&DiffingMenuAction{c: self.c}).Call()
}

func (self *GlobalController) quit() error {
	return (&QuitActions{c: self.c}).Quit()
}

func (self *GlobalController) quitWithoutChangingDirectory() error {
	return (&QuitActions{c: self.c}).QuitWithoutChangingDirectory()
}

func (self *GlobalController) escape() error {
	return (&QuitActions{c: self.c}).Escape()
}

func (self *GlobalController) toggleWhitespace() error {
	return (&ToggleWhitespaceAction{c: self.c}).Call()
}
