package gui

import (
	"fmt"
	"log"
	"strings"

	"unicode/utf8"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/constants"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Contexts    []string
	Handler     func() error
	KeyMod      gocui.KeyMod
	Description string
	Alternative string
	Tag         string // e.g. 'navigation'. Used for grouping things in the cheatsheet
	OpensMenu   bool
}

// GetDisplayStrings returns the display string of a file
func (b *Binding) GetDisplayStrings(isFocused bool) []string {
	return []string{GetKeyDisplay(b.KeyMod), b.Description}
}

var keyMapReversed = map[gocui.KeyMod]string{
	{Key: gocui.KeyF1, Modifier: gocui.ModNone}:         "f1",
	{Key: gocui.KeyF2, Modifier: gocui.ModNone}:         "f2",
	{Key: gocui.KeyF3, Modifier: gocui.ModNone}:         "f3",
	{Key: gocui.KeyF4, Modifier: gocui.ModNone}:         "f4",
	{Key: gocui.KeyF5, Modifier: gocui.ModNone}:         "f5",
	{Key: gocui.KeyF6, Modifier: gocui.ModNone}:         "f6",
	{Key: gocui.KeyF7, Modifier: gocui.ModNone}:         "f7",
	{Key: gocui.KeyF8, Modifier: gocui.ModNone}:         "f8",
	{Key: gocui.KeyF9, Modifier: gocui.ModNone}:         "f9",
	{Key: gocui.KeyF10, Modifier: gocui.ModNone}:        "f10",
	{Key: gocui.KeyF11, Modifier: gocui.ModNone}:        "f11",
	{Key: gocui.KeyF12, Modifier: gocui.ModNone}:        "f12",
	{Key: gocui.KeyInsert, Modifier: gocui.ModNone}:     "insert",
	{Key: gocui.KeyDelete, Modifier: gocui.ModNone}:     "delete",
	{Key: gocui.KeyHome, Modifier: gocui.ModNone}:       "home",
	{Key: gocui.KeyEnd, Modifier: gocui.ModNone}:        "end",
	{Key: gocui.KeyPgup, Modifier: gocui.ModNone}:       "pgup",
	{Key: gocui.KeyPgdn, Modifier: gocui.ModNone}:       "pgdown",
	{Key: gocui.KeyArrowUp, Modifier: gocui.ModNone}:    "▲",
	{Key: gocui.KeyArrowDown, Modifier: gocui.ModNone}:  "▼",
	{Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone}:  "◄",
	{Key: gocui.KeyArrowRight, Modifier: gocui.ModNone}: "►",
	{Key: gocui.KeyTab, Modifier: gocui.ModNone}:        "tab", // ctrl+i
	{Key: gocui.KeyBacktab, Modifier: gocui.ModNone}:    "shift+tab",
	{Key: gocui.KeyEnter, Modifier: gocui.ModNone}:      "enter", // ctrl+m
	{Key: gocui.KeyAltEnter, Modifier: gocui.ModAlt}:    "alt+enter",
	{Key: gocui.KeyEsc, Modifier: gocui.ModNone}:        "esc",        // ctrl+[, ctrl+3
	{Key: gocui.KeyBackspace, Modifier: gocui.ModNone}:  "backspace",  // ctrl+h
	{Key: gocui.KeySpace, Modifier: gocui.ModCtrl}:      "ctrl+space", // ctrl+~, ctrl+2
	{Key: gocui.KeyCtrlSlash, Modifier: gocui.ModCtrl}:  "ctrl+/",     // ctrl+_
	{Key: gocui.KeySpace, Modifier: gocui.ModNone}:      "space",
	{Key: gocui.KeyCtrlA, Modifier: gocui.ModCtrl}:      "ctrl+a",
	{Key: gocui.KeyCtrlB, Modifier: gocui.ModCtrl}:      "ctrl+b",
	{Key: gocui.KeyCtrlC, Modifier: gocui.ModCtrl}:      "ctrl+c",
	{Key: gocui.KeyCtrlD, Modifier: gocui.ModCtrl}:      "ctrl+d",
	{Key: gocui.KeyCtrlE, Modifier: gocui.ModCtrl}:      "ctrl+e",
	{Key: gocui.KeyCtrlF, Modifier: gocui.ModCtrl}:      "ctrl+f",
	{Key: gocui.KeyCtrlG, Modifier: gocui.ModCtrl}:      "ctrl+g",
	{Key: gocui.KeyCtrlJ, Modifier: gocui.ModCtrl}:      "ctrl+j",
	{Key: gocui.KeyCtrlK, Modifier: gocui.ModCtrl}:      "ctrl+k",
	{Key: gocui.KeyCtrlL, Modifier: gocui.ModCtrl}:      "ctrl+l",
	{Key: gocui.KeyCtrlN, Modifier: gocui.ModCtrl}:      "ctrl+n",
	{Key: gocui.KeyCtrlO, Modifier: gocui.ModCtrl}:      "ctrl+o",
	{Key: gocui.KeyCtrlP, Modifier: gocui.ModCtrl}:      "ctrl+p",
	{Key: gocui.KeyCtrlQ, Modifier: gocui.ModCtrl}:      "ctrl+q",
	{Key: gocui.KeyCtrlR, Modifier: gocui.ModCtrl}:      "ctrl+r",
	{Key: gocui.KeyCtrlS, Modifier: gocui.ModCtrl}:      "ctrl+s",
	{Key: gocui.KeyCtrlT, Modifier: gocui.ModCtrl}:      "ctrl+t",
	{Key: gocui.KeyCtrlU, Modifier: gocui.ModCtrl}:      "ctrl+u",
	{Key: gocui.KeyCtrlV, Modifier: gocui.ModCtrl}:      "ctrl+v",
	{Key: gocui.KeyCtrlW, Modifier: gocui.ModCtrl}:      "ctrl+w",
	{Key: gocui.KeyCtrlX, Modifier: gocui.ModCtrl}:      "ctrl+x",
	{Key: gocui.KeyCtrlY, Modifier: gocui.ModCtrl}:      "ctrl+y",
	{Key: gocui.KeyCtrlZ, Modifier: gocui.ModCtrl}:      "ctrl+z",
	{Key: gocui.KeyCtrl4, Modifier: gocui.ModCtrl}:      "ctrl+4", // ctrl+\
	{Key: gocui.KeyCtrl5, Modifier: gocui.ModCtrl}:      "ctrl+5", // ctrl+]
	{Key: gocui.KeyCtrl6, Modifier: gocui.ModCtrl}:      "ctrl+6",
	{Key: gocui.KeyCtrl8, Modifier: gocui.ModCtrl}:      "ctrl+8",
	{Key: gocui.KeyArrowRight, Modifier: gocui.ModCtrl}: "ctrl+►",
	{Key: gocui.KeyArrowLeft, Modifier: gocui.ModCtrl}:  "ctrl+◄",
	{Key: gocui.KeyArrowUp, Modifier: gocui.ModCtrl}:    "ctrl+▲",
	{Key: gocui.KeyArrowDown, Modifier: gocui.ModCtrl}:  "ctrl+▼",
}

var keymap = map[string]gocui.KeyMod{
	"<c-a>":       {Key: gocui.KeyCtrlA, Modifier: gocui.ModCtrl},
	"<c-b>":       {Key: gocui.KeyCtrlB, Modifier: gocui.ModCtrl},
	"<c-c>":       {Key: gocui.KeyCtrlC, Modifier: gocui.ModCtrl},
	"<c-d>":       {Key: gocui.KeyCtrlD, Modifier: gocui.ModCtrl},
	"<c-e>":       {Key: gocui.KeyCtrlE, Modifier: gocui.ModCtrl},
	"<c-f>":       {Key: gocui.KeyCtrlF, Modifier: gocui.ModCtrl},
	"<c-g>":       {Key: gocui.KeyCtrlG, Modifier: gocui.ModCtrl},
	"<c-h>":       {Key: gocui.KeyCtrlH, Modifier: gocui.ModCtrl},
	"<c-i>":       {Key: gocui.KeyCtrlI, Modifier: gocui.ModCtrl},
	"<c-j>":       {Key: gocui.KeyCtrlJ, Modifier: gocui.ModCtrl},
	"<c-k>":       {Key: gocui.KeyCtrlK, Modifier: gocui.ModCtrl},
	"<c-l>":       {Key: gocui.KeyCtrlL, Modifier: gocui.ModCtrl},
	"<c-m>":       {Key: gocui.KeyCtrlM, Modifier: gocui.ModCtrl},
	"<c-n>":       {Key: gocui.KeyCtrlN, Modifier: gocui.ModCtrl},
	"<c-o>":       {Key: gocui.KeyCtrlO, Modifier: gocui.ModCtrl},
	"<c-p>":       {Key: gocui.KeyCtrlP, Modifier: gocui.ModCtrl},
	"<c-q>":       {Key: gocui.KeyCtrlQ, Modifier: gocui.ModCtrl},
	"<c-r>":       {Key: gocui.KeyCtrlR, Modifier: gocui.ModCtrl},
	"<c-s>":       {Key: gocui.KeyCtrlS, Modifier: gocui.ModCtrl},
	"<c-t>":       {Key: gocui.KeyCtrlT, Modifier: gocui.ModCtrl},
	"<c-u>":       {Key: gocui.KeyCtrlU, Modifier: gocui.ModCtrl},
	"<c-v>":       {Key: gocui.KeyCtrlV, Modifier: gocui.ModCtrl},
	"<c-w>":       {Key: gocui.KeyCtrlW, Modifier: gocui.ModCtrl},
	"<c-x>":       {Key: gocui.KeyCtrlX, Modifier: gocui.ModCtrl},
	"<c-y>":       {Key: gocui.KeyCtrlY, Modifier: gocui.ModCtrl},
	"<c-z>":       {Key: gocui.KeyCtrlZ, Modifier: gocui.ModCtrl},
	"<c-~>":       {Key: gocui.KeyCtrlTilde, Modifier: gocui.ModCtrl},
	"<c-2>":       {Key: gocui.KeyCtrl2, Modifier: gocui.ModCtrl},
	"<c-3>":       {Key: gocui.KeyCtrl3, Modifier: gocui.ModCtrl},
	"<c-4>":       {Key: gocui.KeyCtrl4, Modifier: gocui.ModCtrl},
	"<c-5>":       {Key: gocui.KeyCtrl5, Modifier: gocui.ModCtrl},
	"<c-6>":       {Key: gocui.KeyCtrl6, Modifier: gocui.ModCtrl},
	"<c-7>":       {Key: gocui.KeyCtrl7, Modifier: gocui.ModCtrl},
	"<c-8>":       {Key: gocui.KeyCtrl8, Modifier: gocui.ModCtrl},
	"<c-space>":   {Key: gocui.KeyCtrlSpace, Modifier: gocui.ModCtrl},
	"<c-\\>":      {Key: gocui.KeyCtrlBackslash, Modifier: gocui.ModCtrl},
	"<c-[>":       {Key: gocui.KeyCtrlLsqBracket, Modifier: gocui.ModCtrl},
	"<c-]>":       {Key: gocui.KeyCtrlRsqBracket, Modifier: gocui.ModCtrl},
	"<c-/>":       {Key: gocui.KeyCtrlSlash, Modifier: gocui.ModCtrl},
	"<c-_>":       {Key: gocui.KeyCtrlUnderscore, Modifier: gocui.ModCtrl},
	"<backspace>": {Key: gocui.KeyBackspace, Modifier: gocui.ModNone},
	"<tab>":       {Key: gocui.KeyTab, Modifier: gocui.ModNone},
	"<backtab>":   {Key: gocui.KeyBacktab, Modifier: gocui.ModNone},
	"<enter>":     {Key: gocui.KeyEnter, Modifier: gocui.ModNone},
	"<a-enter>":   {Key: gocui.KeyEnter, Modifier: gocui.ModAlt},
	"<esc>":       {Key: gocui.KeyEsc, Modifier: gocui.ModNone},
	"<space>":     {Key: gocui.KeySpace, Modifier: gocui.ModNone},
	"<f1>":        {Key: gocui.KeyF1, Modifier: gocui.ModNone},
	"<f2>":        {Key: gocui.KeyF2, Modifier: gocui.ModNone},
	"<f3>":        {Key: gocui.KeyF3, Modifier: gocui.ModNone},
	"<f4>":        {Key: gocui.KeyF4, Modifier: gocui.ModNone},
	"<f5>":        {Key: gocui.KeyF5, Modifier: gocui.ModNone},
	"<f6>":        {Key: gocui.KeyF6, Modifier: gocui.ModNone},
	"<f7>":        {Key: gocui.KeyF7, Modifier: gocui.ModNone},
	"<f8>":        {Key: gocui.KeyF8, Modifier: gocui.ModNone},
	"<f9>":        {Key: gocui.KeyF9, Modifier: gocui.ModNone},
	"<f10>":       {Key: gocui.KeyF10, Modifier: gocui.ModNone},
	"<f11>":       {Key: gocui.KeyF11, Modifier: gocui.ModNone},
	"<f12>":       {Key: gocui.KeyF12, Modifier: gocui.ModNone},
	"<insert>":    {Key: gocui.KeyInsert, Modifier: gocui.ModNone},
	"<delete>":    {Key: gocui.KeyDelete, Modifier: gocui.ModNone},
	"<home>":      {Key: gocui.KeyHome, Modifier: gocui.ModNone},
	"<end>":       {Key: gocui.KeyEnd, Modifier: gocui.ModNone},
	"<pgup>":      {Key: gocui.KeyPgup, Modifier: gocui.ModNone},
	"<pgdown>":    {Key: gocui.KeyPgdn, Modifier: gocui.ModNone},
	"<up>":        {Key: gocui.KeyArrowUp, Modifier: gocui.ModNone},
	"<down>":      {Key: gocui.KeyArrowDown, Modifier: gocui.ModNone},
	"<left>":      {Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone},
	"<right>":     {Key: gocui.KeyArrowRight, Modifier: gocui.ModNone},
	"<c-right>":   {Key: gocui.KeyArrowRight, Modifier: gocui.ModCtrl},
	"<c-left>":    {Key: gocui.KeyArrowLeft, Modifier: gocui.ModCtrl},
	"<c-up>":      {Key: gocui.KeyArrowUp, Modifier: gocui.ModCtrl},
	"<c-down>":    {Key: gocui.KeyArrowDown, Modifier: gocui.ModCtrl},
}

func (gui *Gui) getKeyDisplay(name string) string {
	key := gui.getKey(name)
	return GetKeyDisplay(key)
}

func GetKeyDisplay(keyMod gocui.KeyMod) string {
	keyInt := 0

	value, ok := keyMapReversed[keyMod]
	if ok {
		return value
	}

	switch key := keyMod.Key.(type) {
	case rune:
		keyInt = int(key)
	case gocui.Key:
		keyInt = int(key)
	}

	return fmt.Sprintf("%c", keyInt)
}

func (gui *Gui) getKey(key string) gocui.KeyMod {
	runeCount := utf8.RuneCountInString(key)
	if runeCount > 1 {
		binding, found := keymap[strings.ToLower(key)]
		if !found {
			log.Fatalf("Unrecognized key %s for keybinding. For permitted values see %s", strings.ToLower(key), constants.Links.Docs.CustomKeybindings)
		} else {
			return binding
		}
	} else if runeCount == 1 {
		return gocui.KeyMod{Key: []rune(key)[0], Modifier: gocui.ModNone}
	}
	log.Fatal("Key empty for keybinding: " + strings.ToLower(key))
	return gocui.KeyMod{Key: nil, Modifier: gocui.ModNone}
}

// GetInitialKeybindings is a function.
func (gui *Gui) GetInitialKeybindings() []*Binding {
	config := gui.UserConfig.Keybinding

	bindings := []*Binding{
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.Quit),
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.QuitWithoutChangingDirectory),
			Handler:  gui.handleQuitWithoutChangingDirectory,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.QuitAlt1),
			Handler:  gui.handleQuit,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.Return),
			Handler:  gui.handleTopLevelReturn,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.OpenRecentRepos),
			Handler:     gui.handleCreateRecentReposMenu,
			Alternative: "<c-r>",
			Description: gui.Tr.SwitchRepo,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.ScrollUpMain),
			Handler:     gui.scrollUpMain,
			Alternative: "fn+up",
			Description: gui.Tr.LcScrollUpMainPanel,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.ScrollDownMain),
			Handler:     gui.scrollDownMain,
			Alternative: "fn+down",
			Description: gui.Tr.LcScrollDownMainPanel,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.ScrollUpMainAlt1),
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.ScrollDownMainAlt1),
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.ScrollUpMainAlt2),
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.ScrollDownMainAlt2),
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.CreateRebaseOptionsMenu),
			Handler:     gui.handleCreateRebaseOptionsMenu,
			Description: gui.Tr.ViewMergeRebaseOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.CreatePatchOptionsMenu),
			Handler:     gui.handleCreatePatchOptionsMenu,
			Description: gui.Tr.ViewPatchOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.PushFiles),
			Handler:     gui.pushFiles,
			Description: gui.Tr.LcPush,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.PullFiles),
			Handler:     gui.handlePullFiles,
			Description: gui.Tr.LcPull,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.Refresh),
			Handler:     gui.handleRefresh,
			Description: gui.Tr.LcRefresh,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.OptionMenu),
			Handler:     gui.handleCreateOptionsMenu,
			Description: gui.Tr.LcOpenMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "",
			KeyMod:   gui.getKey(config.Universal.OptionMenuAlt1),
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName: "",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseMiddle, Modifier: gocui.ModNone},
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.Undo),
			Handler:     gui.reflogUndo,
			Description: gui.Tr.LcUndoReflog,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.Redo),
			Handler:     gui.reflogRedo,
			Description: gui.Tr.LcRedoReflog,
		},
		{
			ViewName:    "status",
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditConfig,
			Description: gui.Tr.EditConfig,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.NextScreenMode),
			Handler:     gui.nextScreenMode,
			Description: gui.Tr.LcNextScreenMode,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.PrevScreenMode),
			Handler:     gui.prevScreenMode,
			Description: gui.Tr.LcPrevScreenMode,
		},
		{
			ViewName:    "status",
			KeyMod:      gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenConfig,
			Description: gui.Tr.OpenConfig,
		},
		{
			ViewName:    "status",
			KeyMod:      gui.getKey(config.Status.CheckForUpdate),
			Handler:     gui.handleCheckForUpdate,
			Description: gui.Tr.LcCheckForUpdate,
		},
		{
			ViewName:    "status",
			KeyMod:      gui.getKey(config.Status.RecentRepos),
			Handler:     gui.handleCreateRecentReposMenu,
			Description: gui.Tr.SwitchRepo,
		},
		{
			ViewName:    "status",
			KeyMod:      gui.getKey(config.Status.AllBranchesLogGraph),
			Handler:     gui.handleShowAllBranchLogs,
			Description: gui.Tr.LcAllBranchesLogGraph,
		},
		{
			ViewName:    "files",
			KeyMod:      gui.getKey("<c-b>"),
			Handler:     gui.handleStatusFilterPressed,
			Description: gui.Tr.LcCommitFileFilter,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChanges),
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.CommitChanges,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChangesWithoutHook),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.LcCommitChangesWithoutHook,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.AmendLastCommit),
			Handler:     gui.handleAmendCommitPress,
			Description: gui.Tr.AmendLastCommit,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChangesWithEditor),
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.CommitChangesWithEditor,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleFilePress,
			Description: gui.Tr.LcToggleStaged,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleCreateDiscardMenu,
			Description: gui.Tr.LcViewDiscardOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleFileEdit,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.IgnoreFile),
			Handler:     gui.handleIgnoreFile,
			Description: gui.Tr.LcIgnoreFile,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.RefreshFiles),
			Handler:     gui.handleRefreshFiles,
			Description: gui.Tr.LcRefreshFiles,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.StashAllChanges),
			Handler:     gui.handleStashChanges,
			Description: gui.Tr.LcStashAllChanges,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.ViewStashOptions),
			Handler:     gui.handleCreateStashMenu,
			Description: gui.Tr.LcViewStashOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.ToggleStagedAll),
			Handler:     gui.handleStageAll,
			Description: gui.Tr.LcToggleStagedAll,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.ViewResetOptions),
			Handler:     gui.handleCreateResetMenu,
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleEnterFile,
			Description: gui.Tr.FileEnter,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.Fetch),
			Handler:     gui.handleGitFetch,
			Description: gui.Tr.LcFetch,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyFileNameToClipboard,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.ExecuteCustomCommand),
			Handler:     gui.handleCustomCommand,
			Description: gui.Tr.LcExecuteCustomCommand,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToUpstreamMenu,
			Description: gui.Tr.LcViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.ToggleTreeView),
			Handler:     gui.handleToggleFileTreeView,
			Description: gui.Tr.LcToggleTreeView,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.OpenMergeTool),
			Handler:     gui.handleOpenMergeTool,
			Description: gui.Tr.LcOpenMergeTool,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleBranchPress,
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.CreatePullRequest),
			Handler:     gui.handleCreatePullRequestPress,
			Description: gui.Tr.LcCreatePullRequest,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.ViewPullRequestOptions),
			Handler:     gui.handleCreatePullRequestMenu,
			Description: gui.Tr.LcCreatePullRequestOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.CopyPullRequestURL),
			Handler:     gui.handleCopyPullRequestURLPress,
			Description: gui.Tr.LcCopyPullRequestURL,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.CheckoutBranchByName),
			Handler:     gui.handleCheckoutByName,
			Description: gui.Tr.LcCheckoutByName,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.ForceCheckoutBranch),
			Handler:     gui.handleForceCheckout,
			Description: gui.Tr.LcForceCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteBranch,
			Description: gui.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.RebaseBranch),
			Handler:     gui.handleRebaseOntoLocalBranch,
			Description: gui.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     gui.handleMerge,
			Description: gui.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.ViewGitFlowOptions),
			Handler:     gui.handleCreateGitFlowMenu,
			Description: gui.Tr.LcGitFlowOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.FastForward),
			Handler:     gui.handleFastForward,
			Description: gui.Tr.FastForward,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToBranchMenu,
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.RenameBranch),
			Handler:     gui.handleRenameBranch,
			Description: gui.Tr.LcRenameBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyBranchNameToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(LOCAL_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleSwitchToSubCommits,
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.withSelectedTag(gui.handleCheckoutTag),
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.withSelectedTag(gui.handleDeleteTag),
			Description: gui.Tr.LcDeleteTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.PushTag),
			Handler:     gui.withSelectedTag(gui.handlePushTag),
			Description: gui.Tr.LcPushTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleCreateTag,
			Description: gui.Tr.LcCreateTag,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.withSelectedTag(gui.handleCreateResetToTagMenu),
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(TAGS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleSwitchToSubCommits,
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Return),
			Handler:     gui.handleRemoteBranchesEscape,
			Description: gui.Tr.ReturnToRemotesList,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateResetToRemoteBranchMenu,
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleSwitchToSubCommits,
			Description: gui.Tr.LcViewCommits,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.FetchRemote),
			Handler:     gui.handleFetchRemote,
			Description: gui.Tr.LcFetchRemote,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.OpenLogMenu),
			Handler:     gui.handleOpenLogMenu,
			Description: gui.Tr.LcOpenLogMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.SquashDown),
			Handler:     gui.handleCommitSquashDown,
			Description: gui.Tr.LcSquashDown,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.RenameCommit),
			Handler:     gui.handleRewordCommit,
			Description: gui.Tr.LcRewordCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.RenameCommitWithEditor),
			Handler:     gui.handleRewordCommitEditor,
			Description: gui.Tr.LcRenameCommitEditor,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateCommitResetMenu,
			Description: gui.Tr.LcResetToThisCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.MarkCommitAsFixup),
			Handler:     gui.handleCommitFixup,
			Description: gui.Tr.LcFixupCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CreateFixupCommit),
			Handler:     gui.handleCreateFixupCommit,
			Description: gui.Tr.LcCreateFixupCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.SquashAboveCommits),
			Handler:     gui.handleSquashAllAboveFixupCommits,
			Description: gui.Tr.LcSquashAboveCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleCommitDelete,
			Description: gui.Tr.LcDeleteCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.MoveDownCommit),
			Handler:     gui.handleCommitMoveDown,
			Description: gui.Tr.LcMoveDownCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.MoveUpCommit),
			Handler:     gui.handleCommitMoveUp,
			Description: gui.Tr.LcMoveUpCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleCommitEdit,
			Description: gui.Tr.LcEditCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.AmendToCommit),
			Handler:     gui.handleCommitAmendTo,
			Description: gui.Tr.LcAmendToCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.PickCommit),
			Handler:     gui.handleCommitPick,
			Description: gui.Tr.LcPickCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.RevertCommit),
			Handler:     gui.handleCommitRevert,
			Description: gui.Tr.LcRevertCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.PasteCommits),
			Handler:     gui.HandlePasteCommits,
			Description: gui.Tr.LcPasteCommits,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewCommitFiles,
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CheckoutCommit),
			Handler:     gui.handleCheckoutCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcCreateNewBranchFromCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.TagCommit),
			Handler:     gui.handleTagCommit,
			Description: gui.Tr.LcTagCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CopyCommitMessageToClipboard),
			Handler:     gui.handleCopySelectedCommitMessageToClipboard,
			Description: gui.Tr.LcCopyCommitMessageToClipboard,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.OpenInBrowser),
			Handler:     gui.handleOpenCommitInBrowser,
			Description: gui.Tr.LcOpenCommitInBrowser,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(BRANCH_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewBisectOptions),
			Handler:     gui.handleOpenBisectMenu,
			Description: gui.Tr.LcViewBisectOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewReflogCommitFiles,
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutReflogCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateReflogResetMenu,
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "commits",
			Contexts:    []string{string(REFLOG_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewSubCommitFiles,
			Description: gui.Tr.LcViewCommitFiles,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleCheckoutSubCommit,
			Description: gui.Tr.LcCheckoutCommit,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ViewResetOptions),
			Handler:     gui.handleCreateSubCommitResetMenu,
			Description: gui.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopy),
			Handler:     gui.handleCopyCommit,
			Description: gui.Tr.LcCherryPickCopy,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.CherryPickCopyRange),
			Handler:     gui.handleCopyCommitRange,
			Description: gui.Tr.LcCherryPickCopyRange,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Commits.ResetCherryPick),
			Handler:     gui.exitCherryPickingMode,
			Description: gui.Tr.LcResetCherryPick,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(SUB_COMMITS_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyCommitShaToClipboard,
		},
		{
			ViewName:    "stash",
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleViewStashFiles,
			Description: gui.Tr.LcViewStashFiles,
		},
		{
			ViewName:    "stash",
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleStashApply,
			Description: gui.Tr.LcApply,
		},
		{
			ViewName:    "stash",
			KeyMod:      gui.getKey(config.Stash.PopStash),
			Handler:     gui.handleStashPop,
			Description: gui.Tr.LcPop,
		},
		{
			ViewName:    "stash",
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleStashDrop,
			Description: gui.Tr.LcDrop,
		},
		{
			ViewName:    "stash",
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName: "commitMessage",
			KeyMod:   gui.getKey(config.Universal.SubmitEditorText),
			Handler:  gui.handleCommitConfirm,
		},
		{
			ViewName: "commitMessage",
			KeyMod:   gui.getKey(config.Universal.Return),
			Handler:  gui.handleCommitClose,
		},
		{
			ViewName: "credentials",
			KeyMod:   gui.getKey(config.Universal.Confirm),
			Handler:  gui.handleSubmitCredential,
		},
		{
			ViewName: "credentials",
			KeyMod:   gui.getKey(config.Universal.Return),
			Handler:  gui.handleCloseCredentialsView,
		},
		{
			ViewName:    "menu",
			KeyMod:      gui.getKey(config.Universal.Return),
			Handler:     gui.handleMenuClose,
			Description: gui.Tr.LcCloseMenu,
		},
		{
			ViewName: "information",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleInfoClick,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopyCommitFileNameToClipboard,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.CommitFiles.CheckoutCommitFile),
			Handler:     gui.handleCheckoutCommitFile,
			Description: gui.Tr.LcCheckoutCommitFile,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDiscardOldFileChange,
			Description: gui.Tr.LcDiscardOldFileChange,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenOldCommitFile,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditCommitFile,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleFileForPatch,
			Description: gui.Tr.LcToggleAddToPatch,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.handleEnterCommitFile,
			Description: gui.Tr.LcEnterFile,
		},
		{
			ViewName:    "commitFiles",
			KeyMod:      gui.getKey(config.Files.ToggleTreeView),
			Handler:     gui.handleToggleCommitFileTreeView,
			Description: gui.Tr.LcToggleTreeView,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.FilteringMenu),
			Handler:     gui.handleCreateFilteringMenuPanel,
			Description: gui.Tr.LcOpenFilteringMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.DiffingMenu),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.DiffingMenuAlt),
			Handler:     gui.handleCreateDiffingMenuPanel,
			Description: gui.Tr.LcOpenDiffingMenu,
			OpensMenu:   true,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.ExtrasMenu),
			Handler:     gui.handleCreateExtrasMenuPanel,
			Description: gui.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "secondary",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelUp, Modifier: gocui.ModNone},
			Handler:  gui.scrollUpSecondary,
		},
		{
			ViewName: "secondary",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelDown, Modifier: gocui.ModNone},
			Handler:  gui.scrollDownSecondary,
		},
		{
			ViewName: "secondary",
			Contexts: []string{string(MAIN_NORMAL_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleMouseDownSecondary,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_NORMAL_CONTEXT_KEY)},
			KeyMod:      gocui.KeyMod{Key: gocui.MouseWheelDown, Modifier: gocui.ModNone},
			Handler:     gui.scrollDownMain,
			Description: gui.Tr.ScrollDown,
			Alternative: "fn+up",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_NORMAL_CONTEXT_KEY)},
			KeyMod:      gocui.KeyMod{Key: gocui.MouseWheelUp, Modifier: gocui.ModNone},
			Handler:     gui.scrollUpMain,
			Description: gui.Tr.ScrollUp,
			Alternative: "fn+down",
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_NORMAL_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleMouseDownMain,
		},
		{
			ViewName: "secondary",
			Contexts: []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleTogglePanelClick,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Return),
			Handler:     gui.handleStagingEscape,
			Description: gui.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleStagedSelection,
			Description: gui.Tr.StageSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleResetSelection,
			Description: gui.Tr.ResetSelection,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.TogglePanel),
			Handler:     gui.handleTogglePanel,
			Description: gui.Tr.TogglePanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Return),
			Handler:     gui.handleEscapePatchBuildingPanel,
			Description: gui.Tr.ExitLineByLineMode,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleOpenFileAtLine,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.PrevItem),
			Handler:     gui.handleSelectPrevLine,
			Description: gui.Tr.PrevLine,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.NextItem),
			Handler:     gui.handleSelectNextLine,
			Description: gui.Tr.NextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevItemAlt),
			Handler:  gui.handleSelectPrevLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextItemAlt),
			Handler:  gui.handleSelectNextLine,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelUp, Modifier: gocui.ModNone},
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelDown, Modifier: gocui.ModNone},
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.handleSelectPrevHunk,
			Description: gui.Tr.PrevHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevBlockAlt),
			Handler:  gui.handleSelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.NextBlock),
			Handler:     gui.handleSelectNextHunk,
			Description: gui.Tr.NextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextBlockAlt),
			Handler:  gui.handleSelectNextHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.copySelectedToClipboard,
			Description: gui.Tr.LcCopySelectedTexToClipboard,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleLineByLineEdit,
			Description: gui.Tr.LcEditFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.OpenFile),
			Handler:     gui.handleFileOpen,
			Description: gui.Tr.LcOpenFile,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.NextPage),
			Handler:     gui.handleLineByLineNextPage,
			Description: gui.Tr.LcNextPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.PrevPage),
			Handler:     gui.handleLineByLinePrevPage,
			Description: gui.Tr.LcPrevPage,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GotoTop),
			Handler:     gui.handleLineByLineGotoTop,
			Description: gui.Tr.LcGotoTop,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GotoBottom),
			Handler:     gui.handleLineByLineGotoBottom,
			Description: gui.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.StartSearch),
			Handler:     func() error { return gui.handleOpenSearch("main") },
			Description: gui.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handleToggleSelectionForPatch,
			Description: gui.Tr.ToggleSelectionForPatch,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Main.ToggleDragSelect),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.Tr.ToggleDragSelect,
		},
		// Alias 'V' -> 'v'
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Main.ToggleDragSelectAlt),
			Handler:     gui.handleToggleSelectRange,
			Description: gui.Tr.ToggleDragSelect,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Main.ToggleSelectHunk),
			Handler:     gui.handleToggleSelectHunk,
			Description: gui.Tr.ToggleSelectHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleLBLMouseDown,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleMouseDrag,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelUp, Modifier: gocui.ModNone},
			Handler:  gui.scrollUpMain,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelDown, Modifier: gocui.ModNone},
			Handler:  gui.scrollDownMain,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY), string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.ScrollLeft),
			Handler:     gui.scrollLeftMain,
			Description: gui.Tr.LcScrollLeft,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_PATCH_BUILDING_CONTEXT_KEY), string(MAIN_STAGING_CONTEXT_KEY), string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.ScrollRight),
			Handler:     gui.scrollRightMain,
			Description: gui.Tr.LcScrollRight,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChanges),
			Handler:     gui.handleCommitPress,
			Description: gui.Tr.CommitChanges,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChangesWithoutHook),
			Handler:     gui.handleWIPCommitPress,
			Description: gui.Tr.LcCommitChangesWithoutHook,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_STAGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.CommitChangesWithEditor),
			Handler:     gui.handleCommitEditorPress,
			Description: gui.Tr.CommitChangesWithEditor,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Return),
			Handler:     gui.handleEscapeMerge,
			Description: gui.Tr.ReturnToFilesPanel,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Files.OpenMergeTool),
			Handler:     gui.handleOpenMergeTool,
			Description: gui.Tr.LcOpenMergeTool,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Select),
			Handler:     gui.handlePickHunk,
			Description: gui.Tr.PickHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Main.PickBothHunks),
			Handler:     gui.handlePickAllHunks,
			Description: gui.Tr.PickAllHunks,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.PrevBlock),
			Handler:     gui.handleSelectPrevConflict,
			Description: gui.Tr.PrevConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.NextBlock),
			Handler:     gui.handleSelectNextConflict,
			Description: gui.Tr.NextConflict,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.PrevItem),
			Handler:     gui.handleSelectPrevConflictHunk,
			Description: gui.Tr.SelectPrevHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.NextItem),
			Handler:     gui.handleSelectNextConflictHunk,
			Description: gui.Tr.SelectNextHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevBlockAlt),
			Handler:  gui.handleSelectPrevConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextBlockAlt),
			Handler:  gui.handleSelectNextConflict,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevItemAlt),
			Handler:  gui.handleSelectPrevConflictHunk,
		},
		{
			ViewName: "main",
			Contexts: []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextItemAlt),
			Handler:  gui.handleSelectNextConflictHunk,
		},
		{
			ViewName:    "main",
			Contexts:    []string{string(MAIN_MERGING_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Undo),
			Handler:     gui.handleMergeConflictUndo,
			Description: gui.Tr.LcUndo,
		},
		{
			ViewName: "branches",
			Contexts: []string{string(REMOTES_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.GoInto),
			Handler:  gui.handleRemoteEnter,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleAddRemote,
			Description: gui.Tr.LcAddNewRemote,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleRemoveRemote,
			Description: gui.Tr.LcRemoveRemote,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.handleEditRemote,
			Description: gui.Tr.LcEditRemote,
		},
		{
			ViewName: "branches",
			Contexts: []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.Select),
			// gonna use the exact same handler as the 'n' keybinding because everybody wants this to happen when they checkout a remote branch
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcCheckout,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleNewBranchOffCurrentItem,
			Description: gui.Tr.LcNewBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.MergeIntoCurrentBranch),
			Handler:     gui.handleMergeRemoteBranch,
			Description: gui.Tr.LcMergeIntoCurrentBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.handleDeleteRemoteBranch,
			Description: gui.Tr.LcDeleteBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.RebaseBranch),
			Handler:     gui.handleRebaseOntoRemoteBranch,
			Description: gui.Tr.LcRebaseBranch,
		},
		{
			ViewName:    "branches",
			Contexts:    []string{string(REMOTE_BRANCHES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Branches.SetUpstream),
			Handler:     gui.handleSetBranchUpstream,
			Description: gui.Tr.LcSetUpstream,
		},
		{
			ViewName: "status",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleStatusClick,
		},
		{
			ViewName: "search",
			KeyMod:   gui.getKey(config.Universal.Confirm),
			Handler:  gui.handleSearch,
		},
		{
			ViewName: "search",
			KeyMod:   gui.getKey(config.Universal.Return),
			Handler:  gui.handleSearchEscape,
		},
		{
			ViewName: "confirmation",
			KeyMod:   gui.getKey(config.Universal.PrevItem),
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			KeyMod:   gui.getKey(config.Universal.NextItem),
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			KeyMod:   gui.getKey(config.Universal.PrevItemAlt),
			Handler:  gui.scrollUpConfirmationPanel,
		},
		{
			ViewName: "confirmation",
			KeyMod:   gui.getKey(config.Universal.NextItemAlt),
			Handler:  gui.scrollDownConfirmationPanel,
		},
		{
			ViewName: "menu",
			KeyMod:   gui.getKey(config.Universal.Select),
			Handler:  gui.onMenuPress,
		},
		{
			ViewName: "menu",
			KeyMod:   gui.getKey(config.Universal.Confirm),
			Handler:  gui.onMenuPress,
		},
		{
			ViewName: "menu",
			KeyMod:   gui.getKey(config.Universal.ConfirmAlt1),
			Handler:  gui.onMenuPress,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.CopyToClipboard),
			Handler:     gui.handleCopySelectedSideContextItemToClipboard,
			Description: gui.Tr.LcCopySubmoduleNameToClipboard,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.GoInto),
			Handler:     gui.forSubmodule(gui.handleSubmoduleEnter),
			Description: gui.Tr.LcEnterSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Remove),
			Handler:     gui.forSubmodule(gui.removeSubmodule),
			Description: gui.Tr.LcRemoveSubmodule,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Submodules.Update),
			Handler:     gui.forSubmodule(gui.handleUpdateSubmodule),
			Description: gui.Tr.LcSubmoduleUpdate,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.New),
			Handler:     gui.handleAddSubmodule,
			Description: gui.Tr.LcAddSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.Edit),
			Handler:     gui.forSubmodule(gui.handleEditSubmoduleUrl),
			Description: gui.Tr.LcEditSubmoduleUrl,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Submodules.Init),
			Handler:     gui.forSubmodule(gui.handleSubmoduleInit),
			Description: gui.Tr.LcInitSubmodule,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(SUBMODULES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Submodules.BulkMenu),
			Handler:     gui.handleBulkSubmoduleActionsMenu,
			Description: gui.Tr.LcViewBulkSubmoduleOptions,
			OpensMenu:   true,
		},
		{
			ViewName:    "files",
			Contexts:    []string{string(FILES_CONTEXT_KEY)},
			KeyMod:      gui.getKey(config.Universal.ToggleWhitespaceInDiffView),
			Handler:     gui.toggleWhitespaceInDiffView,
			Description: gui.Tr.ToggleWhitespaceInDiffView,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.IncreaseContextInDiffView),
			Handler:     gui.IncreaseContextInDiffView,
			Description: gui.Tr.IncreaseContextInDiffView,
		},
		{
			ViewName:    "",
			KeyMod:      gui.getKey(config.Universal.DecreaseContextInDiffView),
			Handler:     gui.DecreaseContextInDiffView,
			Description: gui.Tr.DecreaseContextInDiffView,
		},
		{
			ViewName: "extras",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelUp, Modifier: gocui.ModNone},
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseWheelDown, Modifier: gocui.ModNone},
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName:    "extras",
			KeyMod:      gui.getKey(config.Universal.ExtrasMenu),
			Handler:     gui.handleCreateExtrasMenuPanel,
			Description: gui.Tr.LcOpenExtrasMenu,
			OpensMenu:   true,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(COMMAND_LOG_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevItemAlt),
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(COMMAND_LOG_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.PrevItem),
			Handler:  gui.scrollUpExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(COMMAND_LOG_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextItem),
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			Contexts: []string{string(COMMAND_LOG_CONTEXT_KEY)},
			KeyMod:   gui.getKey(config.Universal.NextItemAlt),
			Handler:  gui.scrollDownExtra,
		},
		{
			ViewName: "extras",
			Tag:      "navigation",
			KeyMod:   gocui.KeyMod{Key: gocui.MouseLeft, Modifier: gocui.ModNone},
			Handler:  gui.handleFocusCommandLog,
		},
	}

	for _, viewName := range []string{"status", "branches", "files", "commits", "commitFiles", "stash", "menu"} {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.PrevBlock), Handler: gui.previousSideWindow},
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.NextBlock), Handler: gui.nextSideWindow},
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.PrevBlockAlt), Handler: gui.previousSideWindow},
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.NextBlockAlt), Handler: gui.nextSideWindow},
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.PrevBlockAlt2), Handler: gui.previousSideWindow},
			{ViewName: viewName, KeyMod: gui.getKey(config.Universal.NextBlockAlt2), Handler: gui.nextSideWindow},
		}...)
	}

	// Appends keybindings to jump to a particular sideView using numbers
	windows := []string{"status", "files", "branches", "commits", "stash"}

	if len(config.Universal.JumpToBlock) != len(windows) {
		log.Fatal("Jump to block keybindings cannot be set. Exactly 5 keybindings must be supplied.")
	} else {
		for i, window := range windows {
			bindings = append(bindings, &Binding{
				ViewName: "",
				KeyMod:   gui.getKey(config.Universal.JumpToBlock[i]),
				Handler:  gui.goToSideWindow(window)})
		}
	}

	for viewName := range gui.State.Contexts.initialViewTabContextMap() {
		bindings = append(bindings, []*Binding{
			{
				ViewName:    viewName,
				KeyMod:      gui.getKey(config.Universal.NextTab),
				Handler:     gui.handleNextTab,
				Description: gui.Tr.LcNextTab,
				Tag:         "navigation",
			},
			{
				ViewName:    viewName,
				KeyMod:      gui.getKey(config.Universal.PrevTab),
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
		if err := gui.g.SetKeybinding(binding.ViewName, binding.Contexts, binding.KeyMod, gui.wrappedHandler(binding.Handler)); err != nil {
			return err
		}
	}

	for viewName := range gui.State.Contexts.initialViewTabContextMap() {
		viewName := viewName
		tabClickCallback := func(tabIndex int) error { return gui.onViewTabClick(viewName, tabIndex) }

		if err := gui.g.SetTabClickBinding(viewName, tabClickCallback); err != nil {
			return err
		}
	}

	return nil
}
