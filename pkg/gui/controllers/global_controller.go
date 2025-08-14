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
			Key:         opts.GetKey(opts.Config.Universal.ExecuteShellCommand),
			Handler:     self.shellCommand,
			Description: self.c.Tr.ExecuteShellCommand,
			Tooltip:     self.c.Tr.ExecuteShellCommandTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.CreatePatchOptionsMenu),
			Handler:     self.createCustomPatchOptionsMenu,
			Description: self.c.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.CreateRebaseOptionsMenu),
			Handler:           opts.Guards.NoPopupPanel(self.c.Helpers().MergeAndRebase.CreateRebaseOptionsMenu),
			Description:       self.c.Tr.ViewMergeRebaseOptions,
			Tooltip:           self.c.Tr.ViewMergeRebaseOptionsTooltip,
			OpensMenu:         true,
			GetDisabledReason: self.canShowRebaseOptions,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Refresh),
			Handler:     opts.Guards.NoPopupPanel(self.refresh),
			Description: self.c.Tr.Refresh,
			Tooltip:     self.c.Tr.RefreshTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.NextScreenMode),
			Handler:     opts.Guards.NoPopupPanel(self.nextScreenMode),
			Description: self.c.Tr.NextScreenMode,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.PrevScreenMode),
			Handler:     opts.Guards.NoPopupPanel(self.prevScreenMode),
			Description: self.c.Tr.PrevScreenMode,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Return),
			Modifier:          gocui.ModNone,
			Handler:           self.escape,
			Description:       self.c.Tr.Cancel,
			DescriptionFunc:   self.escapeDescription,
			GetDisabledReason: self.escapeEnabled,
			DisplayOnScreen:   true,
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
			Description:       self.c.Tr.OpenKeybindingsMenu,
			Handler:           self.createOptionsMenu,
			ShortDescription:  self.c.Tr.Keybindings,
			DisplayOnScreen:   true,
			GetDisabledReason: self.optionsMenuDisabledReason,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.FilteringMenu),
			Handler:     opts.Guards.NoPopupPanel(self.createFilteringMenu),
			Description: self.c.Tr.OpenFilteringMenu,
			Tooltip:     self.c.Tr.OpenFilteringMenuTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenu),
			Handler:     opts.Guards.NoPopupPanel(self.createDiffingMenu),
			Description: self.c.Tr.ViewDiffingOptions,
			Tooltip:     self.c.Tr.ViewDiffingOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenuAlt),
			Handler:     opts.Guards.NoPopupPanel(self.createDiffingMenu),
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
			Key:         opts.GetKey(opts.Config.Universal.SuspendApp),
			Modifier:    gocui.ModNone,
			Handler:     self.c.Helpers().SuspendResume.SuspendApp,
			Description: self.c.Tr.SuspendApp,
			GetDisabledReason: func() *types.DisabledReason {
				if !self.c.Helpers().SuspendResume.CanSuspendApp() {
					return &types.DisabledReason{
						Text: self.c.Tr.CannotSuspendApp,
					}
				}
				return nil
			},
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

func (self *GlobalController) shellCommand() error {
	return (&ShellCommandAction{c: self.c}).Call()
}

func (self *GlobalController) createCustomPatchOptionsMenu() error {
	return (&CustomPatchOptionsMenuAction{c: self.c}).Call()
}

func (self *GlobalController) refresh() error {
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
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

func (self *GlobalController) optionsMenuDisabledReason() *types.DisabledReason {
	ctx := self.c.Context().Current()
	// Don't show options menu while displaying popup.
	if ctx.GetKind() == types.PERSISTENT_POPUP || ctx.GetKind() == types.TEMPORARY_POPUP {
		// The empty error text is intentional. We don't want to show an error
		// toast for this, but only hide it from the options map.
		return &types.DisabledReason{Text: ""}
	}
	return nil
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

func (self *GlobalController) escapeDescription() string {
	return (&QuitActions{c: self.c}).EscapeDescription()
}

func (self *GlobalController) escapeEnabled() *types.DisabledReason {
	if (&QuitActions{c: self.c}).EscapeEnabled() {
		return nil
	}

	// The empty error text is intentional. We don't want to show an error
	// toast for this, but only hide it from the options map.
	return &types.DisabledReason{Text: ""}
}

func (self *GlobalController) toggleWhitespace() error {
	return (&ToggleWhitespaceAction{c: self.c}).Call()
}

func (self *GlobalController) canShowRebaseOptions() *types.DisabledReason {
	if self.c.Model().WorkingTreeStateAtLastCommitRefresh.None() {
		return &types.DisabledReason{
			Text: self.c.Tr.NotMergingOrRebasing,
		}
	}
	return nil
}
