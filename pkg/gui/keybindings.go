package gui

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) noPopupPanel(f func() error) func() error {
	return func() error {
		if gui.popupPanelFocused() {
			return nil
		}

		return f()
	}
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
	self.resetControllers()
	bindings, _ := self.GetInitialKeybindings()
	return bindings
}

// renaming receiver to 'self' to aid refactoring. Will probably end up moving all Gui handlers to this pattern eventually.
func (self *Gui) GetInitialKeybindings() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	config := self.c.UserConfig.Keybinding

	guards := types.KeybindingGuards{
		OutsideFilterMode: self.outsideFilterMode,
		NoPopupPanel:      self.noPopupPanel,
	}

	opts := types.KeybindingsOpts{
		GetKey: keybindings.GetKey,
		Config: config,
		Guards: guards,
	}

	bindings := []*types.Binding{
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.Quit),
			Modifier: gocui.ModNone,
			Handler:  self.handleQuit,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.QuitWithoutChangingDirectory),
			Modifier: gocui.ModNone,
			Handler:  self.handleQuitWithoutChangingDirectory,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.QuitAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.handleQuit,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  self.handleTopLevelReturn,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.ReturnAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.handleTopLevelReturn,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.OpenRecentRepos),
			Handler:     self.handleCreateRecentReposMenu,
			Description: self.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollUpMain),
			Handler:     self.scrollUpMain,
			Alternative: "fn+up/shift+k",
			Description: self.c.Tr.LcScrollUpMainPanel,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ScrollDownMain),
			Handler:     self.scrollDownMain,
			Alternative: "fn+down/shift+j",
			Description: self.c.Tr.LcScrollDownMainPanel,
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
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.CreateRebaseOptionsMenu),
			Handler:     self.helpers.MergeAndRebase.CreateRebaseOptionsMenu,
			Description: self.c.Tr.ViewMergeRebaseOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.CreatePatchOptionsMenu),
			Handler:     self.handleCreatePatchOptionsMenu,
			Description: self.c.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.Refresh),
			Handler:     self.handleRefresh,
			Description: self.c.Tr.LcRefresh,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.OptionMenu),
			Handler:     self.handleCreateOptionsMenu,
			Description: self.c.Tr.LcOpenMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "",
			Key:      opts.GetKey(opts.Config.Universal.OptionMenuAlt1),
			Modifier: gocui.ModNone,
			Handler:  self.handleCreateOptionsMenu,
		},
		{
			ViewName:    "status",
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.handleEditConfig,
			Description: self.c.Tr.EditConfig,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.NextScreenMode),
			Handler:     self.nextScreenMode,
			Description: self.c.Tr.LcNextScreenMode,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.PrevScreenMode),
			Handler:     self.prevScreenMode,
			Description: self.c.Tr.LcPrevScreenMode,
		},
		{
			ViewName:    "status",
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.handleOpenConfig,
			Description: self.c.Tr.OpenConfig,
		},
		{
			ViewName:    "status",
			Key:         opts.GetKey(opts.Config.Status.CheckForUpdate),
			Handler:     self.handleCheckForUpdate,
			Description: self.c.Tr.LcCheckForUpdate,
		},
		{
			ViewName:    "status",
			Key:         opts.GetKey(opts.Config.Status.RecentRepos),
			Handler:     self.handleCreateRecentReposMenu,
			Description: self.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "status",
			Key:         opts.GetKey(opts.Config.Status.AllBranchesLogGraph),
			Handler:     self.handleShowAllBranchLogs,
			Description: self.c.Tr.LcAllBranchesLogGraph,
		},
		{
			ViewName:    "files",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyFileNameToClipboard,
		},
		{
			ViewName:    "localBranches",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyBranchNameToClipboard,
		},
		{
			ViewName:    "commits",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commits",
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     self.helpers.CherryPick.Reset,
			Description: self.c.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "reflogCommits",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "subCommits",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleInfoClick,
		},
		{
			ViewName:    "commitFiles",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitFileNameToClipboard,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.FilteringMenu),
			Handler:     self.handleCreateFilteringMenuPanel,
			Description: self.c.Tr.LcOpenFilteringMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenu),
			Handler:     self.handleCreateDiffingMenuPanel,
			Description: self.c.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.DiffingMenuAlt),
			Handler:     self.handleCreateDiffingMenuPanel,
			Description: self.c.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.ExtrasMenu),
			Handler:     self.handleCreateExtrasMenuPanel,
			Description: self.c.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpSecondary,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownSecondary,
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
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpSecondary,
		},
		{
			ViewName: "status",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleStatusClick,
		},
		{
			ViewName: "search",
			Key:      opts.GetKey(opts.Config.Universal.Confirm),
			Modifier: gocui.ModNone,
			Handler:  self.handleSearch,
		},
		{
			ViewName: "search",
			Key:      opts.GetKey(opts.Config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  self.handleSearchEscape,
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
			ViewName:    "submodules",
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopySubmoduleNameToClipboard,
		},
		{
			ViewName:    "files",
			Key:         opts.GetKey(opts.Config.Universal.ToggleWhitespaceInDiffView),
			Handler:     self.toggleWhitespaceInDiffView,
			Description: self.c.Tr.ToggleWhitespaceInDiffView,
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

	for _, viewName := range []string{"status", "remotes", "tags", "localBranches", "remoteBranches", "files", "submodules", "reflogCommits", "commits", "commitFiles", "subCommits", "stash"} {
		bindings = append(bindings, []*types.Binding{
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.PrevBlock), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.NextBlock), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.PrevBlockAlt), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.NextBlockAlt), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.PrevBlockAlt2), Modifier: gocui.ModNone, Handler: self.previousSideWindow},
			{ViewName: viewName, Key: opts.GetKey(opts.Config.Universal.NextBlockAlt2), Modifier: gocui.ModNone, Handler: self.nextSideWindow},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	windows := []string{"status", "files", "branches", "commits", "stash"}

	if len(config.Universal.JumpToBlock) != len(windows) {
		log.Fatal("Jump to block keybindings cannot be set. Exactly 5 keybindings must be supplied.")
	} else {
		for i, window := range windows {
			bindings = append(bindings, &types.Binding{
				ViewName: "",
				Key:      opts.GetKey(opts.Config.Universal.JumpToBlock[i]),
				Modifier: gocui.ModNone,
				Handler:  self.goToSideWindow(window),
			})
		}
	}

	bindings = append(bindings, []*types.Binding{
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.NextTab),
			Handler:     self.handleNextTab,
			Description: self.c.Tr.LcNextTab,
			Tag:         "navigation",
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.PrevTab),
			Handler:     self.handlePrevTab,
			Description: self.c.Tr.LcPrevTab,
			Tag:         "navigation",
		},
	}...)

	return bindings, mouseKeybindings
}

func (gui *Gui) resetKeybindings() error {
	gui.g.DeleteAllKeybindings()

	bindings, mouseBindings := gui.GetInitialKeybindings()

	// prepending because we want to give our custom keybindings precedence over default keybindings
	customBindings, err := gui.CustomCommandsClient.GetCustomCommandKeybindings()
	if err != nil {
		log.Fatal(err)
	}
	bindings = append(customBindings, bindings...)

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
			tabClickCallback := func(tabIndex int) error { return gui.onViewTabClick(gui.windowForView(viewName), tabIndex) }

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
	handler := binding.Handler
	// TODO: move all mouse-ey stuff into new mouse approach
	if gocui.IsMouseKey(binding.Key) {
		handler = func() error {
			// we ignore click events on views that aren't popup panels, when a popup panel is focused
			if gui.popupPanelFocused() && gui.currentViewName() != binding.ViewName {
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
		// we ignore click events on views that aren't popup panels, when a popup panel is focused
		if gui.popupPanelFocused() && gui.currentViewName() != binding.ViewName {
			return nil
		}

		return baseHandler(opts)
	}
	binding.Handler = newHandler

	return gui.g.SetViewClickBinding(binding)
}
