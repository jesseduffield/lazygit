package gui

import (
	"fmt"
	"log"
	"strings"

	"unicode/utf8"

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
	Tag         string // e.g. 'navigation'. Used for grouping things in the cheatsheet
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

	return fmt.Sprintf("%c", keyInt)
}

func (gui *Gui) getKey(key string) interface{} {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding := keymap[strings.ToLower(key)]
		if binding == nil {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md", strings.ToLower(key))
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return []rune(key)[0]
	}
	log.Fatal("Key empty for keybinding: " + strings.ToLower(key))
	return nil
}

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*Binding {
	config := gui.Config.GetUserConfig().Keybinding

	bindings := []*Binding{
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.Quit),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleQuit),
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.QuitWithoutChangingDirectory),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuitWithoutChangingDirectory,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.QuitAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleQuit),
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  gui.handleTopLevelReturn,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ScrollUpMain),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
			Description: gui.Tr.LcScrollUpMainPanel,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ScrollDownMain),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
			Description: gui.Tr.LcScrollDownMainPanel,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.ScrollUpMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.ScrollDownMainAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.ScrollUpMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.ScrollDownMainAlt2),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.CreateRebaseOptionsMenu),
			Handler:     gui.wrappedHandler(gui.handleCreateRebaseOptionsMenu),
			Description: gui.Tr.ViewMergeRebaseOptions,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.CreatePatchOptionsMenu),
			Handler:     gui.handleCreatePatchOptionsMenu,
			Description: gui.Tr.ViewPatchOptions,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.PushFiles),
			Handler:     gui.pushFiles,
			Description: gui.Tr.LcPush,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.PullFiles),
			Handler:     gui.handlePullFiles,
			Description: gui.Tr.LcPull,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.Refresh),
			Handler:     gui.handleRefresh,
			Description: gui.Tr.LcRefresh,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.OptionMenu),
			Handler:     gui.handleCreateOptionsMenu,
			Description: gui.Tr.LcOpenMenu,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.OptionMenuAlt1),
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
			Key:         gui.getKey(config.Universal.Undo),
			Handler:     gui.reflogUndo,
			Description: gui.Tr.LcUndoReflog,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.Redo),
			Handler:     gui.reflogRedo,
			Description: gui.Tr.LcRedoReflog,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.EditConfig,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.NextScreenMode),
			Handler:     gui.nextScreenMode,
			Description: gui.Tr.LcNextScreenMode,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.PrevScreenMode),
			Handler:     gui.prevScreenMode,
			Description: gui.Tr.LcPrevScreenMode,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.OpenConfig,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.CheckForUpdate),
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.LcCheckForUpdate,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.RecentRepos),
			Handler:     gui.wrappedHandler(gui.handleCreateRecentReposMenu),
			Description: gui.Tr.SwitchRepo,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.AllBranchesLogGraph),
			Handler:     gui.wrappedHandler(gui.handleShowAllBranchLogs),
			Description: gui.Tr.LcAllBranchesLogGraph,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChanges),
			Handler:     gui.wrappedHandler(gui.handleCommitPress),
			Description: gui.Tr.CommitChanges,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChangesWithoutHook),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.LcCommitChangesWithoutHook,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.AmendLastCommit),
			Handler:     gui.wrappedHandler(gui.handleAmendCommitPress),
			Description: gui.Tr.AmendLastCommit,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChangesWithEditor),
			Handler:     gui.wrappedHandler(gui.handleCommitEditorPress),
			Description: gui.Tr.CommitChangesWithEditor,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.wrappedHandler(gui.handleFilePress),
			Description: gui.Tr.LcToggleStaged,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.Tr.LcViewDiscardOptions,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.IgnoreFile),
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.LcIgnoreFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.RefreshFiles),
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.LcRefreshFiles,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.StashAllChanges),
			Handler:     gui.handleStashChanges,
			Description: gui.Tr.LcStashAllChanges,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.ViewStashOptions),
			Handler:     gui.handleCreateStashMenu,
			Description: gui.Tr.LcViewStashOptions,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.ToggleStagedAll),
			Handler:     gui.handleStageAll,
			Description: gui.Tr.LcToggleStagedAll,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.ViewResetOptions),
			Handler:     gui.handleCreateResetMenu,
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.StageLines,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.Fetch),
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.LcFetch,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyFileNameToClipboard,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ExecuteCustomCommand),
			Handler:     gui.handleCustomCommand,
			Description: gui.Tr.LcExecuteCustomCommand,
		},
		{
			ViewName:    "files",
			Contexts:    []string{FILES_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToUpstreamMenu,
			Description: gui.Tr.LcViewResetToUpstreamOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleBranchPress,
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.CreatePullRequest),
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.LcCreatePullRequest,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.CopyPullRequestURL),
			Handler:     gui.handleCopyPullRequestURLPress,
			Description: gui.Tr.LcCopyPullRequestURL,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.CheckoutBranchByName),
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.LcCheckoutByName,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.ForceCheckoutBranch),
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.LcForceCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.RebaseBranch),
			Handler:     gui.handleRebaseOntoLocalBranch,
			Description: gui.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     gui.handleMerge,
			Description: gui.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.ViewGitFlowOptions),
			Handler:     gui.handleCreateGitFlowMenu,
			Description: gui.Tr.LcGitFlowOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.FastForward),
			Handler:     gui.handleFastForward,
			Description: gui.Tr.FastForward,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToBranchMenu,
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.RenameBranch),
			Handler:     gui.handleRenameBranch,
			Description: gui.Tr.LcRenameBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyBranchNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{LOCAL_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleSwitchToSubCommits),
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutTag,
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteTag,
			Description: gui.Tr.LcDeleteTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.PushTag),
			Handler:     gui.handlePushTag,
			Description: gui.Tr.LcPushTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleCreateTag,
			Description: gui.Tr.LcCreateTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToTagMenu,
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{TAGS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleSwitchToSubCommits),
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleRemoteBranchesEscape,
			Description: gui.Tr.ReturnToRemotesList,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToRemoteBranchMenu,
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleSwitchToSubCommits),
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.FetchRemote),
			Handler:     gui.handleFetchRemote,
			Description: gui.Tr.LcFetchRemote,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.SquashDown),
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.LcSquashDown,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.RenameCommit),
			Handler:     gui.handleRenameCommit,
			Description: gui.Tr.LcRenameCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.RenameCommitWithEditor),
			Handler:     gui.handleRenameCommitEditor,
			Description: gui.Tr.LcRenameCommitEditor,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateCommitResetMenu,
			Description: gui.Tr.LcResetToThisCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.MarkCommitAsFixup),
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.LcFixupCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CreateFixupCommit),
			Handler:     gui.handleCreateFixupCommit,
			Description: gui.Tr.LcCreateFixupCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.SquashAboveCommits),
			Handler:     gui.handleSquashAllAboveFixupCommits,
			Description: gui.Tr.LcSquashAboveCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.LcDeleteCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.MoveDownCommit),
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.LcMoveDownCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.MoveUpCommit),
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.LcMoveUpCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.LcEditCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.AmendToCommit),
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.LcAmendToCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.PickCommit),
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.LcPickCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.RevertCommit),
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.LcRevertCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.wrappedHandler(gui.handleCopyCommit),
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.wrappedHandler(gui.handleCopyCommitRange),
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.PasteCommits),
			Handler:     gui.wrappedHandler(gui.HandlePasteCommits),
			Description: gui.Tr.LcPasteCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleViewCommitFiles),
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CheckoutCommit),
			Handler:     gui.handleCheckoutCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcCreateNewBranchFromCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.TagCommit),
			Handler:     gui.handleTagCommit,
			Description: gui.Tr.LcTagCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.wrappedHandler(gui.exitCherryPickingMode),
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{BRANCH_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CopyCommitMessageToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedCommitMessageToClipboard),
			Description: gui.Tr.LcCopyCommitMessageToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleViewReflogCommitFiles),
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutReflogCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateReflogResetMenu,
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.wrappedHandler(gui.handleCopyCommit),
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.wrappedHandler(gui.handleCopyCommitRange),
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.wrappedHandler(gui.exitCherryPickingMode),
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{REFLOG_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyCommitFileNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleViewSubCommitFiles),
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutSubCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.wrappedHandler(gui.handleCreateSubCommitResetMenu),
			Description: gui.Tr.LcViewResetOptions,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.wrappedHandler(gui.handleCopyCommit),
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.wrappedHandler(gui.handleCopyCommitRange),
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.wrappedHandler(gui.exitCherryPickingMode),
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{SUB_COMMITS_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.wrappedHandler(gui.handleViewStashFiles),
			Description: gui.Tr.LcViewStashFiles,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleStashApply,
			Description: gui.Tr.LcApply,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Stash.PopStash),
			Handler:     gui.handleStashPop,
			Description: gui.Tr.LcPop,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleStashDrop,
			Description: gui.Tr.LcDrop,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName: "commitMessage",
			Key:      gui.getKey(config.Universal.SubmitEditorText),
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitConfirm,
		},
		{
			ViewName: "commitMessage",
			Key:      gui.getKey(config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  gui.handleCommitClose,
		},
		{
			ViewName: "credentials",
			Key:      gui.getKey(config.Universal.Confirm),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSubmitCredential,
		},
		{
			ViewName: "credentials",
			Key:      gui.getKey(config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  gui.handleCloseCredentialsView,
		},
		{
			ViewName:    "menu",
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleMenuClose,
			Description: gui.Tr.LcCloseMenu,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.CommitFiles.CheckoutCommitFile),
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.Tr.LcCheckoutCommitFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.Tr.LcDiscardOldFileChange,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditCommitFile,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.Tr.LcToggleAddToPatch,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleEnterCommitFile,
			Description: gui.Tr.LcEnterFile,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.FilteringMenu),
			Handler:     gui.handleCreateFilteringMenuPanel,
			Description: gui.Tr.LcOpenScopingMenu,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.DiffingMenu),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.Tr.LcOpenDiffingMenu,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.DiffingMenuAlt),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.Tr.LcOpenDiffingMenu,
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
			Contexts: []string{MAIN_NORMAL_CONTEXT_KEY},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownSecondary,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_NORMAL_CONTEXT_KEY},
			Key:         gocui.MouseWheelDown,
			Handler:     gui.scrollDownMain,
			Description: gui.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_NORMAL_CONTEXT_KEY},
			Key:         gocui.MouseWheelUp,
			Handler:     gui.scrollUpMain,
			Description: gui.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_NORMAL_CONTEXT_KEY},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownMain,
		},
		{
			ViewName: "secondary",
			Contexts: []string{MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleTogglePanelClick,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.wrappedHandler(gui.handleStagingEscape),
			Description: gui.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.wrappedHandler(gui.handleToggleStagedSelection),
			Description: gui.Tr.StageSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.wrappedHandler(gui.handleResetSelection),
			Description: gui.Tr.ResetSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.TogglePanel),
			Handler:     gui.wrappedHandler(gui.handleTogglePanel),
			Description: gui.Tr.TogglePanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.wrappedHandler(gui.handleEscapePatchBuildingPanel),
			Description: gui.Tr.ExitLineByLineMode,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.wrappedHandler(gui.handleOpenFileAtLine),
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.PrevItem),
			Handler:     gui.wrappedHandler(gui.handleSelectPrevLine),
			Description: gui.Tr.PrevLine,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.NextItem),
			Handler:     gui.wrappedHandler(gui.handleSelectNextLine),
			Description: gui.Tr.NextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectPrevLine),
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectNextLine),
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectPrevLine),
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectNextLine),
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.wrappedHandler(gui.handleSelectPrevHunk),
			Description: gui.Tr.PrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectPrevHunk),
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.NextBlock),
			Handler:     gui.wrappedHandler(gui.handleSelectNextHunk),
			Description: gui.Tr.NextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleSelectNextHunk),
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.NextPage),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleLineByLineNextPage),
			Description: gui.Tr.LcNextPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.PrevPage),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleLineByLinePrevPage),
			Description: gui.Tr.LcPrevPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GotoTop),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleLineByLineGotoTop),
			Description: gui.Tr.LcGotoTop,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GotoBottom),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleLineByLineGotoBottom),
			Description: gui.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.StartSearch),
			Handler:     gui.handleOpenSearch,
			Description: gui.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.wrappedHandler(gui.handleToggleSelectionForPatch),
			Description: gui.Tr.ToggleSelectionForPatch,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Main.ToggleDragSelect),
			Handler:     gui.wrappedHandler(gui.handleToggleSelectRange),
			Description: gui.Tr.ToggleDragSelect,
		},
		// Alias 'V' -> 'v'
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Main.ToggleDragSelectAlt),
			Handler:     gui.wrappedHandler(gui.handleToggleSelectRange),
			Description: gui.Tr.ToggleDragSelect,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Main.ToggleSelectHunk),
			Handler:     gui.wrappedHandler(gui.handleToggleSelectHunk),
			Description: gui.Tr.ToggleSelectHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleLBLMouseDown,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModMotion,
			Handler:  gui.handleMouseDrag,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseScrollUp,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_PATCH_BUILDING_CONTEXT_KEY, MAIN_STAGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseScrollDown,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChanges),
			Handler:     gui.wrappedHandler(gui.handleCommitPress),
			Description: gui.Tr.CommitChanges,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChangesWithoutHook),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.LcCommitChangesWithoutHook,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_STAGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Files.CommitChangesWithEditor),
			Handler:     gui.wrappedHandler(gui.handleCommitEditorPress),
			Description: gui.Tr.CommitChangesWithEditor,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.wrappedHandler(gui.handleEscapeMerge),
			Description: gui.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handlePickHunk,
			Description: gui.Tr.PickHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Main.PickBothHunks),
			Handler:     gui.handlePickBothHunks,
			Description: gui.Tr.PickBothHunks,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.handleSelectPrevConflict,
			Description: gui.Tr.PrevConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.NextBlock),
			Handler:     gui.handleSelectNextConflict,
			Description: gui.Tr.NextConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.PrevItem),
			Handler:     gui.handleSelectTop,
			Description: gui.Tr.SelectTop,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.NextItem),
			Handler:     gui.handleSelectBottom,
			Description: gui.Tr.SelectBottom,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectTop,
		},
		{
			ViewName: "main",
			Contexts: []string{MAIN_MERGING_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectBottom,
		},
		{
			ViewName:    "main",
			Contexts:    []string{MAIN_MERGING_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Undo),
			Handler:     gui.handlePopFileSnapshot,
			Description: gui.Tr.LcUndo,
		},
		{
			ViewName: "branches",
			Contexts: []string{REMOTES_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.GoInto),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.handleRemoteEnter),
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleAddRemote,
			Description: gui.Tr.LcAddNewRemote,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleRemoveRemote,
			Description: gui.Tr.LcRemoveRemote,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditRemote,
			Description: gui.Tr.LcEditRemote,
		},
		{
			ViewName: "branches",
			Contexts: []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:      gui.getKey(config.Universal.Select),
			// gonna use the exact same handler as the 'n' keybinding because everybody wants this to happen when they checkout a remote branch
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.wrappedHandler(gui.handleNewBranchOffCurrentItem),
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     gui.handleMergeRemoteBranch,
			Description: gui.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteRemoteBranch,
			Description: gui.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.RebaseBranch),
			Handler:     gui.handleRebaseOntoRemoteBranch,
			Description: gui.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{REMOTE_BRANCHES_CONTEXT_KEY},
			Key:         gui.getKey(config.Branches.SetUpstream),
			Handler:     gui.handleSetBranchUpstream,
			Description: gui.Tr.LcSetUpstream,
		},
		{
			ViewName: "status",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleStatusClick,
		},
		{
			ViewName: "search",
			Key:      gui.getKey(config.Universal.Confirm),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSearch,
		},
		{
			ViewName: "search",
			Key:      gui.getKey(config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSearchEscape,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey(config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey(config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "menu",
			Key:      gui.getKey(config.Universal.Select),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.onMenuPress),
		},
		{
			ViewName: "menu",
			Key:      gui.getKey(config.Universal.Confirm),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.onMenuPress),
		},
		{
			ViewName: "menu",
			Key:      gui.getKey(config.Universal.ConfirmAlt1),
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedHandler(gui.onMenuPress),
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.wrappedHandler(gui.handleCopySelectedSideContextItemToClipboard),
			Description: gui.Tr.LcCopySubmoduleNameToClipboard,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.forSubmodule(gui.handleSubmoduleEnter),
			Description: gui.Tr.LcEnterSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.forSubmodule(gui.handleResetRemoveSubmodule),
			Description: gui.Tr.LcViewResetAndRemoveOptions,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Submodules.Update),
			Handler:     gui.forSubmodule(gui.handleUpdateSubmodule),
			Description: gui.Tr.LcSubmoduleUpdate,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.wrappedHandler(gui.handleAddSubmodule),
			Description: gui.Tr.LcAddSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.forSubmodule(gui.handleEditSubmoduleUrl),
			Description: gui.Tr.LcEditSubmoduleUrl,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Submodules.Init),
			Handler:     gui.forSubmodule(gui.handleSubmoduleInit),
			Description: gui.Tr.LcInitSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{SUBMODULES_CONTEXT_KEY},
			Key:         gui.getKey(config.Submodules.BulkMenu),
			Handler:     gui.wrappedHandler(gui.handleBulkSubmoduleActionsMenu),
			Description: gui.Tr.LcViewBulkSubmoduleOptions,
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: gui.getKey(config.Universal.PrevBlock), Modifier: gocui.ModNone, Handler: gui.wrappedHandler(gui.previousSideWindow)},
			{ViewName: viewName, Key: gui.getKey(config.Universal.NextBlock), Modifier: gocui.ModNone, Handler: gui.wrappedHandler(gui.nextSideWindow)},
			{ViewName: viewName, Key: gui.getKey(config.Universal.PrevBlockAlt), Modifier: gocui.ModNone, Handler: gui.wrappedHandler(gui.previousSideWindow)},
			{ViewName: viewName, Key: gui.getKey(config.Universal.NextBlockAlt), Modifier: gocui.ModNone, Handler: gui.wrappedHandler(gui.nextSideWindow)},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	for i, window := range []string{"status", "files", "branches", "commits", "stash"} {
		bindings = append(bindings, &Binding{ViewName: "", Key: rune(i+1) + '0', Modifier: gocui.ModNone, Handler: gui.goToSideWindow(window)})
	}

	for viewName := range gui.viewTabContextMap() {
		bindings = append(bindings, []*Binding{
			{
				ViewName:    viewName,
				Key:         gui.getKey(config.Universal.NextTab),
				Handler:     gui.handleNextTab,
				Description: gui.Tr.LcNextTab,
				Tag:         "navigation",
			},
			{
				ViewName:    viewName,
				Key:         gui.getKey(config.Universal.PrevTab),
				Handler:     gui.handlePrevTab,
				Description: gui.Tr.LcPrevTab,
				Tag:         "navigation",
			},
		}...)
	}

	bindings = append(bindings, gui.getListContextKeyBindings()...)

	return bindings
}

func (gui *Gui) keybindings() error {
	bindings := gui.GetCustomCommandKeybindings()

	bindings = append(bindings, gui.GetInitialKeybindings()...)

	for _, binding := range bindings {
		if err := gui.g.SetKeybinding(binding.ViewName, binding.Contexts, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	for viewName := range gui.viewTabContextMap() {
		viewName := viewName
		tabClickCallback := func(tabIndex int) error { return gui.onViewTabClick(viewName, tabIndex) }

		if err := gui.g.SetTabClickBinding(viewName, tabClickCallback); err != nil {
			return err
		}
	}

	return nil
}
