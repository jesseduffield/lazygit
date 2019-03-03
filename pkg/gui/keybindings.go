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
	Description string
}

// GetDisplayStrings returns the display string of a file
func (b *Binding) GetDisplayStrings(isFocused bool) []string {
	return []string{b.GetKey(), b.Description}
}

// GetKey is a function.
func (b *Binding) GetKey() string {
	key := 0

	switch b.Key.(type) {
	case rune:
		key = int(b.Key.(rune))
	case gocui.Key:
		key = int(b.Key.(gocui.Key))
	}

	// special keys
	switch key {
	case 27:
		return "esc"
	case 13:
		return "enter"
	case 32:
		return "space"
	case 65514:
		return "►"
	case 65515:
		return "◄"
	case 65517:
		return "▲"
	case 65516:
		return "▼"
	case 65508:
		return "PgUp"
	case 65507:
		return "PgDn"
	}

	return string(key)
}

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*Binding {
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
			Key:         'm',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.SLocalize("ViewMergeRebaseOptions"),
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
			Description: gui.Tr.SLocalize("toggleStaged"),
		}, {
			ViewName:    "files",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileRemove,
			Description: gui.Tr.SLocalize("removeFile"),
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
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSoftReset,
			Description: gui.Tr.SLocalize("softReset"),
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
		}, {
			ViewName:    "files",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.SLocalize("fetch"),
		}, {
			ViewName:    "branches",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleBranchPress,
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
			Key:         'M',
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
			ViewName:    "commits",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.SLocalize("deleteCommit"),
		}, {
			ViewName:    "commits",
			Key:         'J',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.SLocalize("moveDownCommit"),
		}, {
			ViewName:    "commits",
			Key:         'K',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.SLocalize("moveUpCommit"),
		}, {
			ViewName:    "commits",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.SLocalize("editCommit"),
		}, {
			ViewName:    "commits",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.SLocalize("amendToCommit"),
		}, {
			ViewName:    "commits",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.SLocalize("pickCommit"),
		}, {
			ViewName:    "commits",
			Key:         't',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.SLocalize("revertCommit"),
		}, {
			ViewName:    "commits",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.SLocalize("cherryPickCopy"),
		}, {
			ViewName:    "commits",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.SLocalize("cherryPickCopyRange"),
		}, {
			ViewName:    "commits",
			Key:         'v',
			Modifier:    gocui.ModNone,
			Handler:     gui.HandlePasteCommits,
			Description: gui.Tr.SLocalize("pasteCommits"),
		}, {
			ViewName:    "stash",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashApply,
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
			ViewName: "credentials",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSubmitCredential,
		}, {
			ViewName: "credentials",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCloseCredentialsView,
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
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleDonate,
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
		focus    func(*gocui.Gui, *gocui.View) error
	}{
		"menu":     {prevLine: gui.handleMenuPrevLine, nextLine: gui.handleMenuNextLine, focus: gui.handleMenuSelect},
		"files":    {prevLine: gui.handleFilesPrevLine, nextLine: gui.handleFilesNextLine, focus: gui.handleFilesFocus},
		"branches": {prevLine: gui.handleBranchesPrevLine, nextLine: gui.handleBranchesNextLine, focus: gui.handleBranchSelect},
		"commits":  {prevLine: gui.handleCommitsPrevLine, nextLine: gui.handleCommitsNextLine, focus: gui.handleCommitSelect},
		"stash":    {prevLine: gui.handleStashPrevLine, nextLine: gui.handleStashNextLine, focus: gui.handleStashEntrySelect},
		"status":   {focus: gui.handleStatusSelect},
	}

	for viewName, functions := range listPanelMap {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: functions.prevLine},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: functions.nextLine},
			{ViewName: viewName, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: functions.focus},
		}...)
	}

	return bindings
}

// GetCurrentKeybindings gets the list of keybindings given the current context
func (gui *Gui) GetCurrentKeybindings() []*Binding {
	bindings := gui.GetInitialKeybindings()
	viewName := gui.currentViewName()
	currentContext := gui.State.Contexts[viewName]
	contextBindings := gui.GetContextMap()[viewName][currentContext]

	return append(bindings, contextBindings...)
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	if err := gui.setInitialContexts(); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) GetContextMap() map[string]map[string][]*Binding {
	return map[string]map[string][]*Binding{
		"main": {
			"normal": {
				{
					ViewName:    "main",
					Key:         gocui.MouseWheelDown,
					Modifier:    gocui.ModNone,
					Handler:     gui.scrollDownMain,
					Description: gui.Tr.SLocalize("ScrollDown"),
				}, {
					ViewName:    "main",
					Key:         gocui.MouseWheelUp,
					Modifier:    gocui.ModNone,
					Handler:     gui.scrollUpMain,
					Description: gui.Tr.SLocalize("ScrollUp"),
				},
			},
			"staging": {
				{
					ViewName:    "main",
					Key:         gocui.KeyEsc,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStagingEscape,
					Description: gui.Tr.SLocalize("EscapeStaging"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowUp,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStagingPrevLine,
					Description: gui.Tr.SLocalize("PrevLine"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowDown,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStagingNextLine,
					Description: gui.Tr.SLocalize("NextLine"),
				}, {
					ViewName: "main",
					Key:      'k',
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingPrevLine,
				}, {
					ViewName: "main",
					Key:      'j',
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingNextLine,
				}, {
					ViewName: "main",
					Key:      gocui.MouseWheelUp,
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingPrevLine,
				}, {
					ViewName: "main",
					Key:      gocui.MouseWheelDown,
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingNextLine,
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowLeft,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStagingPrevHunk,
					Description: gui.Tr.SLocalize("PrevHunk"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowRight,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStagingNextHunk,
					Description: gui.Tr.SLocalize("NextHunk"),
				}, {
					ViewName: "main",
					Key:      'h',
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingPrevHunk,
				}, {
					ViewName: "main",
					Key:      'l',
					Modifier: gocui.ModNone,
					Handler:  gui.handleStagingNextHunk,
				}, {
					ViewName:    "main",
					Key:         gocui.KeySpace,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStageLine,
					Description: gui.Tr.SLocalize("StageLine"),
				}, {
					ViewName:    "main",
					Key:         'a',
					Modifier:    gocui.ModNone,
					Handler:     gui.handleStageHunk,
					Description: gui.Tr.SLocalize("StageHunk"),
				},
			},
			"merging": {
				{
					ViewName:    "main",
					Key:         gocui.KeyEsc,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleEscapeMerge,
					Description: gui.Tr.SLocalize("EscapeStaging"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeySpace,
					Modifier:    gocui.ModNone,
					Handler:     gui.handlePickHunk,
					Description: gui.Tr.SLocalize("PickHunk"),
				}, {
					ViewName:    "main",
					Key:         'b',
					Modifier:    gocui.ModNone,
					Handler:     gui.handlePickBothHunks,
					Description: gui.Tr.SLocalize("PickBothHunks"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowLeft,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleSelectPrevConflict,
					Description: gui.Tr.SLocalize("PrevConflict"),
				}, {
					ViewName: "main",
					Key:      gocui.KeyArrowRight,
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectNextConflict,

					Description: gui.Tr.SLocalize("NextConflict"),
				}, {
					ViewName: "main",
					Key:      gocui.KeyArrowUp,

					Modifier:    gocui.ModNone,
					Handler:     gui.handleSelectTop,
					Description: gui.Tr.SLocalize("SelectTop"),
				}, {
					ViewName:    "main",
					Key:         gocui.KeyArrowDown,
					Modifier:    gocui.ModNone,
					Handler:     gui.handleSelectBottom,
					Description: gui.Tr.SLocalize("SelectBottom"),
				}, {
					ViewName: "main",
					Key:      gocui.MouseWheelUp,
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectTop,
				}, {
					ViewName: "main",
					Key:      gocui.MouseWheelDown,
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectBottom,
				}, {
					ViewName: "main",
					Key:      'h',
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectPrevConflict,
				}, {
					ViewName: "main",
					Key:      'l',
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectNextConflict,
				}, {
					ViewName: "main",
					Key:      'k',
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectTop,
				}, {
					ViewName: "main",
					Key:      'j',
					Modifier: gocui.ModNone,
					Handler:  gui.handleSelectBottom,
				}, {
					ViewName:    "main",
					Key:         'z',
					Modifier:    gocui.ModNone,
					Handler:     gui.handlePopFileSnapshot,
					Description: gui.Tr.SLocalize("Undo"),
				},
			},
		},
	}
}
