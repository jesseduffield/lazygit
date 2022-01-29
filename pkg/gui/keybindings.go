package gui

import (
	"fmt"
	"log"
	"strings"

	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

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
	gocui.KeyTab:        "tab", // ctrl+i
	gocui.KeyBacktab:    "shift+tab",
	gocui.KeyEnter:      "enter", // ctrl+m
	gocui.KeyAltEnter:   "alt+enter",
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
	"<backtab>":   gocui.KeyBacktab,
	"<enter>":     gocui.KeyEnter,
	"<a-enter>":   gocui.KeyAltEnter,
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

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*types.Binding {
	config := gui.c.UserConfig.Keybinding

	guards := types.KeybindingGuards{
		OutsideFilterMode: gui.outsideFilterMode,
		NoPopupPanel:      gui.noPopupPanel,
	}

	bindings := []*types.Binding{
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.Quit),
			Modifier: gocui.ModNone,
			Handler:  gui.handleQuit,
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
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			Key:      gui.getKey(config.Universal.Return),
			Modifier: gocui.ModNone,
			Handler:  gui.handleTopLevelReturn,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.OpenRecentRepos),
			Handler:     gui.handleCreateRecentReposMenu,
			Alternative: "<c-r>",
			Description: gui.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ScrollUpMain),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
			Description: gui.c.Tr.LcScrollUpMainPanel,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ScrollDownMain),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
			Description: gui.c.Tr.LcScrollDownMainPanel,
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
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.c.Tr.ViewMergeRebaseOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.CreatePatchOptionsMenu),
			Handler:     gui.handleCreatePatchOptionsMenu,
			Description: gui.c.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.Refresh),
			Handler:     gui.handleRefresh,
			Description: gui.c.Tr.LcRefresh,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.OptionMenu),
			Handler:     gui.handleCreateOptionsMenu,
			Description: gui.c.Tr.LcOpenMenu,
			OpensMenu:   true,
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
			ViewName:    "status",
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditConfig,
			Description: gui.c.Tr.EditConfig,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.NextScreenMode),
			Handler:     gui.nextScreenMode,
			Description: gui.c.Tr.LcNextScreenMode,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.PrevScreenMode),
			Handler:     gui.prevScreenMode,
			Description: gui.c.Tr.LcPrevScreenMode,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenConfig,
			Description: gui.c.Tr.OpenConfig,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.CheckForUpdate),
			Handler:     gui.handleCheckForUpdate,
			Description: gui.c.Tr.LcCheckForUpdate,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.RecentRepos),
			Handler:     gui.handleCreateRecentReposMenu,
			Description: gui.c.Tr.SwitchRepo,
		},
		{
			ViewName:    "status",
			Key:         gui.getKey(config.Status.AllBranchesLogGraph),
			Handler:     gui.handleShowAllBranchLogs,
			Description: gui.c.Tr.LcAllBranchesLogGraph,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.c.Tr.LcViewDiscardOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.ViewResetOptions),
			Handler:     gui.handleCreateResetMenu,
			Description: gui.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.Fetch),
			Handler:     gui.handleGitFetch,
			Description: gui.c.Tr.LcFetch,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyFileNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleBranchPress,
			Description: gui.c.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.CreatePullRequest),
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.c.Tr.LcCreatePullRequest,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.ViewPullRequestOptions),
			Handler:     gui.handleCreatePullRequestMenu,
			Description: gui.c.Tr.LcCreatePullRequestOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.CopyPullRequestURL),
			Handler:     gui.handleCopyPullRequestURLPress,
			Description: gui.c.Tr.LcCopyPullRequestURL,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.CheckoutBranchByName),
			Handler:     gui.handleCheckoutByName,
			Description: gui.c.Tr.LcCheckoutByName,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.ForceCheckoutBranch),
			Handler:     gui.handleForceCheckout,
			Description: gui.c.Tr.LcForceCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteBranch,
			Description: gui.c.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.RebaseBranch),
			Handler:     guards.OutsideFilterMode(gui.handleRebaseOntoLocalBranch),
			Description: gui.c.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     guards.OutsideFilterMode(gui.handleMerge),
			Description: gui.c.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.ViewGitFlowOptions),
			Handler:     gui.handleCreateGitFlowMenu,
			Description: gui.c.Tr.LcGitFlowOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.FastForward),
			Handler:     gui.handleFastForward,
			Description: gui.c.Tr.FastForward,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToBranchMenu,
			Description: gui.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.RenameBranch),
			Handler:     gui.handleRenameBranch,
			Description: gui.c.Tr.LcRenameBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyBranchNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.LOCAL_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleSwitchToSubCommits,
			Description: gui.c.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleRemoteBranchesEscape,
			Description: gui.c.Tr.ReturnToRemotesList,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToRemoteBranchMenu,
			Description: gui.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleSwitchToSubCommits,
			Description: gui.c.Tr.LcViewCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.handleCopyCommit,
			Description: gui.c.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.c.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.PasteCommits),
			Handler:     guards.OutsideFilterMode(gui.HandlePasteCommits),
			Description: gui.c.Tr.LcPasteCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.New),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcCreateNewBranchFromCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.BRANCH_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.c.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewReflogCommitFiles,
			Description: gui.c.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.CheckoutReflogCommit,
			Description: gui.c.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateReflogResetMenu,
			Description: gui.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     guards.OutsideFilterMode(gui.handleCopyCommit),
			Description: gui.c.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     guards.OutsideFilterMode(gui.handleCopyCommitRange),
			Description: gui.c.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.c.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(context.REFLOG_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewSubCommitFiles,
			Description: gui.c.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutSubCommit,
			Description: gui.c.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateSubCommitResetMenu,
			Description: gui.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.handleCopyCommit,
			Description: gui.c.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.c.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.c.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.SUB_COMMITS_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewStashFiles,
			Description: gui.c.Tr.LcViewStashFiles,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleStashApply,
			Description: gui.c.Tr.LcApply,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Stash.PopStash),
			Handler:     gui.handleStashPop,
			Description: gui.c.Tr.LcPop,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleStashDrop,
			Description: gui.c.Tr.LcDrop,
		},
		{
			ViewName:    "stash",
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcNewBranch,
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
			Description: gui.c.Tr.LcCloseMenu,
		},
		{
			ViewName: "information",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopyCommitFileNameToClipboard,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.CommitFiles.CheckoutCommitFile),
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.c.Tr.LcCheckoutCommitFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.c.Tr.LcDiscardOldFileChange,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditCommitFile,
			Description: gui.c.Tr.LcEditFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.c.Tr.LcToggleAddToPatch,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleEnterCommitFile,
			Description: gui.c.Tr.LcEnterFile,
		},
		{
			ViewName:    "commitFiles",
			Key:         gui.getKey(config.Files.ToggleTreeView),
			Handler:     gui.handleToggleCommitFileTreeView,
			Description: gui.c.Tr.LcToggleTreeView,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.FilteringMenu),
			Handler:     gui.handleCreateFilteringMenuPanel,
			Description: gui.c.Tr.LcOpenFilteringMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.DiffingMenu),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.c.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.DiffingMenuAlt),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.c.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.ExtrasMenu),
			Handler:     gui.handleCreateExtrasMenuPanel,
			Description: gui.c.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
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
			Contexts: []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownSecondary,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:         gocui.MouseWheelDown,
			Handler:     gui.scrollDownMain,
			Description: gui.c.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:         gocui.MouseWheelUp,
			Handler:     gui.scrollUpMain,
			Description: gui.c.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_NORMAL_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleMouseDownMain,
		},
		{
			ViewName: "secondary",
			Contexts: []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleTogglePanelClick,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleStagingEscape,
			Description: gui.c.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleStagedSelection,
			Description: gui.c.Tr.StageSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleResetSelection,
			Description: gui.c.Tr.ResetSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.TogglePanel),
			Handler:     gui.handleTogglePanel,
			Description: gui.c.Tr.TogglePanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleEscapePatchBuildingPanel,
			Description: gui.c.Tr.ExitLineByLineMode,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenFileAtLine,
			Description: gui.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.PrevItem),
			Handler:     gui.handleSelectPrevLine,
			Description: gui.c.Tr.PrevLine,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.NextItem),
			Handler:     gui.handleSelectNextLine,
			Description: gui.c.Tr.NextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.handleSelectPrevHunk,
			Description: gui.c.Tr.PrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.NextBlock),
			Handler:     gui.handleSelectNextHunk,
			Description: gui.c.Tr.NextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Modifier:    gocui.ModNone,
			Handler:     gui.copySelectedToClipboard,
			Description: gui.c.Tr.LcCopySelectedTexToClipboard,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Edit),
			Handler:     gui.handleLineByLineEdit,
			Description: gui.c.Tr.LcEditFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.OpenFile),
			Handler:     gui.Controllers.Files.Open,
			Description: gui.c.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.NextPage),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLineByLineNextPage,
			Description: gui.c.Tr.LcNextPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.PrevPage),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLineByLinePrevPage,
			Description: gui.c.Tr.LcPrevPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GotoTop),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLineByLineGotoTop,
			Description: gui.c.Tr.LcGotoTop,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.GotoBottom),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLineByLineGotoBottom,
			Description: gui.c.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.StartSearch),
			Handler:     func() error { return gui.handleOpenSearch("main") },
			Description: gui.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleSelectionForPatch,
			Description: gui.c.Tr.ToggleSelectionForPatch,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Main.ToggleDragSelect),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.c.Tr.ToggleDragSelect,
		},
		// Alias 'V' -> 'v'
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Main.ToggleDragSelectAlt),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.c.Tr.ToggleDragSelect,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Main.ToggleSelectHunk),
			Handler:     gui.handleToggleSelectHunk,
			Description: gui.c.Tr.ToggleSelectHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleLBLMouseDown,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModMotion,
			Handler:  gui.handleMouseDrag,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelUp,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:      gocui.MouseWheelDown,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY), string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.ScrollLeft),
			Handler:     gui.scrollLeftMain,
			Description: gui.c.Tr.LcScrollLeft,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_PATCH_BUILDING_CONTEXT_KEY), string(context.MAIN_STAGING_CONTEXT_KEY), string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.ScrollRight),
			Handler:     gui.scrollRightMain,
			Description: gui.c.Tr.LcScrollRight,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.CommitChanges),
			Handler:     gui.Controllers.Files.HandleCommitPress,
			Description: gui.c.Tr.CommitChanges,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.CommitChangesWithoutHook),
			Handler:     gui.Controllers.Files.HandleWIPCommitPress,
			Description: gui.c.Tr.LcCommitChangesWithoutHook,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_STAGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.CommitChangesWithEditor),
			Handler:     gui.Controllers.Files.HandleCommitEditorPress,
			Description: gui.c.Tr.CommitChangesWithEditor,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Return),
			Handler:     gui.handleEscapeMerge,
			Description: gui.c.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Files.OpenMergeTool),
			Handler:     gui.Controllers.Files.OpenMergeTool,
			Description: gui.c.Tr.LcOpenMergeTool,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Select),
			Handler:     gui.handlePickHunk,
			Description: gui.c.Tr.PickHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Main.PickBothHunks),
			Handler:     gui.handlePickAllHunks,
			Description: gui.c.Tr.PickAllHunks,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.handleSelectPrevConflict,
			Description: gui.c.Tr.PrevConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.NextBlock),
			Handler:     gui.handleSelectNextConflict,
			Description: gui.c.Tr.NextConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.PrevItem),
			Handler:     gui.handleSelectPrevConflictHunk,
			Description: gui.c.Tr.SelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.NextItem),
			Handler:     gui.handleSelectNextConflictHunk,
			Description: gui.c.Tr.SelectNextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextBlockAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectPrevConflictHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.handleSelectNextConflictHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(context.MAIN_MERGING_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Undo),
			Handler:     gui.handleMergeConflictUndo,
			Description: gui.c.Tr.LcUndo,
		},
		{
			ViewName: "branches",
			Contexts: []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.Select),
			// gonna use the exact same handler as the 'n' keybinding because everybody wants this to happen when they checkout a remote branch
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.c.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     guards.OutsideFilterMode(gui.handleMergeRemoteBranch),
			Description: gui.c.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteRemoteBranch,
			Description: gui.c.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.RebaseBranch),
			Handler:     guards.OutsideFilterMode(gui.handleRebaseOntoRemoteBranch),
			Description: gui.c.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(context.REMOTE_BRANCHES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Branches.SetUpstream),
			Handler:     gui.handleSetBranchUpstream,
			Description: gui.c.Tr.LcSetUpstream,
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
			ViewName:    "files",
			Contexts:    []string{string(context.SUBMODULES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.c.Tr.LcCopySubmoduleNameToClipboard,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(context.FILES_CONTEXT_KEY)},
			Key:         gui.getKey(config.Universal.ToggleWhitespaceInDiffView),
			Handler:     gui.toggleWhitespaceInDiffView,
			Description: gui.c.Tr.ToggleWhitespaceInDiffView,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.IncreaseContextInDiffView),
			Handler:     gui.IncreaseContextInDiffView,
			Description: gui.c.Tr.IncreaseContextInDiffView,
		},
		{
			ViewName:    "",
			Key:         gui.getKey(config.Universal.DecreaseContextInDiffView),
			Handler:     gui.DecreaseContextInDiffView,
			Description: gui.c.Tr.DecreaseContextInDiffView,
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
			ViewName:    "extras",
			Key:         gui.getKey(config.Universal.ExtrasMenu),
			Handler:     gui.handleCreateExtrasMenuPanel,
			Description: gui.c.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.PrevItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextItem),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(context.COMMAND_LOG_CONTEXT_KEY)},
			Key:      gui.getKey(config.Universal.NextItemAlt),
			Modifier: gocui.ModNone,
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Key:      gocui.MouseLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.handleFocusCommandLog,
		},
	}

	for _, controller := range []types.IController{
		gui.Controllers.LocalCommits,
		gui.Controllers.Submodules,
		gui.Controllers.Files,
		gui.Controllers.Remotes,
		gui.Controllers.Menu,
		gui.Controllers.Bisect,
		gui.Controllers.Undo,
		gui.Controllers.Sync,
		gui.Controllers.Tags,
	} {
		context := controller.Context()
		viewName := ""
		var contextKeys []string
		// nil context means global keybinding
		if context != nil {
			viewName = context.GetViewName()
			contextKeys = []string{string(context.GetKey())}
		}

		for _, binding := range controller.Keybindings(gui.getKey, config, guards) {
			binding.Contexts = contextKeys
			binding.ViewName = viewName
			bindings = append(bindings, binding)
		}
	}

	// while migrating we'll continue providing keybindings from the list contexts themselves.
	// for each controller we add above we need to remove the corresponding list context from here.
	for _, listContext := range []types.IListContext{
		gui.State.Contexts.Branches,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.Suggestions,
	} {
		viewName := listContext.GetViewName()
		contextKey := listContext.GetKey()
		for _, binding := range listContext.Keybindings(gui.getKey, config, guards) {
			binding.Contexts = []string{string(contextKey)}
			binding.ViewName = viewName
			bindings = append(bindings, binding)
		}
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*types.Binding{
			{ViewName: viewName, Key: gui.getKey(config.Universal.PrevBlock), Modifier: gocui.ModNone, Handler: gui.previousSideWindow},
			{ViewName: viewName, Key: gui.getKey(config.Universal.NextBlock), Modifier: gocui.ModNone, Handler: gui.nextSideWindow},
			{ViewName: viewName, Key: gui.getKey(config.Universal.PrevBlockAlt), Modifier: gocui.ModNone, Handler: gui.previousSideWindow},
			{ViewName: viewName, Key: gui.getKey(config.Universal.NextBlockAlt), Modifier: gocui.ModNone, Handler: gui.nextSideWindow},
			{ViewName: viewName, Key: gui.getKey(config.Universal.PrevBlockAlt2), Modifier: gocui.ModNone, Handler: gui.previousSideWindow},
			{ViewName: viewName, Key: gui.getKey(config.Universal.NextBlockAlt2), Modifier: gocui.ModNone, Handler: gui.nextSideWindow},
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
				Key:      gui.getKey(config.Universal.JumpToBlock[i]),
				Modifier: gocui.ModNone,
				Handler:  gui.goToSideWindow(window)})
		}
	}

	for viewName := range gui.State.Contexts.InitialViewTabContextMap() {
		bindings = append(bindings, []*types.Binding{
			{
				ViewName:    viewName,
				Key:         gui.getKey(config.Universal.NextTab),
				Handler:     gui.handleNextTab,
				Description: gui.c.Tr.LcNextTab,
				Tag:         "navigation",
			},
			{
				ViewName:    viewName,
				Key:         gui.getKey(config.Universal.PrevTab),
				Handler:     gui.handlePrevTab,
				Description: gui.c.Tr.LcPrevTab,
				Tag:         "navigation",
			},
		}...)
	}

	return bindings
}

func (gui *Gui) keybindings() error {
	bindings := gui.GetCustomCommandKeybindings()

	bindings = append(bindings, gui.GetInitialKeybindings()...)

	for _, binding := range bindings {
		if err := gui.SetKeybinding(binding); err != nil {
			return err
		}
	}

	for viewName := range gui.State.Contexts.InitialViewTabContextMap() {
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
	if isMouseKey(binding.Key) {
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

func isMouseKey(key interface{}) bool {
	return key == gocui.MouseLeft || key == gocui.MouseRight || key == gocui.MouseMiddle || key == gocui.MouseRelease || key == gocui.MouseWheelUp || key == gocui.MouseWheelDown || key == gocui.MouseWheelLeft || key == gocui.MouseWheelRight
}
