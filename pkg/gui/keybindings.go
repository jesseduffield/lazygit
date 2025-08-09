package gui

import (
	"errors"
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
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
	config := gui.c.UserConfig().Keybinding

	guards := types.KeybindingGuards{
		OutsideFilterMode: gui.outsideFilterMode,
		NoPopupPanel:      gui.noPopupPanel,
	}

	return types.KeybindingsOpts{
		GetKey: keybindings.GetKey,
		Config: config,
		Guards: guards,
	}
}

func (gui *Gui) GetInitialKeybindings() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	opts := gui.c.KeybindingsOpts()

	bindings := []*types.Binding{
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.OpenRecentRepos),
			Handler:     opts.Guards.NoPopupPanel(gui.helpers.Repos.CreateRecentReposMenu),
			Description: gui.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollUpMain),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up/shift+k",
			Description: gui.c.Tr.ScrollUpMainWindow,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollDownMain),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down/shift+j",
			Description: gui.c.Tr.ScrollDownMainWindow,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollUpMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollDownMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollUpMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollDownMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:          "files",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:          "localBranches",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "remoteBranches",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "tags",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyTagToClipboard,
		},
		{
			ViewName:          "commits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:    "commits",
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     gui.helpers.CherryPick.Reset,
			Description: gui.c.Tr.ResetCherryPick,
		},
		{
			ViewName:          "reflogCommits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:          "subCommits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:          "commitFiles",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ExtrasMenu),
			Handler:     opts.Guards.NoPopupPanel(gui.handleCreateExtrasMenuPanel),
			Description: gui.c.Tr.OpenCommandLogMenu,
			Tooltip:     gui.c.Tr.OpenCommandLogMenuTooltip,
			OpensMenu:   true,
		},
		{
			ViewName:    "main",
			Key:         gocui.MouseWheelDown,
			Handler:     gui.scrollDownMain,
			Description: gui.c.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Key:         gocui.MouseWheelUp,
			Handler:     gui.scrollUpMain,
			Description: gui.c.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownSecondary,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpSecondary,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gocui.MouseWheelUp,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gocui.MouseWheelDown,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.NextPage),
			Modifier: gocui.ModNone,
			Handler:  gui.pageDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.PrevPage),
			Modifier: gocui.ModNone,
			Handler:  gui.pageUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.GotoTop),
			Modifier: gocui.ModNone,
			Handler:  gui.goToConfirmationPanelTop,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.GotoTopAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.goToConfirmationPanelTop,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.GotoBottom),
			Modifier: gocui.ModNone,
			Handler:  gui.goToConfirmationPanelBottom,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.GotoBottomAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.goToConfirmationPanelBottom,
		},
		{
			ViewName:          "submodules",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           gui.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: gui.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       gui.c.Tr.CopySubmoduleNameToClipboard,
		},
		{
			ViewName: "extras",
			Key:      gocui.MouseWheelUp,
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Key:      gocui.MouseWheelDown,
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.NextPage),
			Modifier: gocui.ModNone,
			Handler:  gui.pageDownExtrasPanel,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.PrevPage),
			Modifier: gocui.ModNone,
			Handler:  gui.pageUpExtrasPanel,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.GotoTop),
			Modifier: gocui.ModNone,
			Handler:  gui.goToExtrasPanelTop,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.GotoTopAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.goToExtrasPanelTop,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.GotoBottom),
			Modifier: gocui.ModNone,
			Handler:  gui.goToExtrasPanelBottom,
		},
		{
			ViewName: "extras",
			Key:      opts.GetKey(opts.Config.Universal.GotoBottomAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.goToExtrasPanelBottom,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
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
			Key:         opts.GetKey(opts.Config.Universal.NextTab),
			Handler:     opts.Guards.NoPopupPanel(gui.handleNextTab),
			Description: gui.c.Tr.NextTab,
			Tag:         "navigation",
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.PrevTab),
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
		if err := gui.SetKeybinding(binding); err != nil {
			return err
		}
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

func (gui *Gui) wrappedHandler(f func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}

func (gui *Gui) SetKeybinding(binding *types.Binding) error {
	handler := func() error {
		return gui.callKeybindingHandler(binding)
	}

	// TODO: move all mouse-ey stuff into new mouse approach
	if gocui.IsMouseKey(binding.Key) {
		handler = func() error {
			// we ignore click events on views that aren't popup panels, when a popup panel is focused
			if gui.helpers.Confirmation.IsPopupPanelFocused() && gui.currentViewName() != binding.ViewName {
				return nil
			}

			return binding.Handler()
		}
	}

	return gui.g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, gui.wrappedHandler(handler))
}

// warning: mutates the binding
func (gui *Gui) SetMouseKeybinding(binding *gocui.ViewMouseBinding) error {
	baseHandler := binding.Handler
	newHandler := func(opts gocui.ViewMouseBindingOpts) error {
		if gui.helpers.Confirmation.IsPopupPanelFocused() && gui.currentViewName() != binding.ViewName &&
			!gocui.IsMouseScrollKey(opts.Key) {
			// we ignore click events on views that aren't popup panels, when a popup panel is focused.
			// Unless both the current view and the clicked-on view are either commit message or commit
			// description, or a confirmation and the suggestions view, because we want to allow switching
			// between those two views by clicking.
			isCommitMessageOrSuggestionsView := func(viewName string) bool {
				return viewName == "commitMessage" || viewName == "commitDescription" ||
					viewName == "confirmation" || viewName == "suggestions"
			}
			if !isCommitMessageOrSuggestionsView(gui.currentViewName()) || !isCommitMessageOrSuggestionsView(binding.ViewName) {
				return nil
			}
		}

		return baseHandler(opts)
	}
	binding.Handler = newHandler

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
