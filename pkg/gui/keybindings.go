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
func (self *Gui) GetCheatsheetKeybindings() []*types.Binding {
	self.g = &gocui.Gui{}
	if err := self.createAllViews(); err != nil {
		panic(err)
	}
	// need to instantiate views
	self.helpers = helpers.NewStubHelpers()
	self.State = &GuiRepoState{}
	self.State.Contexts = self.contextTree()
	self.State.ContextMgr = NewContextMgr(self, self.State.Contexts)
	self.resetHelpersAndControllers()
	bindings, _ := self.GetInitialKeybindings()
	return bindings
}

func (self *Gui) keybindingOpts() types.KeybindingsOpts {
	config := self.c.UserConfig().Keybinding

	guards := types.KeybindingGuards{
		OutsideFilterMode: self.outsideFilterMode,
		NoPopupPanel:      self.noPopupPanel,
	}

	return types.KeybindingsOpts{
		GetKey: keybindings.GetKey,
		Config: config,
		Guards: guards,
	}
}

// renaming receiver to 'self' to aid refactoring. Will probably end up moving all Gui handlers to this pattern eventually.
func (self *Gui) GetInitialKeybindings() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	opts := self.c.KeybindingsOpts()

	bindings := []*types.Binding{
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.OpenRecentRepos),
			Handler:     self.helpers.Repos.CreateRecentReposMenu,
			Description: self.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollUpMain),
			Handler:     self.scrollUpMain,
			Alternative: "fn+up/shift+k",
			Description: self.c.Tr.ScrollUpMainWindow,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollDownMain),
			Handler:     self.scrollDownMain,
			Alternative: "fn+down/shift+j",
			Description: self.c.Tr.ScrollDownMainWindow,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollUpMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollDownMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollUpMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ScrollDownMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownMain,
		},
		{
			ViewName:          "files",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:          "localBranches",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "remoteBranches",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyBranchNameToClipboard,
		},
		{
			ViewName:          "commits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:    "commits",
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     self.helpers.CherryPick.Reset,
			Description: self.c.Tr.ResetCherryPick,
		},
		{
			ViewName:          "reflogCommits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName:          "subCommits",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemCommitHashToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyCommitHashToClipboard,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleInfoClick,
		},
		{
			ViewName:          "commitFiles",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopyPathToClipboard,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ExtrasMenu),
			Handler:     self.handleCreateExtrasMenuPanel,
			Description: self.c.Tr.OpenCommandLogMenu,
			Tooltip:     self.c.Tr.OpenCommandLogMenuTooltip,
			OpensMenu:   true,
		},
		{
			ViewName:    "main",
			Key:         gocui.MouseWheelDown,
			Handler:     self.scrollDownMain,
			Description: self.c.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Key:         gocui.MouseWheelUp,
			Handler:     self.scrollUpMain,
			Description: self.c.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownSecondary,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpSecondary,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gocui.MouseWheelUp,
			Handler:  self.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gocui.MouseWheelDown,
			Handler:  self.scrollDownConfirmationPanel,
		},
		{
			ViewName:          "submodules",
			Key:               opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:           self.handleCopySelectedSideContextItemToClipboard,
			GetDisabledReason: self.getCopySelectedSideContextItemToClipboardDisabledReason,
			Description:       self.c.Tr.CopySubmoduleNameToClipboard,
		},
		{
			ViewName: "extras",
			Key:      gocui.MouseWheelUp,
			Handler:  self.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Key:      gocui.MouseWheelDown,
			Handler:  self.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleFocusCommandLog,
		},
	}

	mouseKeybindings := []*gocui.ViewMouseBinding{}
	for _, c := range self.State.Contexts.Flatten() {
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
			Handler:     self.handleNextTab,
			Description: self.c.Tr.NextTab,
			Tag:         "navigation",
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.PrevTab),
			Handler:     self.handlePrevTab,
			Description: self.c.Tr.PrevTab,
			Tag:         "navigation",
		},
	}...)

	return bindings, mouseKeybindings
}

func (self *Gui) GetInitialKeybindingsWithCustomCommands() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	// if the search or filter prompt is open, we only want the keybindings for
	// that context. It shouldn't be possible, for example, to open a menu while
	// the prompt is showing; you first need to confirm or cancel the search/filter.
	if currentContext := self.State.ContextMgr.Current(); currentContext.GetKey() == context.SEARCH_CONTEXT_KEY {
		bindings := currentContext.GetKeybindings(self.c.KeybindingsOpts())
		viewName := currentContext.GetViewName()
		for _, binding := range bindings {
			binding.ViewName = viewName
		}
		return bindings, nil
	}

	bindings, mouseBindings := self.GetInitialKeybindings()
	customBindings, err := self.CustomCommandsClient.GetCustomCommandKeybindings()
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
		// we ignore click events on views that aren't popup panels, when a popup panel is focused.
		// Unless both the current view and the clicked-on view are either commit message or commit
		// description, because we want to allow switching between those two views by clicking.
		isCommitMessageView := func(viewName string) bool {
			return viewName == "commitMessage" || viewName == "commitDescription"
		}
		if gui.helpers.Confirmation.IsPopupPanelFocused() && gui.currentViewName() != binding.ViewName &&
			(!isCommitMessageView(gui.currentViewName()) || !isCommitMessageView(binding.ViewName)) {
			return nil
		}

		return baseHandler(opts)
	}
	binding.Handler = newHandler

	return gui.g.SetViewClickBinding(binding)
}

func (gui *Gui) callKeybindingHandler(binding *types.Binding) error {
	var disabledReason *types.DisabledReason
	if binding.GetDisabledReason != nil {
		disabledReason = binding.GetDisabledReason()
	}
	if disabledReason != nil {
		if disabledReason.ShowErrorInPanel {
			return errors.New(disabledReason.Text)
		}

		if len(disabledReason.Text) > 0 {
			gui.c.ErrorToast(gui.Tr.DisabledMenuItemPrefix + disabledReason.Text)
		}
		return nil
	}
	return binding.Handler()
}
