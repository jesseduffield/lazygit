package gui

import (
	"github.com/jesseduffield/gocui"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Contexts    []string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{} // FIXME: find out how to get `gocui.Key | rune`
	Modifier    gocui.Modifier
	Description string
	Alternative string
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
		if b.Key.(gocui.Key) == gocui.KeyCtrlJ {
			return "ctrl+j"
		}
		if b.Key.(gocui.Key) == gocui.KeyCtrlK {
			return "ctrl+k"
		}
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
	case 9:
		return "tab"
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
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			Key:      'Q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuitWithoutChangingDirectory,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		},
		{
			ViewName:    "",
			Key:         gocui.KeyPgup,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
		},
		{
			ViewName:    "",
			Key:         gocui.KeyPgdn,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
		},
		{
			ViewName: "",
			Key:      'K',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      'J',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlU,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlD,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "",
			Key:         'm',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.SLocalize("ViewMergeRebaseOptions"),
		},
		{
			ViewName:    "",
			Key:         'P',
			Modifier:    gocui.ModNone,
			Handler:     gui.pushFiles,
			Description: gui.Tr.SLocalize("push"),
		},
		{
			ViewName:    "",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePullFiles,
			Description: gui.Tr.SLocalize("pull"),
		},
		{
			ViewName:    "",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefresh,
			Description: gui.Tr.SLocalize("refresh"),
		},
		{
			ViewName: "",
			Key:      'x',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName: "",
			Key:      '?',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName: "",
			Key:      gocui.MouseMiddle,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlP,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreatePatchOptionsMenu,
		},
		{
			ViewName:    "status",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.SLocalize("EditConfig"),
		},
		{
			ViewName:    "status",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.SLocalize("OpenConfig"),
		},
		{
			ViewName:    "status",
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.SLocalize("checkForUpdate"),
		},
		{
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
		},
		{
			ViewName:    "files",
			Key:         'w',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.SLocalize("commitChangesWithoutHook"),
		},
		{
			ViewName:    "files",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAmendCommitPress,
			Description: gui.Tr.SLocalize("AmendLastCommit"),
		},
		{
			ViewName:    "files",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.SLocalize("CommitChangesWithEditor"),
		},
		{
			ViewName:    "files",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFilePress,
			Description: gui.Tr.SLocalize("toggleStaged"),
		},
		{
			ViewName:    "files",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.Tr.SLocalize("viewDiscardOptions"),
		},
		{
			ViewName:    "files",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		},
		{
			ViewName:    "files",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "files",
			Key:         'i',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.SLocalize("ignoreFile"),
		},
		{
			ViewName:    "files",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.SLocalize("refreshFiles"),
		},
		{
			ViewName:    "files",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashChanges,
			Description: gui.Tr.SLocalize("stashAllChanges"),
		},
		{
			ViewName:    "files",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateStashMenu,
			Description: gui.Tr.SLocalize("viewStashOptions"),
		},
		{
			ViewName:    "files",
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageAll,
			Description: gui.Tr.SLocalize("toggleStagedAll"),
		},
		{
			ViewName:    "files",
			Key:         'D',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateResetMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "files",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.SLocalize("StageLines"),
		},
		{
			ViewName:    "files",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.SLocalize("fetch"),
		},
		{
			ViewName:    "files",
			Key:         'X',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCustomCommand,
			Description: gui.Tr.SLocalize("executeCustomCommand"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleBranchPress,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.SLocalize("createPullRequest"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.SLocalize("checkoutByName"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'F',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.SLocalize("forceCheckout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleNewBranch,
			Description: gui.Tr.SLocalize("newBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRebaseOntoLocalBranch,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'M',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleMerge,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFastForward,
			Description: gui.Tr.SLocalize("FastForward"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutTag,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDeleteTag,
			Description: gui.Tr.SLocalize("deleteTag"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         'P',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePushTag,
			Description: gui.Tr.SLocalize("pushTag"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateTag,
			Description: gui.Tr.SLocalize("createTag"),
		},
		{
			ViewName: "branches",
			Key:      ']',
			Modifier: gocui.ModNone,
			Handler:  gui.handleNextBranchesTab,
		},
		{
			ViewName: "branches",
			Key:      '[',
			Modifier: gocui.ModNone,
			Handler:  gui.handlePrevBranchesTab,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRemoteBranchesEscape,
			Description: gui.Tr.SLocalize("ReturnToRemotesList"),
		},
		{
			ViewName:    "commits",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.SLocalize("squashDown"),
		},
		{
			ViewName:    "commits",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommit,
			Description: gui.Tr.SLocalize("renameCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'R',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRenameCommitEditor,
			Description: gui.Tr.SLocalize("renameCommitEditor"),
		},
		{
			ViewName:    "commits",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateCommitResetMenu,
			Description: gui.Tr.SLocalize("resetToThisCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'f',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.SLocalize("fixupCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'F',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCreateFixupCommit,
			Description: gui.Tr.SLocalize("createFixupCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSquashAllAboveFixupCommits,
			Description: gui.Tr.SLocalize("squashAboveCommits"),
		},
		{
			ViewName:    "commits",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.SLocalize("deleteCommit"),
		},
		{
			ViewName:    "commits",
			Key:         gocui.KeyCtrlJ,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.SLocalize("moveDownCommit"),
		},
		{
			ViewName:    "commits",
			Key:         gocui.KeyCtrlK,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.SLocalize("moveUpCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.SLocalize("editCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'A',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.SLocalize("amendToCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'p',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.SLocalize("pickCommit"),
		},
		{
			ViewName:    "commits",
			Key:         't',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.SLocalize("revertCommit"),
		},
		{
			ViewName:    "commits",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.SLocalize("cherryPickCopy"),
		},
		{
			ViewName:    "commits",
			Key:         'C',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.SLocalize("cherryPickCopyRange"),
		},
		{
			ViewName:    "commits",
			Key:         'v',
			Modifier:    gocui.ModNone,
			Handler:     gui.HandlePasteCommits,
			Description: gui.Tr.SLocalize("pasteCommits"),
		},
		{
			ViewName:    "commits",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSwitchToCommitFilesPanel,
			Description: gui.Tr.SLocalize("viewCommitFiles"),
		},
		{
			ViewName:    "commits",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleDiffCommit,
			Description: gui.Tr.SLocalize("CommitsDiff"),
		},
		{
			ViewName:    "commits",
			Key:         'T',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleTagCommit,
			Description: gui.Tr.SLocalize("tagCommit"),
		},
		{
			ViewName:    "stash",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashApply,
			Description: gui.Tr.SLocalize("apply"),
		},
		{
			ViewName:    "stash",
			Key:         'g',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashPop,
			Description: gui.Tr.SLocalize("pop"),
		},
		{
			ViewName:    "stash",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStashDrop,
			Description: gui.Tr.SLocalize("drop"),
		},
		{
			ViewName: "commitMessage",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitConfirm,
		},
		{
			ViewName: "commitMessage",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitClose,
		},
		{
			ViewName: "credentials",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSubmitCredential,
		},
		{
			ViewName: "credentials",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCloseCredentialsView,
		},
		{
			ViewName: "menu",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		},
		{
			ViewName: "menu",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.handleMenuClose,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleDonate,
		},
		{
			ViewName:    "commitFiles",
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSwitchToCommitsPanel,
			Description: gui.Tr.SLocalize("goBack"),
		},
		{
			ViewName:    "commitFiles",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.Tr.SLocalize("checkoutCommitFile"),
		},
		{
			ViewName:    "commitFiles",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.Tr.SLocalize("discardOldFileChange"),
		},
		{
			ViewName:    "commitFiles",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.Tr.SLocalize("toggleAddToPatch"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gocui.KeyEnter,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEnterCommitFile,
			Description: gui.Tr.SLocalize("enterFile"),
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpSecondary,
		},
		{
			ViewName: "secondary",
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownSecondary,
		},
		{
			ViewName: "secondary",
			Contexts: []string{"normal"},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownSecondary,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"normal"},
			Key:         gocui.MouseWheelDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollDownMain,
			Description: gui.Tr.SLocalize("ScrollDown"),
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{"normal"},
			Key:         gocui.MouseWheelUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.scrollUpMain,
			Description: gui.Tr.SLocalize("ScrollUp"),
			Alternative: "fn+down",
		},
		{
			ViewName: "main",
			Contexts: []string{"normal"},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownMain,
		},
		{
			ViewName: "secondary",
			Contexts: []string{"staging"},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleTogglePanelClick,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStagingEscape,
			Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleStageSelection,
			Description: gui.Tr.SLocalize("StageSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleResetSelection,
			Description: gui.Tr.SLocalize("ResetSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gocui.KeyTab,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleTogglePanel,
			Description: gui.Tr.SLocalize("TogglePanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building"},
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEscapePatchBuildingPanel,
			Description: gui.Tr.SLocalize("ExitLineByLineMode"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectPrevLine,
			Description: gui.Tr.SLocalize("PrevLine"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectNextLine,
			Description: gui.Tr.SLocalize("NextLine"),
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextLine,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gocui.KeyArrowLeft,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectPrevHunk,
			Description: gui.Tr.SLocalize("PrevHunk"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gocui.KeyArrowRight,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectNextHunk,
			Description: gui.Tr.SLocalize("NextHunk"),
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      'h',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      'l',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAddSelectionToPatch,
			Description: gui.Tr.SLocalize("StageSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRemoveSelectionFromPatch,
			Description: gui.Tr.SLocalize("ResetSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         'v',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleSelectRange,
			Description: gui.Tr.SLocalize("ToggleDragSelect"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         'a',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleToggleSelectHunk,
			Description: gui.Tr.SLocalize("ToggleSelectHunk"),
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDown,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModMotion,
			Handler:  gui.handleMouseDrag,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseScrollUp,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseScrollDown,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEscapeMerge,
			Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePickHunk,
			Description: gui.Tr.SLocalize("PickHunk"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         'b',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePickBothHunks,
			Description: gui.Tr.SLocalize("PickBothHunks"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeyArrowLeft,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectPrevConflict,
			Description: gui.Tr.SLocalize("PrevConflict"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeyArrowRight,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectNextConflict,
			Description: gui.Tr.SLocalize("NextConflict"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeyArrowUp,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectTop,
			Description: gui.Tr.SLocalize("SelectTop"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gocui.KeyArrowDown,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSelectBottom,
			Description: gui.Tr.SLocalize("SelectBottom"),
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      'h',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      'l',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         'z',
			Modifier:    gocui.ModNone,
			Handler:     gui.handlePopFileSnapshot,
			Description: gui.Tr.SLocalize("Undo"),
		},
		{
			ViewName: "branches",
			Contexts: []string{"remotes"},
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleRemoteEnter,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleAddRemote,
			Description: gui.Tr.SLocalize("addNewRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRemoveRemote,
			Description: gui.Tr.SLocalize("removeRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleEditRemote,
			Description: gui.Tr.SLocalize("editRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gocui.KeySpace,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCheckoutRemoteBranch,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         'M',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleMergeRemoteBranch,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDeleteRemoteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleRebaseOntoRemoteBranch,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         'u',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleSetBranchUpstream,
			Description: gui.Tr.SLocalize("setUpstream"),
		},
		{
			ViewName: "stash",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStashEntrySelect,
		},
		{
			ViewName: "status",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStatusClick,
		},
		{
			ViewName: "commitFiles",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitFilesClick,
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: gui.nextView},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	for i, viewName := range []string{"status", "files", "branches", "commits", "stash"} {
		bindings = append(bindings, &Binding{ViewName: "", Key: rune(i+1) + '0', Modifier: gocui.ModNone, Handler: gui.goToSideView(viewName)})
	}

	for _, listView := range gui.getListViews() {
		bindings = append(bindings, []*Binding{
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: 'k', Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: 'j', Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: listView.handleClick},
		}...)
	}

	return bindings
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Contexts, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	if err := g.SetTabClickBinding("branches", gui.onBranchesTabClick); err != nil {
		return err
	}

	return nil
}
