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
			Keys:              opts.GetKeys(opts.Config.Universal.CycleDiffRenderers),
			Handler:           opts.Guards.NoPopupPanel(self.cycleDiffRenderers),
			GetDisabledReason: self.canCycleDiffRenderers,
			Description:       self.c.Tr.CycleDiffRenderers,
			Tooltip:           self.c.Tr.CycleDiffRenderersTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.CycleDiffRenderersReverse),
			Handler:           opts.Guards.NoPopupPanel(self.cycleDiffRenderersBackward),
			GetDisabledReason: self.canCycleDiffRenderers,
			Description:       self.c.Tr.CycleDiffRenderersReverse,
			Tooltip:           self.c.Tr.CycleDiffRenderersReverseTooltip,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.JumpToFile),
			Handler:         opts.Guards.NoPopupPanel(self.jumpToFileInDiff),
			Description:     self.c.Tr.JumpToFileInDiff,
			DescriptionFunc: self.jumpToFileInDiffDescription,
			Tooltip:         self.c.Tr.JumpToFileInDiffTooltip,
			OpensMenu:       true,
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
	self.c.Refresh(types.RefreshOptions{})
	return nil
}

func (self *GlobalController) nextScreenMode() error {
	return (&ScreenModeActions{c: self.c}).Next()
}

func (self *GlobalController) prevScreenMode() error {
	return (&ScreenModeActions{c: self.c}).Prev()
}

func (self *GlobalController) cycleDiffRenderers() error {
	self.c.State().GetDiffRendererConfigManager().CycleDiffRenderers()
	self.onDiffRenderersChanged()
	return nil
}

func (self *GlobalController) cycleDiffRenderersBackward() error {
	self.c.State().GetDiffRendererConfigManager().CycleDiffRenderersBackward()
	self.onDiffRenderersChanged()
	return nil
}

// onDiffRenderersChanged re-renders the main view so the newly selected diff renderer
// takes effect, and shows a toast naming it.
func (self *GlobalController) onDiffRenderersChanged() {
	currentSide := self.c.Context().CurrentSide()
	currentKey := self.c.Context().Current().GetKey()
	if currentSide.GetKey() == currentKey ||
		currentKey == context.NORMAL_MAIN_CONTEXT_KEY ||
		currentKey == context.NORMAL_SECONDARY_CONTEXT_KEY {
		// The new renderer lays the same diff out its own way, so the line you were
		// looking at ends up elsewhere in the view; keep it in front of you.
		self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().Normal.GetView())
		self.c.Helpers().DiffLine.PreserveDiffPositionOnRerender(self.c.Contexts().NormalSecondary.GetView())
		currentSide.HandleRenderToMain()
	}

	diffRendererConfigManager := self.c.State().GetDiffRendererConfigManager()
	current, total := diffRendererConfigManager.CurrentDiffRendererIndex()
	name := diffRendererConfigManager.CurrentDiffRendererName(self.c.Tr)
	self.c.Toast(utils.ResolvePlaceholderString(self.c.Tr.SelectedDiffRenderers, map[string]string{
		"name":    name,
		"current": strconv.Itoa(current + 1),
		"total":   strconv.Itoa(total),
	}))
}

func (self *GlobalController) canCycleDiffRenderers() *types.DisabledReason {
	_, total := self.c.State().GetDiffRendererConfigManager().CurrentDiffRendererIndex()
	if total <= 1 {
		return &types.DisabledReason{
			Text: self.c.Tr.CycleDiffRenderersDisabledReason,
		}
	}
	return nil
}

// jumpToFileInDiff offers the files of the diff the main section is showing in a menu,
// and scrolls that pane to the file picked. The panel the user is in keeps the focus;
// they are reading the diff from there, and the next commit or file to read is picked
// there too.
func (self *GlobalController) jumpToFileInDiff() error {
	pane := self.diffPane()
	if pane == nil {
		return nil
	}

	return self.c.Helpers().DiffLine.OpenJumpToFileMenu(pane, self.c.Tr.JumpToFileInDiff)
}

// jumpToFileInDiffDescription qualifies the command's description so that it is listed
// only where it applies. A command with no description is left out of the keybindings
// menu.
//
// It doesn't apply where the main section is showing content that is no diff of the
// panel's — a branch's commit log, the status dashboard, a message. Nor does it while
// the focus is in one of the panes, which bind the key themselves; the menu would
// otherwise offer it twice there, once for the pane and once among the global keys.
//
// The static Description stays as it is: the cheatsheets are generated from that, and
// they document what a key does rather than when it applies.
func (self *GlobalController) jumpToFileInDiffDescription() string {
	_, focusIsInAPane := self.c.Context().Current().(*context.MainContext)
	if focusIsInAPane || self.diffPane() == nil {
		return ""
	}
	return self.c.Tr.JumpToFileInDiff
}

// diffPane returns the pane of the main section showing the diff of the panel the user
// is in, and nil when neither of them is showing one. A pane is cleared as it is
// emptied, so a pane that says it is showing a diff is showing one. Its window also has
// to be showing the pane. Resolving a conflicted file puts the merge conflicts view
// there instead, and the pane behind it goes on holding the diff it last rendered.
//
// Where both panes show a diff — the unstaged and staged sides of a file — the answer
// is the upper one, the pane the keys for scrolling the section act on.
func (self *GlobalController) diffPane() *context.MainContext {
	for _, pane := range []*context.MainContext{
		self.c.Contexts().Normal, self.c.Contexts().NormalSecondary,
	} {
		onScreen := self.c.Helpers().Window.GetContextForWindow(pane.GetWindowName()) == pane
		if onScreen && pane.ContentIsDiff() {
			return pane
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
