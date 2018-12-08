package gui

import (
	"github.com/jesseduffield/gocui"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{} // FIXME: find out how to get `gocui.Key | rune`
	Modifier    gocui.Modifier
	KeyReadable string
	Description string
}

// GetDisplayStrings returns the display string of a file
func (b *Binding) GetDisplayStrings() []string {
	return []string{b.GetKey(), b.Description}
}

// GetKey is a function.
func (b *Binding) GetKey() string {
	r, ok := b.Key.(rune)
	key := ""

	if ok {
		key = string(r)
	} else if b.KeyReadable != "" {
		key = b.KeyReadable
	}

	return key
}

// GetKeybindings is a function.
func (gui *Gui) GetKeybindings() []*Binding {
	bindings := []*Binding{
		{
			ViewName: "",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		}, {
			ViewName: "",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		}, {
			ViewName: "",
			Key:      gocui.KeyPgup,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		}, {
			ViewName: "",
			Key:      gocui.KeyPgdn,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlU,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		}, {
			ViewName: "",
			Key:      gocui.KeyCtrlD,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		}, {
			ViewName:    "",
			Key:         'P',
			Modifier:    gocui.ModNone,
			Handler:     gui.pushFiles,
			Description: gui.Tr.SLocalize("push"),
		}, {
			ViewName:    "",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.pullFiles,
			Description: gui.Tr.SLocalize("pull"),
		}, {
			ViewName:    "",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefresh,
			Description: gui.Tr.SLocalize("refresh"),
		}, {
			ViewName: "",
			Key:      'x',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		}, {
			ViewName:    "status",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.SLocalize("EditConfig"),
		}, {
			ViewName:    "status",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.SLocalize("OpenConfig"),
		}, {
			ViewName:    "status",
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.SLocalize("checkForUpdate"),
		}, {
			ViewName:    "status",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRecentReposMenu,
			Description: gui.Tr.SLocalize("SwitchRepo"),
		},
		{
			ViewName:    "files",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.SLocalize("CommitChanges"),
		}, {
			ViewName:    "files",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAmendCommitPress,
			Description: gui.Tr.SLocalize("AmendLastCommit"),
		}, {
			ViewName:    "files",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.SLocalize("CommitChangesWithEditor"),
		}, {
			ViewName:    "files",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFilePress,
			KeyReadable: "space",
			Description: gui.Tr.SLocalize("toggleStaged"),
		}, {
			ViewName:    "files",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileRemove,
			Description: gui.Tr.SLocalize("removeFile"),
		}, {
			ViewName:    "files", // TODO: might make this for more views as well
			Key:         'm',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.SLocalize("ViewMergeRebaseOptions"),
		}, {
			ViewName:    "files",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		}, {
			ViewName:    "files",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		}, {
			ViewName:    "files",
			Key:         'i',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.SLocalize("ignoreFile"),
		}, {
			ViewName:    "files",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.SLocalize("refreshFiles"),
		}, {
			ViewName:    "files",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashSave,
			Description: gui.Tr.SLocalize("stashFiles"),
		}, {
			ViewName:    "files",
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageAll,
			Description: gui.Tr.SLocalize("toggleStagedAll"),
		}, {
			ViewName:    "files",
			Key:         't',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAddPatch,
			Description: gui.Tr.SLocalize("addPatch"),
		}, {
			ViewName:    "files",
			Key:         'D',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleResetAndClean,
			Description: gui.Tr.SLocalize("resetHard"),
		}, {
			ViewName:    "files",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.SLocalize("StageLines"),
			KeyReadable: "enter",
		}, {
			ViewName: "merging",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleEscapeMerge,
		}, {
			ViewName:    "merging",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePickHunk,
			Description: gui.Tr.SLocalize("PickHunk"),
			KeyReadable: "space",
		}, {
			ViewName:    "merging",
			Key:         'b',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePickBothHunks,
			Description: gui.Tr.SLocalize("PickBothHunks"),
		}, {
			ViewName: "merging",
			Key:      gocui.KeyArrowLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		}, {
			ViewName: "merging",
			Key:      gocui.KeyArrowRight,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		}, {
			ViewName: "merging",
			Key:      gocui.KeyArrowUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		}, {
			ViewName: "merging",
			Key:      gocui.KeyArrowDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		}, {
			ViewName: "merging",
			Key:      'h',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		}, {
			ViewName: "merging",
			Key:      'l',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		}, {
			ViewName: "merging",
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		}, {
			ViewName: "merging",
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		}, {
			ViewName:    "merging",
			Key:         'z',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePopFileSnapshot,
			Description: gui.Tr.SLocalize("Undo"),
		}, {
			ViewName:    "branches",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleBranchPress,
			KeyReadable: "space",
			Description: gui.Tr.SLocalize("checkout"),
		}, {
			ViewName:    "branches",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.SLocalize("createPullRequest"),
		}, {
			ViewName:    "branches",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.SLocalize("checkoutByName"),
		}, {
			ViewName:    "branches",
			Key:         'F',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.SLocalize("forceCheckout"),
		}, {
			ViewName:    "branches",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleNewBranch,
			Description: gui.Tr.SLocalize("newBranch"),
		}, {
			ViewName:    "branches",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		}, {
			ViewName:    "branches",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRebase,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		}, {
			ViewName:    "branches",
			Key:         'm',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleMerge,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		}, {
			ViewName:    "branches",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFastForward,
			Description: gui.Tr.SLocalize("FastForward"),
		}, {
			ViewName:    "commits",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.SLocalize("squashDown"),
		}, {
			ViewName:    "commits",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommit,
			Description: gui.Tr.SLocalize("renameCommit"),
		}, {
			ViewName:    "commits",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommitEditor,
			Description: gui.Tr.SLocalize("renameCommitEditor"),
		}, {
			ViewName:    "commits",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleResetToCommit,
			Description: gui.Tr.SLocalize("resetToThisCommit"),
		}, {
			ViewName:    "commits",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.SLocalize("fixupCommit"),
		}, {
			ViewName:    "stash",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashApply,
			KeyReadable: "space",
			Description: gui.Tr.SLocalize("apply"),
		}, {
			ViewName:    "stash",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashPop,
			Description: gui.Tr.SLocalize("pop"),
		}, {
			ViewName:    "stash",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashDrop,
			Description: gui.Tr.SLocalize("drop"),
		}, {
			ViewName: "commitMessage",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitConfirm,
		}, {
			ViewName: "commitMessage",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitClose,
		}, {
			ViewName: "menu",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		}, {
			ViewName: "menu",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		}, {
			ViewName:    "staging",
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStagingEscape,
			KeyReadable: "esc",
			Description: gui.Tr.SLocalize("EscapeStaging"),
		}, {
			ViewName: "staging",
			Key:      gocui.KeyArrowUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingPrevLine,
		}, {
			ViewName: "staging",
			Key:      gocui.KeyArrowDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingNextLine,
		}, {
			ViewName: "staging",
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingPrevLine,
		}, {
			ViewName: "staging",
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingNextLine,
		}, {
			ViewName: "staging",
			Key:      gocui.KeyArrowLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingPrevHunk,
		}, {
			ViewName: "staging",
			Key:      gocui.KeyArrowRight,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingNextHunk,
		}, {
			ViewName: "staging",
			Key:      'h',
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingPrevHunk,
		}, {
			ViewName: "staging",
			Key:      'l',
			Modifier: gocui.ModNone,
			Handler:  gui.handleStagingNextHunk,
		}, {
			ViewName:    "staging",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageLine,
			Description: gui.Tr.SLocalize("StageLine"),
		}, {
			ViewName:    "staging",
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageHunk,
			Description: gui.Tr.SLocalize("StageHunk"),
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: gui.nextView},
		}...)
	}

	listPanelMap := map[string]struct {
		prevLine func(*gocui.Gui, *gocui.View) error
		nextLine func(*gocui.Gui, *gocui.View) error
	}{
		"menu":     {prevLine: gui.handleMenuPrevLine, nextLine: gui.handleMenuNextLine},
		"files":    {prevLine: gui.handleFilesPrevLine, nextLine: gui.handleFilesNextLine},
		"branches": {prevLine: gui.handleBranchesPrevLine, nextLine: gui.handleBranchesNextLine},
		"commits":  {prevLine: gui.handleCommitsPrevLine, nextLine: gui.handleCommitsNextLine},
		"stash":    {prevLine: gui.handleStashPrevLine, nextLine: gui.handleStashNextLine},
	}

	for viewName, functions := range listPanelMap {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: functions.nextLine},
		}...)
	}

	return bindings
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}
