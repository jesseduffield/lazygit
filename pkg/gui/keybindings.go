package gui

import (
	"errors"
	"log"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) noPopupPanel(f func() error) func() error {
	return func() error {
		if gui.helpers.Confirmation.IsPopupPanelFocused() {
			return nil
		}

		return f()
	}
}

func (gui *Gui) outsideFilterMode(f func() error) func() error {
	return func() error {
		if !gui.validateNotInFilterMode() {
			return nil
		}

		return f()
	}
}

func (gui *Gui) validateNotInFilterMode() bool {
	if gui.State.Modes.Filtering.Active() {
		gui.c.Confirm(types.ConfirmOpts{
			Title:         gui.c.Tr.MustExitFilterModeTitle,
			Prompt:        gui.c.Tr.MustExitFilterModePrompt,
			HandleConfirm: gui.helpers.Mode.ExitFilterMode,
		})

		return false
	}
	return true
}

// only to be called from the cheatsheet generate script. This mutates the Gui struct.
func (gui *Gui) GetCheatsheetKeybindings() []*types.Binding {
	gui.g = &gocui.Gui{}
	if err := gui.createAllViews(); err != nil {
		panic(err)
	}
	// need to instantiate views
	gui.helpers = helpers.NewStubHelpers()
	gui.State = &GuiRepoState{}
	gui.State.Contexts = gui.contextTree()
	gui.State.ContextMgr = NewContextMgr(gui, gui.State.Contexts)
	gui.resetHelpersAndControllers()
	bindings, _ := gui.GetInitialKeybindings()
	return bindings
}

func (gui *Gui) keybindingOpts() types.KeybindingsOpts {
	keybindingConfig := gui.c.UserConfig().Keybinding

	guards := types.KeybindingGuards{
		OutsideFilterMode: gui.outsideFilterMode,
		NoPopupPanel:      gui.noPopupPanel,
	}

	return types.KeybindingsOpts{
		GetKeys: config.GetValidatedKeyBindingKeys,
		Config:  keybindingConfig,
		Guards:  guards,
	}
}

func (gui *Gui) GetInitialKeybindings() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	opts := gui.c.KeybindingsOpts()

	bindings := []*types.Binding{
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.OpenRecentRepos),
			Handler:     opts.Guards.NoPopupPanel(gui.helpers.Repos.CreateRecentReposMenu),
			Description: gui.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.ScrollUpMain),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up/shift+k",
			Description: gui.c.Tr.ScrollUpMainWindow,
		},
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.ScrollDownMain),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down/shift+j",
			Description: gui.c.Tr.ScrollDownMainWindow,
		},
		{
			ViewName:          "files",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:          "localBranches",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "remoteBranches",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "tags",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyTagToClipboard,
		},
		{
			ViewName:          "commits",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:    "commits",
			Keys:        opts.GetKeys(opts.Config.Commits.ResetCherryPick),
			Handler:     gui.helpers.CherryPick.Reset,
			Description: gui.c.Tr.ResetCherryPick,
		},
		{
			ViewName:          "reflogCommits",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:          "subCommits",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName: "information",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseLeft)},
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:          "commitFiles",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.ExtrasMenu),
			Handler:     opts.Guards.NoPopupPanel(gui.handleCreateExtrasMenuPanel),
			Description: gui.c.Tr.OpenCommandLogMenu,
			Tooltip:     gui.c.Tr.OpenCommandLogMenuTooltip,
			OpensMenu:   true,
		},
		{
			ViewName:    "main",
			Keys:        []gocui.Key{gocui.NewKeyName(gocui.MouseWheelDown)},
			Handler:     gui.scrollDownMain,
			Description: gui.c.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Keys:        []gocui.Key{gocui.NewKeyName(gocui.MouseWheelUp)},
			Handler:     gui.scrollUpMain,
			Description: gui.c.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "secondary",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelDown)},
			Handler:  gui.scrollDownSecondary,
		},
		{
			ViewName: "secondary",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelUp)},
			Handler:  gui.scrollUpSecondary,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.PrevItem),
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.NextItem),
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelUp)},
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelDown)},
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.NextPage),
			Handler:  gui.pageDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.PrevPage),
			Handler:  gui.pageUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.GotoTop),
			Handler:  gui.goToConfirmationPanelTop,
		},
		{
			ViewName: "confirmation",
			Keys:     opts.GetKeys(opts.Config.Universal.GotoBottom),
			Handler:  gui.goToConfirmationPanelBottom,
		},
		{
			ViewName:          "submodules",
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopySubmoduleNameToClipboard,
		},
		{
			ViewName: "extras",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelUp)},
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseWheelDown)},
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Keys:     opts.GetKeys(opts.Config.Universal.PrevItem),
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Keys:     opts.GetKeys(opts.Config.Universal.NextItem),
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Keys:     opts.GetKeys(opts.Config.Universal.NextPage),
			Handler:  gui.pageDownExtrasPanel,
		},
		{
			ViewName: "extras",
			Keys:     opts.GetKeys(opts.Config.Universal.PrevPage),
			Handler:  gui.pageUpExtrasPanel,
		},
		{
			ViewName: "extras",
			Keys:     opts.GetKeys(opts.Config.Universal.GotoTop),
			Handler:  gui.goToExtrasPanelTop,
		},
		{
			ViewName: "extras",
			Keys:     opts.GetKeys(opts.Config.Universal.GotoBottom),
			Handler:  gui.goToExtrasPanelBottom,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Keys:     []gocui.Key{gocui.NewKeyName(gocui.MouseLeft)},
			Handler:  gui.handleFocusCommandLog,
		},
	}

	mouseKeybindings := []*gocui.ViewMouseBinding{}
	for _, c := range gui.State.Contexts.Flatten() {
		viewName := c.GetViewName()
		for _, binding := range c.GetKeybindings(opts) {
			// TODO: move all mouse keybindings into the mouse keybindings approach below
			binding.ViewName = viewName
			bindings = append(bindings, binding)
		}

		mouseKeybindings = append(mouseKeybindings, c.GetMouseKeybindings(opts)...)
	}

	bindings = append(bindings, []*types.Binding{
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.NextTab),
			Handler:     opts.Guards.NoPopupPanel(gui.handleNextTab),
			Description: gui.c.Tr.NextTab,
			Tag:         "navigation",
		},
		{
			ViewName:    "",
			Keys:        opts.GetKeys(opts.Config.Universal.PrevTab),
			Handler:     opts.Guards.NoPopupPanel(gui.handlePrevTab),
			Description: gui.c.Tr.PrevTab,
			Tag:         "navigation",
		},
	}...)

	return bindings, mouseKeybindings
}

func (gui *Gui) GetInitialKeybindingsWithCustomCommands() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	// if the search or filter prompt is open, we only want the keybindings for
	// that context. It shouldn't be possible, for example, to open a menu while
	// the prompt is showing; you first need to confirm or cancel the search/filter.
	if currentContext := gui.State.ContextMgr.Current(); currentContext.GetKey() == context.SEARCH_CONTEXT_KEY {
		bindings := currentContext.GetKeybindings(gui.c.KeybindingsOpts())
		viewName := currentContext.GetViewName()
		for _, binding := range bindings {
			binding.ViewName = viewName
		}
		return bindings, nil
	}

	bindings, mouseBindings := gui.GetInitialKeybindings()
	customBindings, err := gui.CustomCommandsClient.GetCustomCommandKeybindings()
	if err != nil {
		log.Fatal(err)
	}
	// prepending because we want to give our custom keybindings precedence over default keybindings
	bindings = append(customBindings, bindings...)
	return bindings, mouseBindings
}

func (gui *Gui) resetKeybindings() error {
	gui.g.DeleteAllKeybindings()

	bindings, mouseBindings := gui.GetInitialKeybindingsWithCustomCommands()

	for _, binding := range bindings {
		gui.SetKeybinding(binding)
	}

	for _, binding := range mouseBindings {
		if err := gui.SetMouseKeybinding(binding); err != nil {
			return err
		}
	}

	for _, values := range gui.viewTabMap() {
		for _, value := range values {
			viewName := value.ViewName
			tabClickCallback := func(tabIndex int) error {
				return gui.onViewTabClick(gui.helpers.Window.WindowForView(viewName), tabIndex)
			}

			if err := gui.g.SetTabClickBinding(viewName, tabClickCallback); err != nil {
				return err
			}
		}
	}

	return nil
}

func (gui *Gui) SetKeybinding(binding *types.Binding) {
	handler := func(g *gocui.Gui, v *gocui.View) error {
		return gui.callKeybindingHandler(binding)
	}

	for _, key := range binding.Keys {
		gui.g.SetKeybinding(binding.ViewName, key, handler)
	}
}

func (gui *Gui) SetMouseKeybinding(binding *gocui.ViewMouseBinding) error {
	return gui.g.SetViewClickBinding(binding)
}

func (gui *Gui) callKeybindingHandler(binding *types.Binding) error {
	if binding.GetDisabledReason != nil {
		if disabledReason := binding.GetDisabledReason(); disabledReason != nil {
			if disabledReason.AllowFurtherDispatching {
				return &types.ErrKeybindingNotHandled{DisabledReason: disabledReason}
			}

			if disabledReason.ShowErrorInPanel {
				return errors.New(disabledReason.Text)
			}

			if len(disabledReason.Text) > 0 {
				gui.c.ErrorToast(gui.Tr.DisabledMenuItemPrefix + disabledReason.Text)
			}
			return nil
		}
	}

	return binding.Handler()
}
