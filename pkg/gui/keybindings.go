package gui

import (
	"log"
	"strings"

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
	return []string{GetKeyDisplay(b.Key), b.Description}
}

var keyMapReversed = map[gocui.Key]string{
	gocui.KeyF1:         "f1",
	gocui.KeyF2:         "f2",
	gocui.KeyF3:         "f3",
	gocui.KeyF4:         "f4",
	gocui.KeyF5:         "f5",
	gocui.KeyF6:         "f6",
	gocui.KeyF7:         "f7",
	gocui.KeyF8:         "f8",
	gocui.KeyF9:         "f9",
	gocui.KeyF10:        "f10",
	gocui.KeyF11:        "f11",
	gocui.KeyF12:        "f12",
	gocui.KeyInsert:     "insert",
	gocui.KeyDelete:     "delete",
	gocui.KeyHome:       "home",
	gocui.KeyEnd:        "end",
	gocui.KeyPgup:       "pgup",
	gocui.KeyPgdn:       "pgdown",
	gocui.KeyArrowUp:    "▲",
	gocui.KeyArrowDown:  "▼",
	gocui.KeyArrowLeft:  "◄",
	gocui.KeyArrowRight: "►",
	gocui.KeyTab:        "tab",        // ctrl+i
	gocui.KeyEnter:      "enter",      // ctrl+m
	gocui.KeyEsc:        "esc",        // ctrl+[, ctrl+3
	gocui.KeyBackspace:  "backspace",  // ctrl+h
	gocui.KeyCtrlSpace:  "ctrl+space", // ctrl+~, ctrl+2
	gocui.KeyCtrlSlash:  "ctrl+/",     // ctrl+_
	gocui.KeySpace:      "space",
	gocui.KeyCtrlA:      "ctrl+a",
	gocui.KeyCtrlB:      "ctrl+b",
	gocui.KeyCtrlC:      "ctrl+c",
	gocui.KeyCtrlD:      "ctrl+d",
	gocui.KeyCtrlE:      "ctrl+e",
	gocui.KeyCtrlF:      "ctrl+f",
	gocui.KeyCtrlG:      "ctrl+g",
	gocui.KeyCtrlJ:      "ctrl+j",
	gocui.KeyCtrlK:      "ctrl+k",
	gocui.KeyCtrlL:      "ctrl+l",
	gocui.KeyCtrlN:      "ctrl+n",
	gocui.KeyCtrlO:      "ctrl+o",
	gocui.KeyCtrlP:      "ctrl+p",
	gocui.KeyCtrlQ:      "ctrl+q",
	gocui.KeyCtrlR:      "ctrl+r",
	gocui.KeyCtrlS:      "ctrl+s",
	gocui.KeyCtrlT:      "ctrl+t",
	gocui.KeyCtrlU:      "ctrl+u",
	gocui.KeyCtrlV:      "ctrl+v",
	gocui.KeyCtrlW:      "ctrl+w",
	gocui.KeyCtrlX:      "ctrl+x",
	gocui.KeyCtrlY:      "ctrl+y",
	gocui.KeyCtrlZ:      "ctrl+z",
	gocui.KeyCtrl4:      "ctrl+4", // ctrl+\
	gocui.KeyCtrl5:      "ctrl+5", // ctrl+]
	gocui.KeyCtrl6:      "ctrl+6",
	gocui.KeyCtrl8:      "ctrl+8",
}

var keymap = map[string]interface{}{
	"<c-a>":       gocui.KeyCtrlA,
	"<c-b>":       gocui.KeyCtrlB,
	"<c-c>":       gocui.KeyCtrlC,
	"<c-d>":       gocui.KeyCtrlD,
	"<c-e>":       gocui.KeyCtrlE,
	"<c-f>":       gocui.KeyCtrlF,
	"<c-g>":       gocui.KeyCtrlG,
	"<c-h>":       gocui.KeyCtrlH,
	"<c-i>":       gocui.KeyCtrlI,
	"<c-j>":       gocui.KeyCtrlJ,
	"<c-k>":       gocui.KeyCtrlK,
	"<c-l>":       gocui.KeyCtrlL,
	"<c-m>":       gocui.KeyCtrlM,
	"<c-n>":       gocui.KeyCtrlN,
	"<c-o>":       gocui.KeyCtrlO,
	"<c-p>":       gocui.KeyCtrlP,
	"<c-q>":       gocui.KeyCtrlQ,
	"<c-r>":       gocui.KeyCtrlR,
	"<c-s>":       gocui.KeyCtrlS,
	"<c-t>":       gocui.KeyCtrlT,
	"<c-u>":       gocui.KeyCtrlU,
	"<c-v>":       gocui.KeyCtrlV,
	"<c-w>":       gocui.KeyCtrlW,
	"<c-x>":       gocui.KeyCtrlX,
	"<c-y>":       gocui.KeyCtrlY,
	"<c-z>":       gocui.KeyCtrlZ,
	"<c-~>":       gocui.KeyCtrlTilde,
	"<c-2>":       gocui.KeyCtrl2,
	"<c-3>":       gocui.KeyCtrl3,
	"<c-4>":       gocui.KeyCtrl4,
	"<c-5>":       gocui.KeyCtrl5,
	"<c-6>":       gocui.KeyCtrl6,
	"<c-7>":       gocui.KeyCtrl7,
	"<c-8>":       gocui.KeyCtrl8,
	"<c-space>":   gocui.KeyCtrlSpace,
	"<c-\\>":      gocui.KeyCtrlBackslash,
	"<c-[>":       gocui.KeyCtrlLsqBracket,
	"<c-]>":       gocui.KeyCtrlRsqBracket,
	"<c-/>":       gocui.KeyCtrlSlash,
	"<c-_>":       gocui.KeyCtrlUnderscore,
	"<backspace>": gocui.KeyBackspace,
	"<tab>":       gocui.KeyTab,
	"<enter>":     gocui.KeyEnter,
	"<esc>":       gocui.KeyEsc,
	"<space>":     gocui.KeySpace,
	"<f1>":        gocui.KeyF1,
	"<f2>":        gocui.KeyF2,
	"<f3>":        gocui.KeyF3,
	"<f4>":        gocui.KeyF4,
	"<f5>":        gocui.KeyF5,
	"<f6>":        gocui.KeyF6,
	"<f7>":        gocui.KeyF7,
	"<f8>":        gocui.KeyF8,
	"<f9>":        gocui.KeyF9,
	"<f10>":       gocui.KeyF10,
	"<f11>":       gocui.KeyF11,
	"<f12>":       gocui.KeyF12,
	"<insert>":    gocui.KeyInsert,
	"<delete>":    gocui.KeyDelete,
	"<home>":      gocui.KeyHome,
	"<end>":       gocui.KeyEnd,
	"<pgup>":      gocui.KeyPgup,
	"<pgdown>":    gocui.KeyPgdn,
	"<up>":        gocui.KeyArrowUp,
	"<down>":      gocui.KeyArrowDown,
	"<left>":      gocui.KeyArrowLeft,
	"<right>":     gocui.KeyArrowRight,
}

func (gui *Gui) getKeyDisplay(name string) string {
	key := gui.getKey(name)
	return GetKeyDisplay(key)
}

func GetKeyDisplay(key interface{}) string {
	keyInt := 0

	switch key := key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		value, ok := keyMapReversed[key]
		if ok {
			return value
		}
		keyInt = int(key)
	}

	return string(keyInt)
}

func (gui *Gui) getKey(name string) interface{} {
	key := gui.Config.GetUserConfig().GetString("keybinding." + name)
	if len(key) > 1 {
		binding := keymap[strings.ToLower(key)]
		if binding == nil {
			log.Fatalf("Unrecognized key %s for keybinding %s", strings.ToLower(key), name)
		} else {
			return binding
		}
	} else if len(key) == 1 {
		return []rune(key)[0]
	}
	log.Fatal("Key empty for keybinding: " + strings.ToLower(name))
	return nil
}

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*Binding {
	bindings := []*Binding{
		{
			ViewName: "",
			Key:      gui.getKey("universal.quit"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.quitWithoutChangingDirectory"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuitWithoutChangingDirectory,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.quit-alt1"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.return"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.scrollUpMain"),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
			Description: gui.Tr.SLocalize("scrollUpMainPanel"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.scrollDownMain"),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
			Description: gui.Tr.SLocalize("scrollDownMainPanel"),
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.scrollUpMain-alt1"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.scrollDownMain-alt1"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.scrollUpMain-alt2"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.scrollDownMain-alt2"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.createRebaseOptionsMenu"),
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.SLocalize("ViewMergeRebaseOptions"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.createPatchOptionsMenu"),
			Handler:     gui.handleCreatePatchOptionsMenu,
			Description: gui.Tr.SLocalize("ViewPatchOptions"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.pushFiles"),
			Handler:     gui.pushFiles,
			Description: gui.Tr.SLocalize("push"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.pullFiles"),
			Handler:     gui.handlePullFiles,
			Description: gui.Tr.SLocalize("pull"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.refresh"),
			Handler:     gui.handleRefresh,
			Description: gui.Tr.SLocalize("refresh"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.optionMenu"),
			Handler:     gui.handleCreateOptionsMenu,
			Description: gui.Tr.SLocalize("openMenu"),
		},
		{
			ViewName: "",
			Key:      gui.getKey("universal.optionMenu-alt1"),
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
			ViewName:    "",
			Key:         gui.getKey("universal.undo"),
			Handler:     gui.reflogUndo,
			Description: gui.Tr.SLocalize("undoReflog"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.redo"),
			Handler:     gui.reflogRedo,
			Description: gui.Tr.SLocalize("redoReflog"),
		},
		{
			ViewName:    "status",
			Key:         gui.getKey("universal.edit"),
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.SLocalize("EditConfig"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.nextScreenMode"),
			Handler:     gui.nextScreenMode,
			Description: gui.Tr.SLocalize("nextScreenMode"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.prevScreenMode"),
			Handler:     gui.prevScreenMode,
			Description: gui.Tr.SLocalize("prevScreenMode"),
		},
		{
			ViewName:    "status",
			Key:         gui.getKey("universal.openFile"),
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.SLocalize("OpenConfig"),
		},
		{
			ViewName:    "status",
			Key:         gui.getKey("status.checkForUpdate"),
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.SLocalize("checkForUpdate"),
		},
		{
			ViewName:    "status",
			Key:         gui.getKey("status.recentRepos"),
			Handler:     gui.handleCreateRecentReposMenu,
			Description: gui.Tr.SLocalize("SwitchRepo"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.commitChanges"),
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.SLocalize("CommitChanges"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.commitChangesWithoutHook"),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.SLocalize("commitChangesWithoutHook"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.amendLastCommit"),
			Handler:     gui.handleAmendCommitPress,
			Description: gui.Tr.SLocalize("AmendLastCommit"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.commitChangesWithEditor"),
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.SLocalize("CommitChangesWithEditor"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleFilePress,
			Description: gui.Tr.SLocalize("toggleStaged"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.Tr.SLocalize("viewDiscardOptions"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("universal.edit"),
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("universal.openFile"),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.ignoreFile"),
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.SLocalize("ignoreFile"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.refreshFiles"),
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.SLocalize("refreshFiles"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.stashAllChanges"),
			Handler:     gui.handleStashChanges,
			Description: gui.Tr.SLocalize("stashAllChanges"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.viewStashOptions"),
			Handler:     gui.handleCreateStashMenu,
			Description: gui.Tr.SLocalize("viewStashOptions"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.toggleStagedAll"),
			Handler:     gui.handleStageAll,
			Description: gui.Tr.SLocalize("toggleStagedAll"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.viewResetOptions"),
			Handler:     gui.handleCreateResetMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("universal.goInto"),
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.SLocalize("StageLines"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("files.fetch"),
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.SLocalize("fetch"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.executeCustomCommand"),
			Handler:     gui.handleCustomCommand,
			Description: gui.Tr.SLocalize("executeCustomCommand"),
		},
		{
			ViewName:    "files",
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateResetToUpstreamMenu,
			Description: gui.Tr.SLocalize("viewResetToUpstreamOptions"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleBranchPress,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.createPullRequest"),
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.SLocalize("createPullRequest"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.checkoutBranchByName"),
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.SLocalize("checkoutByName"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.forceCheckoutBranch"),
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.SLocalize("forceCheckout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("universal.new"),
			Handler:     gui.handleNewBranch,
			Description: gui.Tr.SLocalize("newBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.rebaseBranch"),
			Handler:     gui.handleRebaseOntoLocalBranch,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.mergeIntoCurrentBranch"),
			Handler:     gui.handleMerge,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.viewGitFlowOptions"),
			Handler:     gui.handleCreateGitFlowMenu,
			Description: gui.Tr.SLocalize("gitFlowOptions"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.FastForward"),
			Handler:     gui.handleFastForward,
			Description: gui.Tr.SLocalize("FastForward"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateResetToBranchMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("branches.renameBranch"),
			Handler:     gui.handleRenameBranch,
			Description: gui.Tr.SLocalize("renameBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"local-branches"},
			Key:         gui.getKey("universal.copyToClipboard"),
			Handler:     gui.handleClipboardCopyBranch,
			Description: gui.Tr.SLocalize("copyBranchNameToClipboard"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleCheckoutTag,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleDeleteTag,
			Description: gui.Tr.SLocalize("deleteTag"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gui.getKey("branches.pushTag"),
			Handler:     gui.handlePushTag,
			Description: gui.Tr.SLocalize("pushTag"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gui.getKey("universal.new"),
			Handler:     gui.handleCreateTag,
			Description: gui.Tr.SLocalize("createTag"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"tags"},
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateResetToTagMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "branches",
			Key:         gui.getKey("universal.nextTab"),
			Handler:     gui.handleNextBranchesTab,
			Description: gui.Tr.SLocalize("nextTab"),
		},
		{
			ViewName:    "branches",
			Key:         gui.getKey("universal.prevTab"),
			Handler:     gui.handlePrevBranchesTab,
			Description: gui.Tr.SLocalize("prevTab"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleRemoteBranchesEscape,
			Description: gui.Tr.SLocalize("ReturnToRemotesList"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateResetToRemoteBranchMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         gui.getKey("branches.fetchRemote"),
			Handler:     gui.handleFetchRemote,
			Description: gui.Tr.SLocalize("fetchRemote"),
		},
		{
			ViewName:    "commits",
			Key:         gui.getKey("universal.nextTab"),
			Handler:     gui.handleNextCommitsTab,
			Description: gui.Tr.SLocalize("nextTab"),
		},
		{
			ViewName:    "commits",
			Key:         gui.getKey("universal.prevTab"),
			Handler:     gui.handlePrevCommitsTab,
			Description: gui.Tr.SLocalize("prevTab"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.squashDown"),
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.SLocalize("squashDown"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.renameCommit"),
			Handler:     gui.handleRenameCommit,
			Description: gui.Tr.SLocalize("renameCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.renameCommitWithEditor"),
			Handler:     gui.handleRenameCommitEditor,
			Description: gui.Tr.SLocalize("renameCommitEditor"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateCommitResetMenu,
			Description: gui.Tr.SLocalize("resetToThisCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.markCommitAsFixup"),
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.SLocalize("fixupCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.createFixupCommit"),
			Handler:     gui.handleCreateFixupCommit,
			Description: gui.Tr.SLocalize("createFixupCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.squashAboveCommits"),
			Handler:     gui.handleSquashAllAboveFixupCommits,
			Description: gui.Tr.SLocalize("squashAboveCommits"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.SLocalize("deleteCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.moveDownCommit"),
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.SLocalize("moveDownCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.moveUpCommit"),
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.SLocalize("moveUpCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("universal.edit"),
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.SLocalize("editCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.amendToCommit"),
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.SLocalize("amendToCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.pickCommit"),
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.SLocalize("pickCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.revertCommit"),
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.SLocalize("revertCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.cherryPickCopy"),
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.SLocalize("cherryPickCopy"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("universal.copyToClipboard"),
			Handler:     gui.handleClipboardCopyCommit,
			Description: gui.Tr.SLocalize("copyCommitShaToClipboard"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.cherryPickCopyRange"),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.SLocalize("cherryPickCopyRange"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.pasteCommits"),
			Handler:     gui.HandlePasteCommits,
			Description: gui.Tr.SLocalize("pasteCommits"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("universal.goInto"),
			Handler:     gui.handleSwitchToCommitFilesPanel,
			Description: gui.Tr.SLocalize("viewCommitFiles"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.checkoutCommit"),
			Handler:     gui.handleCheckoutCommit,
			Description: gui.Tr.SLocalize("checkoutCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.tagCommit"),
			Handler:     gui.handleTagCommit,
			Description: gui.Tr.SLocalize("tagCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"branch-commits"},
			Key:         gui.getKey("commits.resetCherryPick"),
			Handler:     gui.handleResetCherryPick,
			Description: gui.Tr.SLocalize("resetCherryPick"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"reflog-commits"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleCheckoutReflogCommit,
			Description: gui.Tr.SLocalize("checkoutCommit"),
		},
		{
			ViewName:    "commits",
			Contexts:    []string{"reflog-commits"},
			Key:         gui.getKey("commits.viewResetOptions"),
			Handler:     gui.handleCreateReflogResetMenu,
			Description: gui.Tr.SLocalize("viewResetOptions"),
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleStashApply,
			Description: gui.Tr.SLocalize("apply"),
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey("stash.popStash"),
			Handler:     gui.handleStashPop,
			Description: gui.Tr.SLocalize("pop"),
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey("universal.remove"),
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
			ViewName:    "menu",
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleMenuClose,
			Description: gui.Tr.SLocalize("closeMenu"),
		},
		{
			ViewName:    "menu",
			Key:         gui.getKey("universal.quit"),
			Handler:     gui.handleMenuClose,
			Description: gui.Tr.SLocalize("closeMenu"),
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleSwitchToCommitsPanel,
			Description: gui.Tr.SLocalize("goBack"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("commitFiles.checkoutCommitFile"),
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.Tr.SLocalize("checkoutCommitFile"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.Tr.SLocalize("discardOldFileChange"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("universal.openFile"),
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.Tr.SLocalize("toggleAddToPatch"),
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey("universal.goInto"),
			Handler:     gui.handleEnterCommitFile,
			Description: gui.Tr.SLocalize("enterFile"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.filteringMenu"),
			Handler:     gui.handleCreateFilteringMenuPanel,
			Description: gui.Tr.SLocalize("openScopingMenu"),
		},
		{
			ViewName:    "",
			Key:         gui.getKey("universal.diffingMenu"),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.Tr.SLocalize("openDiffingMenu"),
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
			Handler:     gui.scrollDownMain,
			Description: gui.Tr.SLocalize("ScrollDown"),
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{"normal"},
			Key:         gocui.MouseWheelUp,
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
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleStagingEscape,
			Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleToggleStagedSelection,
			Description: gui.Tr.SLocalize("StageSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleResetSelection,
			Description: gui.Tr.SLocalize("ResetSelection"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("universal.togglePanel"),
			Handler:     gui.handleTogglePanel,
			Description: gui.Tr.SLocalize("TogglePanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building"},
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleEscapePatchBuildingPanel,
			Description: gui.Tr.SLocalize("ExitLineByLineMode"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("universal.prevItem"),
			Handler:     gui.handleSelectPrevLine,
			Description: gui.Tr.SLocalize("PrevLine"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("universal.nextItem"),
			Handler:     gui.handleSelectNextLine,
			Description: gui.Tr.SLocalize("NextLine"),
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gui.getKey("universal.prevItem-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gui.getKey("universal.nextItem-alt"),
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
			Key:         gui.getKey("universal.prevBlock"),
			Handler:     gui.handleSelectPrevHunk,
			Description: gui.Tr.SLocalize("PrevHunk"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("universal.nextBlock"),
			Handler:     gui.handleSelectNextHunk,
			Description: gui.Tr.SLocalize("NextHunk"),
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gui.getKey("universal.prevBlock-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{"patch-building", "staging"},
			Key:      gui.getKey("universal.nextBlock-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("universal.edit"),
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.SLocalize("editFile"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("universal.openFile"),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.SLocalize("openFile"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleToggleSelectionForPatch,
			Description: gui.Tr.SLocalize("ToggleSelectionForPatch"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("main.toggleDragSelect"),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.Tr.SLocalize("ToggleDragSelect"),
		},
		// Alias 'V' -> 'v'
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("main.toggleDragSelect-alt"),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.Tr.SLocalize("ToggleDragSelect"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"patch-building", "staging"},
			Key:         gui.getKey("main.toggleSelectHunk"),
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
			Contexts:    []string{"staging"},
			Key:         gui.getKey("files.commitChanges"),
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.SLocalize("CommitChanges"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("files.commitChangesWithoutHook"),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.SLocalize("commitChangesWithoutHook"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"staging"},
			Key:         gui.getKey("files.commitChangesWithEditor"),
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.SLocalize("CommitChangesWithEditor"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.return"),
			Handler:     gui.handleEscapeMerge,
			Description: gui.Tr.SLocalize("ReturnToFilesPanel"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handlePickHunk,
			Description: gui.Tr.SLocalize("PickHunk"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("main.pickBothHunks"),
			Handler:     gui.handlePickBothHunks,
			Description: gui.Tr.SLocalize("PickBothHunks"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.prevBlock"),
			Handler:     gui.handleSelectPrevConflict,
			Description: gui.Tr.SLocalize("PrevConflict"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.nextBlock"),
			Handler:     gui.handleSelectNextConflict,
			Description: gui.Tr.SLocalize("NextConflict"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.prevItem"),
			Handler:     gui.handleSelectTop,
			Description: gui.Tr.SLocalize("SelectTop"),
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.nextItem"),
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
			Key:      gui.getKey("universal.prevBlock-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      gui.getKey("universal.nextBlock-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      gui.getKey("universal.prevItem-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		},
		{
			ViewName: "main",
			Contexts: []string{"mergin"},
			Key:      gui.getKey("universal.nextItem-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		},
		{
			ViewName:    "main",
			Contexts:    []string{"merging"},
			Key:         gui.getKey("universal.undo"),
			Handler:     gui.handlePopFileSnapshot,
			Description: gui.Tr.SLocalize("undo"),
		},
		{
			ViewName: "branches",
			Contexts: []string{"remotes"},
			Key:      gui.getKey("universal.goInto"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleRemoteEnter,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         gui.getKey("universal.new"),
			Handler:     gui.handleAddRemote,
			Description: gui.Tr.SLocalize("addNewRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleRemoveRemote,
			Description: gui.Tr.SLocalize("removeRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remotes"},
			Key:         gui.getKey("universal.edit"),
			Handler:     gui.handleEditRemote,
			Description: gui.Tr.SLocalize("editRemote"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("universal.select"),
			Handler:     gui.handleCheckoutRemoteBranch,
			Description: gui.Tr.SLocalize("checkout"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("universal.new"),
			Handler:     gui.handleNewBranchOffRemote,
			Description: gui.Tr.SLocalize("newBranch"),
		},

		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("branches.mergeIntoCurrentBranch"),
			Handler:     gui.handleMergeRemoteBranch,
			Description: gui.Tr.SLocalize("mergeIntoCurrentBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("universal.remove"),
			Handler:     gui.handleDeleteRemoteBranch,
			Description: gui.Tr.SLocalize("deleteBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("branches.rebaseBranch"),
			Handler:     gui.handleRebaseOntoRemoteBranch,
			Description: gui.Tr.SLocalize("rebaseBranch"),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{"remote-branches"},
			Key:         gui.getKey("branches.setUpstream"),
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
		{
			ViewName: "search",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSearch,
		},
		{
			ViewName: "search",
			Key:      gui.getKey("universal.return"),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSearchEscape,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey("universal.prevItem"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey("universal.nextItem"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey("universal.prevItem-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey("universal.nextItem-alt"),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: gui.getKey("universal.togglePanel"), Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gui.getKey("universal.prevBlock"), Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gui.getKey("universal.nextBlock"), Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gui.getKey("universal.prevBlock-alt"), Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gui.getKey("universal.nextBlock-alt"), Modifier: gocui.ModNone, Handler: gui.nextView},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	for i, viewName := range []string{"status", "files", "branches", "commits", "stash"} {
		bindings = append(bindings, &Binding{ViewName: "", Key: rune(i+1) + '0', Modifier: gocui.ModNone, Handler: gui.goToSideView(viewName)})
	}

	for _, listView := range gui.getListViews() {
		bindings = append(bindings, []*Binding{
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.prevItem-alt"), Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.prevItem"), Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: listView.handlePrevLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.nextItem-alt"), Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.nextItem"), Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.prevPage"), Modifier: gocui.ModNone, Handler: listView.handlePrevPage, Description: gui.Tr.SLocalize("prevPage")},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.nextPage"), Modifier: gocui.ModNone, Handler: listView.handleNextPage, Description: gui.Tr.SLocalize("nextPage")},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gui.getKey("universal.gotoTop"), Modifier: gocui.ModNone, Handler: listView.handleGotoTop, Description: gui.Tr.SLocalize("gotoTop")},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: listView.handleNextLine},
			{ViewName: listView.viewName, Contexts: []string{listView.context}, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: listView.handleClick},
		}...)

		// the commits panel needs to lazyload things so it has a couple of its own handlers
		openSearchHandler := gui.handleOpenSearch
		gotoBottomHandler := listView.handleGotoBottom
		if listView.viewName == "commits" {
			openSearchHandler = gui.handleOpenSearchForCommitsPanel
			gotoBottomHandler = gui.handleGotoBottomForCommitsPanel
		}

		bindings = append(bindings, []*Binding{
			{
				ViewName:    listView.viewName,
				Contexts:    []string{listView.context},
				Key:         gui.getKey("universal.startSearch"),
				Handler:     openSearchHandler,
				Description: gui.Tr.SLocalize("startSearch"),
			},
			{
				ViewName:    listView.viewName,
				Contexts:    []string{listView.context},
				Key:         gui.getKey("universal.gotoBottom"),
				Handler:     gotoBottomHandler,
				Description: gui.Tr.SLocalize("gotoBottom"),
			},
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

	tabClickBindings := map[string]func(int) error{
		"branches": gui.onBranchesTabClick,
		"commits":  gui.onCommitsTabClick,
	}

	for viewName, binding := range tabClickBindings {
		if err := g.SetTabClickBinding(viewName, binding); err != nil {
			return err
		}
	}

	return nil
}
