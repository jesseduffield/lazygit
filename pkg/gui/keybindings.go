package gui

import (
	"log"
	"strings"
	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) getKeyDisplay(name string) string {
	key := gui.getKey(name)
	return keybindings.GetKeyDisplay(key)
}

func (gui *Gui) getKey(key string) types.Key {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding := keybindings.Keymap[strings.ToLower(key)]
		if binding == nil {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see %s", strings.ToLower(key), constants.Links.Docs.CustomKeybindings)
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return []rune(key)[0]
	}
	log.Fatal("Key empty for keybinding: " + strings.ToLower(key))
	return nil
}

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
		GetKey: self.getKey,
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
			Contexts:    []string{string(context.STATUS_CONTEXT_KEY)},
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
			Contexts:    []string{string(context.STATUS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.handleOpenConfig,
			Description: self.c.Tr.OpenConfig,
		},
		{
			ViewName:    "status",
			Contexts:    []string{string(context.STATUS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Status.CheckForUpdate),
			Handler:     self.handleCheckForUpdate,
			Description: self.c.Tr.LcCheckForUpdate,
		},
		{
			ViewName:    "status",
			Contexts:    []string{string(context.STATUS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Status.RecentRepos),
			Handler:     self.handleCreateRecentReposMenu,
			Description: self.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "status",
			Contexts:    []string{string(context.STATUS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Status.AllBranchesLogGraph),
			Handler:     self.handleShowAllBranchLogs,
			Description: self.c.Tr.LcAllBranchesLogGraph,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyFileNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyBranchNameToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.LOCAL_COMMITS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.LOCAL_COMMITS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Commits.ResetCherryPick),
			Handler:     self.helpers.CherryPick.Reset,
			Description: self.c.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "subCommits",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
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
			Contexts:    []string{string(context.COMMIT_FILES_CONTEXT_KEY)},
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
			Contexts:    []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:         gocui.MouseWheelDown,
			Handler:     self.scrollDownMain,
			Description: self.c.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:         gocui.MouseWheelUp,
			Handler:     self.scrollUpMain,
			Description: self.c.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "secondary",
			Contexts: []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleTogglePanelClick,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.handleStagingEscape,
			Description: self.c.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.handleToggleStagedSelection,
			Description: self.c.Tr.StageSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.handleResetSelection,
			Description: self.c.Tr.ResetSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler:     self.handleTogglePanel,
			Description: self.c.Tr.TogglePanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.handleEscapePatchBuildingPanel,
			Description: self.c.Tr.ExitLineByLineMode,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.handleOpenFileAtLine,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.PrevItem),
			Handler:     self.handleSelectPrevLine,
			Description: self.c.Tr.PrevLine,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.NextItem),
			Handler:     self.handleSelectNextLine,
			Description: self.c.Tr.NextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectNextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.PrevBlock),
			Handler:     self.handleSelectPrevHunk,
			Description: self.c.Tr.PrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.NextBlock),
			Handler:     self.handleSelectNextHunk,
			Description: self.c.Tr.NextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectNextHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Modifier:    gocui.ModNone,
			Handler:     self.copySelectedToClipboard,
			Description: self.c.Tr.LcCopySelectedTexToClipboard,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.handleLineByLineEdit,
			Description: self.c.Tr.LcEditFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.HandleOpenFile,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.NextPage),
			Modifier:    gocui.ModNone,
			Handler:     self.handleLineByLineNextPage,
			Description: self.c.Tr.LcNextPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.PrevPage),
			Modifier:    gocui.ModNone,
			Handler:     self.handleLineByLinePrevPage,
			Description: self.c.Tr.LcPrevPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.GotoTop),
			Modifier:    gocui.ModNone,
			Handler:     self.handleLineByLineGotoTop,
			Description: self.c.Tr.LcGotoTop,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.GotoBottom),
			Modifier:    gocui.ModNone,
			Handler:     self.handleLineByLineGotoBottom,
			Description: self.c.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     func() error { return self.handleOpenSearch("main") },
			Description: self.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.handleToggleSelectionForPatch,
			Description: self.c.Tr.ToggleSelectionForPatch,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Main.ToggleDragSelect),
			Handler:     self.handleToggleSelectRange,
			Description: self.c.Tr.ToggleDragSelect,
		},
		// Alias 'V' -> 'v'
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Main.ToggleDragSelectAlt),
			Handler:     self.handleToggleSelectRange,
			Description: self.c.Tr.ToggleDragSelect,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Main.ToggleSelectHunk),
			Handler:     self.handleToggleSelectHunk,
			Description: self.c.Tr.ToggleSelectHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Main.EditSelectHunk),
			Handler:     self.handleEditHunk,
			Description: self.c.Tr.EditHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  self.handleLBLMouseDown,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModMotion,
			Handler:  self.handleMouseDrag,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY), string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.ScrollLeft),
			Handler:     self.scrollLeftMain,
			Description: self.c.Tr.LcScrollLeft,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY), string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.ScrollRight),
			Handler:     self.scrollRightMain,
			Description: self.c.Tr.LcScrollRight,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.handleEscapeMerge,
			Description: self.c.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Files.OpenMergeTool),
			Handler:     self.helpers.WorkingTree.OpenMergeTool,
			Description: self.c.Tr.LcOpenMergeTool,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.handlePickHunk,
			Description: self.c.Tr.PickHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Main.PickBothHunks),
			Handler:     self.handlePickAllHunks,
			Description: self.c.Tr.PickAllHunks,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.PrevBlock),
			Handler:     self.handleSelectPrevConflict,
			Description: self.c.Tr.PrevConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.NextBlock),
			Handler:     self.handleSelectNextConflict,
			Description: self.c.Tr.NextConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.PrevItem),
			Handler:     self.handleSelectPrevConflictHunk,
			Description: self.c.Tr.SelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.NextItem),
			Handler:     self.handleSelectNextConflictHunk,
			Description: self.c.Tr.SelectNextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectPrevConflictHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.handleSelectNextConflictHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.handleMergeConflictEditFileAtLine,
			Description: self.c.Tr.LcEditFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.handleMergeConflictOpenFileAtLine,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.Undo),
			Handler:     self.handleMergeConflictUndo,
			Description: self.c.Tr.LcUndo,
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
			ViewName:    "files",
			Contexts:    []string{string(context.SUBMODULES_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.CopyToClipboard),
			Handler:     self.handleCopySelectedSideContextItemToClipboard,
			Description: self.c.Tr.LcCopySubmoduleNameToClipboard,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         opts.GetKey(opts.Config.Universal.ToggleWhitespaceInDiffView),
			Handler:     self.toggleWhitespaceInDiffView,
			Description: self.c.Tr.ToggleWhitespaceInDiffView,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.IncreaseContextInDiffView),
			Handler:     self.IncreaseContextInDiffView,
			Description: self.c.Tr.IncreaseContextInDiffView,
		},
		{
			ViewName:    "",
			Key:         opts.GetKey(opts.Config.Universal.DecreaseContextInDiffView),
			Handler:     self.DecreaseContextInDiffView,
			Description: self.c.Tr.DecreaseContextInDiffView,
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
			ViewName:    "extras",
			Key:         opts.GetKey(opts.Config.Universal.ExtrasMenu),
			Handler:     self.handleCreateExtrasMenuPanel,
			Description: self.c.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      opts.GetKey(opts.Config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  self.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
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
		contextKey := c.GetKey()
		for _, binding := range c.GetKeybindings(opts) {
			// TODO: move all mouse keybindings into the mouse keybindings approach below
			if !gocui.IsMouseKey(binding.Key) && contextKey != context.GLOBAL_CONTEXT_KEY {
				binding.Contexts = []string{string(contextKey)}
			}
			binding.ViewName = viewName
			bindings = append(bindings, binding)
		}

		mouseKeybindings = append(mouseKeybindings, c.GetMouseKeybindings(opts)...)
	}

	for _, viewName := range []string{"status", "branches", "remoteBranches", "files", "commits", "commitFiles", "subCommits", "stash"} {
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

	for viewName := range self.initialViewTabContextMap(self.State.Contexts) {
		bindings = append(bindings, []*types.Binding{
			{
				ViewName:    viewName,
				Key:         opts.GetKey(opts.Config.Universal.NextTab),
				Handler:     self.handleNextTab,
				Description: self.c.Tr.LcNextTab,
				Tag:         "navigation",
			},
			{
				ViewName:    viewName,
				Key:         opts.GetKey(opts.Config.Universal.PrevTab),
				Handler:     self.handlePrevTab,
				Description: self.c.Tr.LcPrevTab,
				Tag:         "navigation",
			},
		}...)
	}

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

	for viewName := range gui.initialViewTabContextMap(gui.State.Contexts) {
		viewName := viewName
		tabClickCallback := func(tabIndex int) error { return gui.onViewTabClick(viewName, tabIndex) }

		if err := gui.g.SetTabClickBinding(viewName, tabClickCallback); err != nil {
			return err
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

	return gui.g.SetKeybinding(binding.ViewName, binding.Contexts, binding.Key, binding.Modifier, gui.wrappedHandler(handler))
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
