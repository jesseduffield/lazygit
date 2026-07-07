package controllers

import (
	"strconv"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
			Keys:        opts.GetKeys(opts.Config.Universal.ExecuteShellCommand),
			Handler:     self.shellCommand,
			Description: self.c.Tr.ExecuteShellCommand,
			Tooltip:     self.c.Tr.ExecuteShellCommandTooltip,
			OpensMenu:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.CreatePatchOptionsMenu),
			Handler:     self.createCustomPatchOptionsMenu,
			Description: self.c.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.CreateRebaseOptionsMenu),
			Handler:           opts.Guards.NoPopupPanel(self.c.Helpers().MergeAndRebase.CreateRebaseOptionsMenu),
			Description:       self.c.Tr.ViewMergeRebaseOptions,
			Tooltip:           self.c.Tr.ViewMergeRebaseOptionsTooltip,
			OpensMenu:         true,
			GetDisabledReason: self.canShowRebaseOptions,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.Refresh),
			Handler:     opts.Guards.NoPopupPanel(self.refresh),
			Description: self.c.Tr.Refresh,
			Tooltip:     self.c.Tr.RefreshTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.NextScreenMode),
			Handler:     opts.Guards.NoPopupPanel(self.nextScreenMode),
			Description: self.c.Tr.NextScreenMode,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.PrevScreenMode),
			Handler:     opts.Guards.NoPopupPanel(self.prevScreenMode),
			Description: self.c.Tr.PrevScreenMode,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.CyclePagers),
			Handler:           opts.Guards.NoPopupPanel(self.cyclePagers),
			GetDisabledReason: self.canCyclePagers,
			Description:       self.c.Tr.CyclePagers,
			Tooltip:           self.c.Tr.CyclePagersTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.CyclePagersReverse),
			Handler:           opts.Guards.NoPopupPanel(self.cyclePagersBackward),
			GetDisabledReason: self.canCyclePagers,
			Description:       self.c.Tr.CyclePagersReverse,
			Tooltip:           self.c.Tr.CyclePagersReverseTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Return),
			Handler:           self.escape,
			Description:       self.c.Tr.Cancel,
			DescriptionFunc:   self.escapeDescription,
			GetDisabledReason: self.escapeEnabled,
			DisplayOnScreen:   true,
		},
		{
			ViewName:          "",
			Keys:              opts.GetKeys(opts.Config.Universal.OptionMenu),
			Description:       self.c.Tr.OpenKeybindingsMenu,
			ShortDescription:  self.c.Tr.Keybindings,
			Handler:           self.createOptionsMenu,
			GetDisabledReason: self.optionsMenuDisabledReason,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.FilteringMenu),
			Handler:     opts.Guards.NoPopupPanel(self.createFilteringMenu),
			Description: self.c.Tr.OpenFilteringMenu,
			Tooltip:     self.c.Tr.OpenFilteringMenuTooltip,
			OpensMenu:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.DiffingMenu),
			Handler:     opts.Guards.NoPopupPanel(self.createDiffingMenu),
			Description: self.c.Tr.ViewDiffingOptions,
			Tooltip:     self.c.Tr.ViewDiffingOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.Quit),
			Description: self.c.Tr.Quit,
			Handler:     self.quit,
		},
		{
			Keys:    opts.GetKeys(opts.Config.Universal.QuitWithoutChangingDirectory),
			Handler: self.quitWithoutChangingDirectory,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.SuspendApp),
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
			Keys:        opts.GetKeys(opts.Config.Universal.ToggleWhitespaceInDiffView),
			Handler:     self.toggleWhitespace,
			Description: self.c.Tr.ToggleWhitespaceInDiffView,
			Tooltip:     self.c.Tr.ToggleWhitespaceInDiffViewTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.ToggleWordDiffInDiffView),
			Handler:     self.toggleWordDiff,
			Description: self.c.Tr.ToggleWordDiff,
			Tooltip:     self.c.Tr.ToggleWordDiffTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.EditConfig),
			Handler:     self.editConfig,
			Description: self.c.Tr.EditConfig,
			Tooltip:     self.c.Tr.EditFileTooltip,
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

func (self *GlobalController) cyclePagers() error {
	self.c.State().GetPagerConfig().CyclePagers()
	self.onPagerChanged()
	return nil
}

func (self *GlobalController) cyclePagersBackward() error {
	self.c.State().GetPagerConfig().CyclePagersBackward()
	self.onPagerChanged()
	return nil
}

// onPagerChanged re-renders the main view so the newly selected pager takes
// effect, and shows a toast naming it.
func (self *GlobalController) onPagerChanged() {
	currentSide := self.c.Context().CurrentSide()
	currentKey := self.c.Context().Current().GetKey()
	if currentSide.GetKey() == currentKey ||
		currentKey == context.NORMAL_MAIN_CONTEXT_KEY ||
		currentKey == context.NORMAL_SECONDARY_CONTEXT_KEY {
		currentSide.HandleRenderToMain()
	}

	pagerConfig := self.c.State().GetPagerConfig()
	current, total := pagerConfig.CurrentPagerIndex()
	name := pagerConfig.CurrentPagerName()
	if name == "" {
		if pagerConfig.CurrentPagerUsesGitConfigDiff() {
			name = self.c.Tr.ExternalDiffPagerName
		} else {
			name = self.c.Tr.DefaultPagerName
		}
	}
	self.c.Toast(utils.ResolvePlaceholderString(self.c.Tr.SelectedPager, map[string]string{
		"name":    name,
		"current": strconv.Itoa(current + 1),
		"total":   strconv.Itoa(total),
	}))
}

func (self *GlobalController) canCyclePagers() *types.DisabledReason {
	_, total := self.c.State().GetPagerConfig().CurrentPagerIndex()
	if total <= 1 {
		return &types.DisabledReason{
			Text: self.c.Tr.CyclePagersDisabledReason,
		}
	}
	return nil
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

func (self *GlobalController) toggleWordDiff() error {
	return (&ToggleWordDiffAction{c: self.c}).Call()
}

func (self *GlobalController) editConfig() error {
	return (&EditConfigAction{c: self.c}).Call()
}

func (self *GlobalController) canShowRebaseOptions() *types.DisabledReason {
	if self.c.Model().WorkingTreeStateAtLastCommitRefresh.None() {
		return &types.DisabledReason{
			Text: self.c.Tr.NotMergingOrRebasing,
		}
	}
	return nil
}
