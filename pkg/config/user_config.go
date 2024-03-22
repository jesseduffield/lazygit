package config

import (
	"time"

	"github.com/karimkhaleel/jsonschema"
)

type UserConfig struct {
	// Config relating to the Lazygit UI
	Gui GuiConfig `yaml:"gui"`
	// Config relating to git
	Git GitConfig `yaml:"git"`
	// Periodic update checks
	Update UpdateConfig `yaml:"update"`
	// Background refreshes
	Refresher RefresherConfig `yaml:"refresher"`
	// If true, show a confirmation popup before quitting Lazygit
	ConfirmOnQuit bool `yaml:"confirmOnQuit"`
	// If true, exit Lazygit when the user presses escape in a context where there is nothing to cancel/close
	QuitOnTopLevelReturn bool `yaml:"quitOnTopLevelReturn"`
	// Keybindings
	Keybinding KeybindingConfig `yaml:"keybinding"`
	// Config relating to things outside of Lazygit like how files are opened, copying to clipboard, etc
	OS OSConfig `yaml:"os,omitempty"`
	// If true, don't display introductory popups upon opening Lazygit.
	// Lazygit sets this to true upon first runninng the program so that you don't see introductory popups every time you open the program.
	DisableStartupPopups bool `yaml:"disableStartupPopups"`
	// User-configured commands that can be invoked from within Lazygit
	CustomCommands []CustomCommand `yaml:"customCommands" jsonschema:"uniqueItems=true"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-pull-request-urls
	Services map[string]string `yaml:"services"`
	// What to do when opening Lazygit outside of a git repo.
	// - 'prompt': (default) ask whether to initialize a new repo or open in the most recent repo
	// - 'create': initialize a new repo
	// - 'skip': open most recent repo
	// - 'quit': exit Lazygit
	NotARepository string `yaml:"notARepository" jsonschema:"enum=prompt,enum=create,enum=skip,enum=quit"`
	// If true, display a confirmation when subprocess terminates. This allows you to view the output of the subprocess before returning to Lazygit.
	PromptToReturnFromSubprocess bool `yaml:"promptToReturnFromSubprocess"`
}

type RefresherConfig struct {
	// File/submodule refresh interval in seconds.
	// Auto-refresh can be disabled via option 'git.autoRefresh'.
	RefreshInterval int `yaml:"refreshInterval" jsonschema:"minimum=0"`
	// Re-fetch interval in seconds.
	// Auto-fetch can be disabled via option 'git.autoFetch'.
	FetchInterval int `yaml:"fetchInterval" jsonschema:"minimum=0"`
}

type GuiConfig struct {
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-author-color
	AuthorColors map[string]string `yaml:"authorColors"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-branch-color
	BranchColors map[string]string `yaml:"branchColors"`
	// The number of lines you scroll by when scrolling the main window
	ScrollHeight int `yaml:"scrollHeight" jsonschema:"minimum=1"`
	// If true, allow scrolling past the bottom of the content in the main window
	ScrollPastBottom bool `yaml:"scrollPastBottom"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#scroll-off-margin
	ScrollOffMargin int `yaml:"scrollOffMargin"`
	// One of: 'margin' (default) | 'jump'
	ScrollOffBehavior string `yaml:"scrollOffBehavior"`
	// If true, capture mouse events.
	// When mouse events are captured, it's a little harder to select text: e.g. requiring you to hold the option key when on macOS.
	MouseEvents bool `yaml:"mouseEvents"`
	// If true, do not show a warning when discarding changes in the staging view.
	SkipDiscardChangeWarning bool `yaml:"skipDiscardChangeWarning"`
	// If true, do not show warning when applying/popping the stash
	SkipStashWarning bool `yaml:"skipStashWarning"`
	// If true, do not show a warning when attempting to commit without any staged files; instead stage all unstaged files.
	SkipNoStagedFilesWarning bool `yaml:"skipNoStagedFilesWarning"`
	// If true, do not show a warning when rewording a commit via an external editor
	SkipRewordInEditorWarning bool `yaml:"skipRewordInEditorWarning"`
	// Fraction of the total screen width to use for the left side section. You may want to pick a small number (e.g. 0.2) if you're using a narrow screen, so that you can see more of the main section.
	// Number from 0 to 1.0.
	SidePanelWidth float64 `yaml:"sidePanelWidth" jsonschema:"maximum=1,minimum=0"`
	// If true, increase the height of the focused side window; creating an accordion effect.
	ExpandFocusedSidePanel bool `yaml:"expandFocusedSidePanel"`
	// Sometimes the main window is split in two (e.g. when the selected file has both staged and unstaged changes). This setting controls how the two sections are split.
	// Options are:
	// - 'horizontal': split the window horizontally
	// - 'vertical': split the window vertically
	// - 'flexible': (default) split the window horizontally if the window is wide enough, otherwise split vertically
	MainPanelSplitMode string `yaml:"mainPanelSplitMode" jsonschema:"enum=horizontal,enum=flexible,enum=vertical"`
	// How the window is split when in half screen mode (i.e. after hitting '+' once).
	// Possible values:
	// - 'left': split the window horizontally (side panel on the left, main view on the right)
	// - 'top': split the window vertically (side panel on top, main view below)
	EnlargedSideViewLocation string `yaml:"enlargedSideViewLocation"`
	// One of 'auto' (default) | 'en' | 'zh-CN' | 'zh-TW' | 'pl' | 'nl' | 'ja' | 'ko' | 'ru'
	Language string `yaml:"language" jsonschema:"enum=auto,enum=en,enum=zh-TW,enum=zh-CN,enum=pl,enum=nl,enum=ja,enum=ko,enum=ru"`
	// Format used when displaying time e.g. commit time.
	// Uses Go's time format syntax: https://pkg.go.dev/time#Time.Format
	TimeFormat string `yaml:"timeFormat"`
	// Format used when displaying time if the time is less than 24 hours ago.
	// Uses Go's time format syntax: https://pkg.go.dev/time#Time.Format
	ShortTimeFormat string `yaml:"shortTimeFormat"`
	// Config relating to colors and styles.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#color-attributes
	Theme ThemeConfig `yaml:"theme"`
	// Config relating to the commit length indicator
	CommitLength CommitLengthConfig `yaml:"commitLength"`
	// If true, show the '5 of 20' footer at the bottom of list views
	ShowListFooter bool `yaml:"showListFooter"`
	// If true, display the files in the file views as a tree. If false, display the files as a flat list.
	// This can be toggled from within Lazygit with the '~' key, but that will not change the default.
	ShowFileTree bool `yaml:"showFileTree"`
	// If true, show a random tip in the command log when Lazygit starts
	ShowRandomTip bool `yaml:"showRandomTip"`
	// If true, show the command log
	ShowCommandLog bool `yaml:"showCommandLog"`
	// If true, show the bottom line that contains keybinding info and useful buttons. If false, this line will be hidden except to display a loader for an in-progress action.
	ShowBottomLine bool `yaml:"showBottomLine"`
	// If true, show jump-to-window keybindings in window titles.
	ShowPanelJumps bool `yaml:"showPanelJumps"`
	// Deprecated: use nerdFontsVersion instead
	ShowIcons bool `yaml:"showIcons"`
	// Nerd fonts version to use.
	// One of: '2' | '3' | empty string (default)
	// If empty, do not show icons.
	NerdFontsVersion string `yaml:"nerdFontsVersion" jsonschema:"enum=2,enum=3,enum="`
	// If true (default), file icons are shown in the file views. Only relevant if NerdFontsVersion is not empty.
	ShowFileIcons bool `yaml:"showFileIcons"`
	// If true, show commit hashes alongside branch names in the branches view.
	ShowBranchCommitHash bool `yaml:"showBranchCommitHash"`
	// Height of the command log view
	CommandLogSize int `yaml:"commandLogSize" jsonschema:"minimum=0"`
	// Whether to split the main window when viewing file changes.
	// One of: 'auto' | 'always'
	// If 'auto', only split the main window when a file has both staged and unstaged changes
	SplitDiff string `yaml:"splitDiff" jsonschema:"enum=auto,enum=always"`
	// Default size for focused window. Window size can be changed from within Lazygit with '+' and '_' (but this won't change the default).
	// One of: 'normal' (default) | 'half' | 'full'
	WindowSize string `yaml:"windowSize" jsonschema:"enum=normal,enum=half,enum=full"`
	// Window border style.
	// One of 'rounded' (default) | 'single' | 'double' | 'hidden'
	Border string `yaml:"border" jsonschema:"enum=single,enum=double,enum=rounded,enum=hidden"`
	// If true, show a seriously epic explosion animation when nuking the working tree.
	AnimateExplosion bool `yaml:"animateExplosion"`
	// Whether to stack UI components on top of each other.
	// One of 'auto' (default) | 'always' | 'never'
	PortraitMode string `yaml:"portraitMode"`
	// How things are filtered when typing '/'.
	// One of 'substring' (default) | 'fuzzy'
	FilterMode string `yaml:"filterMode" jsonschema:"enum=substring,enum=fuzzy"`
}

func (c *GuiConfig) UseFuzzySearch() bool {
	return c.FilterMode == "fuzzy"
}

type ThemeConfig struct {
	// Border color of focused window
	ActiveBorderColor []string `yaml:"activeBorderColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Border color of non-focused windows
	InactiveBorderColor []string `yaml:"inactiveBorderColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Border color of focused window when searching in that window
	SearchingActiveBorderColor []string `yaml:"searchingActiveBorderColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Color of keybindings help text in the bottom line
	OptionsTextColor []string `yaml:"optionsTextColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Background color of selected line.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#highlighting-the-selected-line
	SelectedLineBgColor []string `yaml:"selectedLineBgColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Foreground color of copied commit
	CherryPickedCommitFgColor []string `yaml:"cherryPickedCommitFgColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Background color of copied commit
	CherryPickedCommitBgColor []string `yaml:"cherryPickedCommitBgColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Foreground color of marked base commit (for rebase)
	MarkedBaseCommitFgColor []string `yaml:"markedBaseCommitFgColor"`
	// Background color of marked base commit (for rebase)
	MarkedBaseCommitBgColor []string `yaml:"markedBaseCommitBgColor"`
	// Color for file with unstaged changes
	UnstagedChangesColor []string `yaml:"unstagedChangesColor" jsonschema:"minItems=1,uniqueItems=true"`
	// Default text color
	DefaultFgColor []string `yaml:"defaultFgColor" jsonschema:"minItems=1,uniqueItems=true"`
}

type CommitLengthConfig struct {
	// If true, show an indicator of commit message length
	Show bool `yaml:"show"`
}

type GitConfig struct {
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Pagers.md
	Paging PagingConfig `yaml:"paging"`
	// Config relating to committing
	Commit CommitConfig `yaml:"commit"`
	// Config relating to merging
	Merging MergingConfig `yaml:"merging"`
	// list of branches that are considered 'main' branches, used when displaying commits
	MainBranches []string `yaml:"mainBranches" jsonschema:"uniqueItems=true"`
	// Prefix to use when skipping hooks. E.g. if set to 'WIP', then pre-commit hooks will be skipped when the commit message starts with 'WIP'
	SkipHookPrefix string `yaml:"skipHookPrefix"`
	// If true, periodically fetch from remote
	AutoFetch bool `yaml:"autoFetch"`
	// If true, periodically refresh files and submodules
	AutoRefresh bool `yaml:"autoRefresh"`
	// If true, pass the --all arg to git fetch
	FetchAll bool `yaml:"fetchAll"`
	// Command used when displaying the current branch git log in the main window
	BranchLogCmd string `yaml:"branchLogCmd"`
	// Command used to display git log of all branches in the main window
	AllBranchesLogCmd string `yaml:"allBranchesLogCmd"`
	// If true, do not spawn a separate process when using GPG
	OverrideGpg bool `yaml:"overrideGpg"`
	// If true, do not allow force pushes
	DisableForcePushing bool `yaml:"disableForcePushing"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-commit-message-prefix
	CommitPrefixes map[string]CommitPrefixConfig `yaml:"commitPrefixes"`
	// If true, parse emoji strings in commit messages e.g. render :rocket: as 🚀
	// (This should really be under 'gui', not 'git')
	ParseEmoji bool `yaml:"parseEmoji"`
	// Config for showing the log in the commits view
	Log LogConfig `yaml:"log"`
	// When copying commit hashes to the clipboard, truncate them to this
	// length. Set to 40 to disable truncation.
	TruncateCopiedCommitHashesTo int `yaml:"truncateCopiedCommitHashesTo"`
}

type PagerType string

func (PagerType) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.Examples = []any{
		"delta --dark --paging=never",
		"diff-so-fancy",
		"ydiff -p cat -s --wrap --width={{columnWidth}}",
	}
}

type PagingConfig struct {
	// Value of the --color arg in the git diff command. Some pagers want this to be set to 'always' and some want it set to 'never'
	ColorArg string `yaml:"colorArg" jsonschema:"enum=always,enum=never"`
	// e.g.
	// diff-so-fancy
	// delta --dark --paging=never
	// ydiff -p cat -s --wrap --width={{columnWidth}}
	Pager PagerType `yaml:"pager" jsonschema:"minLength=1"`
	// If true, Lazygit will use whatever pager is specified in `$GIT_PAGER`, `$PAGER`, or your *git config*. If the pager ends with something like ` | less` we will strip that part out, because less doesn't play nice with our rendering approach. If the custom pager uses less under the hood, that will also break rendering (hence the `--paging=never` flag for the `delta` pager).
	UseConfig bool `yaml:"useConfig"`
	// e.g. 'difft --color=always'
	ExternalDiffCommand string `yaml:"externalDiffCommand"`
}

type CommitConfig struct {
	// If true, pass '--signoff' flag when committing
	SignOff bool `yaml:"signOff"`
	// Automatic WYSIWYG wrapping of the commit message as you type
	AutoWrapCommitMessage bool `yaml:"autoWrapCommitMessage"`
	// If autoWrapCommitMessage is true, the width to wrap to
	AutoWrapWidth int `yaml:"autoWrapWidth"`
}

type MergingConfig struct {
	// If true, run merges in a subprocess so that if a commit message is required, Lazygit will not hang
	// Only applicable to unix users.
	ManualCommit bool `yaml:"manualCommit"`
	// Extra args passed to `git merge`, e.g. --no-ff
	Args string `yaml:"args" jsonschema:"example=--no-ff"`
}

type LogConfig struct {
	// One of: 'date-order' | 'author-date-order' | 'topo-order' | 'default'
	// 'topo-order' makes it easier to read the git log graph, but commits may not
	// appear chronologically. See https://git-scm.com/docs/
	//
	// Deprecated: Configure this with `Log menu -> Commit sort order` (<c-l> in the commits window by default).
	Order string `yaml:"order" jsonschema:"deprecated,enum=date-order,enum=author-date-order,enum=topo-order,enum=default,deprecated"`
	// This determines whether the git graph is rendered in the commits panel
	// One of 'always' | 'never' | 'when-maximised'
	//
	// Deprecated: Configure this with `Log menu -> Show git graph` (<c-l> in the commits window by default).
	ShowGraph string `yaml:"showGraph" jsonschema:"deprecated,enum=always,enum=never,enum=when-maximised"`
	// displays the whole git graph by default in the commits view (equivalent to passing the `--all` argument to `git log`)
	ShowWholeGraph bool `yaml:"showWholeGraph"`
}

type CommitPrefixConfig struct {
	// pattern to match on. E.g. for 'feature/AB-123' to match on the AB-123 use "^\\w+\\/(\\w+-\\w+).*"
	Pattern string `yaml:"pattern" jsonschema:"example=^\\w+\\/(\\w+-\\w+).*,minLength=1"`
	// Replace directive. E.g. for 'feature/AB-123' to start the commit message with 'AB-123 ' use "[$1] "
	Replace string `yaml:"replace" jsonschema:"example=[$1] ,minLength=1"`
}

type UpdateConfig struct {
	// One of: 'prompt' (default) | 'background' | 'never'
	Method string `yaml:"method" jsonschema:"enum=prompt,enum=background,enum=never"`
	// Period in days between update checks
	Days int64 `yaml:"days" jsonschema:"minimum=0"`
}

type KeybindingConfig struct {
	Universal      KeybindingUniversalConfig      `yaml:"universal"`
	Status         KeybindingStatusConfig         `yaml:"status"`
	Files          KeybindingFilesConfig          `yaml:"files"`
	Branches       KeybindingBranchesConfig       `yaml:"branches"`
	Worktrees      KeybindingWorktreesConfig      `yaml:"worktrees"`
	Commits        KeybindingCommitsConfig        `yaml:"commits"`
	AmendAttribute KeybindingAmendAttributeConfig `yaml:"amendAttribute"`
	Stash          KeybindingStashConfig          `yaml:"stash"`
	CommitFiles    KeybindingCommitFilesConfig    `yaml:"commitFiles"`
	Main           KeybindingMainConfig           `yaml:"main"`
	Submodules     KeybindingSubmodulesConfig     `yaml:"submodules"`
	CommitMessage  KeybindingCommitMessageConfig  `yaml:"commitMessage"`
}

// damn looks like we have some inconsistencies here with -alt and -alt1
type KeybindingUniversalConfig struct {
	Quit                         string   `yaml:"quit"`
	QuitAlt1                     string   `yaml:"quit-alt1"`
	Return                       string   `yaml:"return"`
	QuitWithoutChangingDirectory string   `yaml:"quitWithoutChangingDirectory"`
	TogglePanel                  string   `yaml:"togglePanel"`
	PrevItem                     string   `yaml:"prevItem"`
	NextItem                     string   `yaml:"nextItem"`
	PrevItemAlt                  string   `yaml:"prevItem-alt"`
	NextItemAlt                  string   `yaml:"nextItem-alt"`
	PrevPage                     string   `yaml:"prevPage"`
	NextPage                     string   `yaml:"nextPage"`
	ScrollLeft                   string   `yaml:"scrollLeft"`
	ScrollRight                  string   `yaml:"scrollRight"`
	GotoTop                      string   `yaml:"gotoTop"`
	GotoBottom                   string   `yaml:"gotoBottom"`
	ToggleRangeSelect            string   `yaml:"toggleRangeSelect"`
	RangeSelectDown              string   `yaml:"rangeSelectDown"`
	RangeSelectUp                string   `yaml:"rangeSelectUp"`
	PrevBlock                    string   `yaml:"prevBlock"`
	NextBlock                    string   `yaml:"nextBlock"`
	PrevBlockAlt                 string   `yaml:"prevBlock-alt"`
	NextBlockAlt                 string   `yaml:"nextBlock-alt"`
	NextBlockAlt2                string   `yaml:"nextBlock-alt2"`
	PrevBlockAlt2                string   `yaml:"prevBlock-alt2"`
	JumpToBlock                  []string `yaml:"jumpToBlock"`
	NextMatch                    string   `yaml:"nextMatch"`
	PrevMatch                    string   `yaml:"prevMatch"`
	StartSearch                  string   `yaml:"startSearch"`
	OptionMenu                   string   `yaml:"optionMenu"`
	OptionMenuAlt1               string   `yaml:"optionMenu-alt1"`
	Select                       string   `yaml:"select"`
	GoInto                       string   `yaml:"goInto"`
	Confirm                      string   `yaml:"confirm"`
	ConfirmInEditor              string   `yaml:"confirmInEditor"`
	Remove                       string   `yaml:"remove"`
	New                          string   `yaml:"new"`
	Edit                         string   `yaml:"edit"`
	OpenFile                     string   `yaml:"openFile"`
	ScrollUpMain                 string   `yaml:"scrollUpMain"`
	ScrollDownMain               string   `yaml:"scrollDownMain"`
	ScrollUpMainAlt1             string   `yaml:"scrollUpMain-alt1"`
	ScrollDownMainAlt1           string   `yaml:"scrollDownMain-alt1"`
	ScrollUpMainAlt2             string   `yaml:"scrollUpMain-alt2"`
	ScrollDownMainAlt2           string   `yaml:"scrollDownMain-alt2"`
	ExecuteCustomCommand         string   `yaml:"executeCustomCommand"`
	CreateRebaseOptionsMenu      string   `yaml:"createRebaseOptionsMenu"`
	Push                         string   `yaml:"pushFiles"` // 'Files' appended for legacy reasons
	Pull                         string   `yaml:"pullFiles"` // 'Files' appended for legacy reasons
	Refresh                      string   `yaml:"refresh"`
	CreatePatchOptionsMenu       string   `yaml:"createPatchOptionsMenu"`
	NextTab                      string   `yaml:"nextTab"`
	PrevTab                      string   `yaml:"prevTab"`
	NextScreenMode               string   `yaml:"nextScreenMode"`
	PrevScreenMode               string   `yaml:"prevScreenMode"`
	Undo                         string   `yaml:"undo"`
	Redo                         string   `yaml:"redo"`
	FilteringMenu                string   `yaml:"filteringMenu"`
	DiffingMenu                  string   `yaml:"diffingMenu"`
	DiffingMenuAlt               string   `yaml:"diffingMenu-alt"`
	CopyToClipboard              string   `yaml:"copyToClipboard"`
	OpenRecentRepos              string   `yaml:"openRecentRepos"`
	SubmitEditorText             string   `yaml:"submitEditorText"`
	ExtrasMenu                   string   `yaml:"extrasMenu"`
	ToggleWhitespaceInDiffView   string   `yaml:"toggleWhitespaceInDiffView"`
	IncreaseContextInDiffView    string   `yaml:"increaseContextInDiffView"`
	DecreaseContextInDiffView    string   `yaml:"decreaseContextInDiffView"`
	OpenDiffTool                 string   `yaml:"openDiffTool"`
}

type KeybindingStatusConfig struct {
	CheckForUpdate      string `yaml:"checkForUpdate"`
	RecentRepos         string `yaml:"recentRepos"`
	AllBranchesLogGraph string `yaml:"allBranchesLogGraph"`
}

type KeybindingFilesConfig struct {
	CommitChanges            string `yaml:"commitChanges"`
	CommitChangesWithoutHook string `yaml:"commitChangesWithoutHook"`
	AmendLastCommit          string `yaml:"amendLastCommit"`
	CommitChangesWithEditor  string `yaml:"commitChangesWithEditor"`
	FindBaseCommitForFixup   string `yaml:"findBaseCommitForFixup"`
	ConfirmDiscard           string `yaml:"confirmDiscard"`
	IgnoreFile               string `yaml:"ignoreFile"`
	RefreshFiles             string `yaml:"refreshFiles"`
	StashAllChanges          string `yaml:"stashAllChanges"`
	ViewStashOptions         string `yaml:"viewStashOptions"`
	ToggleStagedAll          string `yaml:"toggleStagedAll"`
	ViewResetOptions         string `yaml:"viewResetOptions"`
	Fetch                    string `yaml:"fetch"`
	ToggleTreeView           string `yaml:"toggleTreeView"`
	OpenMergeTool            string `yaml:"openMergeTool"`
	OpenStatusFilter         string `yaml:"openStatusFilter"`
	CopyFileInfoToClipboard  string `yaml:"copyFileInfoToClipboard"`
}

type KeybindingBranchesConfig struct {
	CreatePullRequest      string `yaml:"createPullRequest"`
	ViewPullRequestOptions string `yaml:"viewPullRequestOptions"`
	CopyPullRequestURL     string `yaml:"copyPullRequestURL"`
	CheckoutBranchByName   string `yaml:"checkoutBranchByName"`
	ForceCheckoutBranch    string `yaml:"forceCheckoutBranch"`
	RebaseBranch           string `yaml:"rebaseBranch"`
	RenameBranch           string `yaml:"renameBranch"`
	MergeIntoCurrentBranch string `yaml:"mergeIntoCurrentBranch"`
	ViewGitFlowOptions     string `yaml:"viewGitFlowOptions"`
	FastForward            string `yaml:"fastForward"`
	CreateTag              string `yaml:"createTag"`
	PushTag                string `yaml:"pushTag"`
	SetUpstream            string `yaml:"setUpstream"`
	FetchRemote            string `yaml:"fetchRemote"`
	SortOrder              string `yaml:"sortOrder"`
}

type KeybindingWorktreesConfig struct {
	ViewWorktreeOptions string `yaml:"viewWorktreeOptions"`
}

type KeybindingCommitsConfig struct {
	SquashDown                     string `yaml:"squashDown"`
	RenameCommit                   string `yaml:"renameCommit"`
	RenameCommitWithEditor         string `yaml:"renameCommitWithEditor"`
	ViewResetOptions               string `yaml:"viewResetOptions"`
	MarkCommitAsFixup              string `yaml:"markCommitAsFixup"`
	CreateFixupCommit              string `yaml:"createFixupCommit"`
	SquashAboveCommits             string `yaml:"squashAboveCommits"`
	MoveDownCommit                 string `yaml:"moveDownCommit"`
	MoveUpCommit                   string `yaml:"moveUpCommit"`
	AmendToCommit                  string `yaml:"amendToCommit"`
	ResetCommitAuthor              string `yaml:"resetCommitAuthor"`
	PickCommit                     string `yaml:"pickCommit"`
	RevertCommit                   string `yaml:"revertCommit"`
	CherryPickCopy                 string `yaml:"cherryPickCopy"`
	PasteCommits                   string `yaml:"pasteCommits"`
	MarkCommitAsBaseForRebase      string `yaml:"markCommitAsBaseForRebase"`
	CreateTag                      string `yaml:"tagCommit"`
	CheckoutCommit                 string `yaml:"checkoutCommit"`
	ResetCherryPick                string `yaml:"resetCherryPick"`
	CopyCommitAttributeToClipboard string `yaml:"copyCommitAttributeToClipboard"`
	OpenLogMenu                    string `yaml:"openLogMenu"`
	OpenInBrowser                  string `yaml:"openInBrowser"`
	ViewBisectOptions              string `yaml:"viewBisectOptions"`
	StartInteractiveRebase         string `yaml:"startInteractiveRebase"`
}

type KeybindingAmendAttributeConfig struct {
	ResetAuthor string `yaml:"resetAuthor"`
	SetAuthor   string `yaml:"setAuthor"`
	AddCoAuthor string `yaml:"addCoAuthor"`
}

type KeybindingStashConfig struct {
	PopStash    string `yaml:"popStash"`
	RenameStash string `yaml:"renameStash"`
}

type KeybindingCommitFilesConfig struct {
	CheckoutCommitFile string `yaml:"checkoutCommitFile"`
}

type KeybindingMainConfig struct {
	ToggleSelectHunk string `yaml:"toggleSelectHunk"`
	PickBothHunks    string `yaml:"pickBothHunks"`
	EditSelectHunk   string `yaml:"editSelectHunk"`
}

type KeybindingSubmodulesConfig struct {
	Init     string `yaml:"init"`
	Update   string `yaml:"update"`
	BulkMenu string `yaml:"bulkMenu"`
}

type KeybindingCommitMessageConfig struct {
	CommitMenu string `yaml:"commitMenu"`
}

// OSConfig contains config on the level of the os
type OSConfig struct {
	// Command for editing a file. Should contain "{{filename}}".
	Edit string `yaml:"edit,omitempty"`

	// Command for editing a file at a given line number. Should contain
	// "{{filename}}", and may optionally contain "{{line}}".
	EditAtLine string `yaml:"editAtLine,omitempty"`

	// Same as EditAtLine, except that the command needs to wait until the
	// window is closed.
	EditAtLineAndWait string `yaml:"editAtLineAndWait,omitempty"`

	// Whether lazygit suspends until an edit process returns
	// Pointer to bool so that we can distinguish unset (nil) from false.
	// We're naming this `editInTerminal` for backwards compatibility
	SuspendOnEdit *bool `yaml:"editInTerminal,omitempty"`

	// For opening a directory in an editor
	OpenDirInEditor string `yaml:"openDirInEditor,omitempty"`

	// A built-in preset that sets all of the above settings. Supported presets
	// are defined in the getPreset function in editor_presets.go.
	EditPreset string `yaml:"editPreset,omitempty" jsonschema:"example=vim,example=nvim,example=emacs,example=nano,example=vscode,example=sublime,example=kakoune,example=helix,example=xcode"`

	// Command for opening a file, as if the file is double-clicked. Should
	// contain "{{filename}}", but doesn't support "{{line}}".
	Open string `yaml:"open,omitempty"`

	// Command for opening a link. Should contain "{{link}}".
	OpenLink string `yaml:"openLink,omitempty"`

	// --------

	// The following configs are all deprecated and kept for backward
	// compatibility. They will be removed in the future.

	// EditCommand is the command for editing a file.
	// Deprecated: use Edit instead. Note that semantics are different:
	// EditCommand is just the command itself, whereas Edit contains a
	// "{{filename}}" variable.
	EditCommand string `yaml:"editCommand,omitempty"`

	// EditCommandTemplate is the command template for editing a file
	// Deprecated: use EditAtLine instead.
	EditCommandTemplate string `yaml:"editCommandTemplate,omitempty"`

	// OpenCommand is the command for opening a file
	// Deprecated: use Open instead.
	OpenCommand string `yaml:"openCommand,omitempty"`

	// OpenLinkCommand is the command for opening a link
	// Deprecated: use OpenLink instead.
	OpenLinkCommand string `yaml:"openLinkCommand,omitempty"`

	// CopyToClipboardCmd is the command for copying to clipboard.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-command-for-copying-to-clipboard
	CopyToClipboardCmd string `yaml:"copyToClipboardCmd,omitempty"`
}

type CustomCommandAfterHook struct {
	CheckForConflicts bool `yaml:"checkForConflicts"`
}

type CustomCommand struct {
	// The key to trigger the command. Use a single letter or one of the values from https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md
	Key string `yaml:"key"`
	// The context in which to listen for the key
	Context string `yaml:"context" jsonschema:"enum=status,enum=files,enum=worktrees,enum=localBranches,enum=remotes,enum=remoteBranches,enum=tags,enum=commits,enum=reflogCommits,enum=subCommits,enum=commitFiles,enum=stash,enum=global"`
	// The command to run (using Go template syntax for placeholder values)
	Command string `yaml:"command" jsonschema:"example=git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"`
	// If true, run the command in a subprocess (e.g. if the command requires user input)
	Subprocess bool `yaml:"subprocess"`
	// A list of prompts that will request user input before running the final command
	Prompts []CustomCommandPrompt `yaml:"prompts"`
	// Text to display while waiting for command to finish
	LoadingText string `yaml:"loadingText" jsonschema:"example=Loading..."`
	// Label for the custom command when displayed in the keybindings menu
	Description string `yaml:"description"`
	// If true, stream the command's output to the Command Log panel
	Stream bool `yaml:"stream"`
	// If true, show the command's output in a popup within Lazygit
	ShowOutput bool `yaml:"showOutput"`
	// Actions to take after the command has completed
	After CustomCommandAfterHook `yaml:"after"`
}

type CustomCommandPrompt struct {
	// One of: 'input' | 'menu' | 'confirm' | 'menuFromCommand'
	Type string `yaml:"type"`
	// Used to reference the entered value from within the custom command. E.g. a prompt with `key: 'Branch'` can be referred to as `{{.Form.Branch}}` in the command
	Key string `yaml:"key"`
	// The title to display in the popup panel
	Title string `yaml:"title"`

	// The initial value to appear in the text box.
	// Only for input prompts.
	InitialValue string `yaml:"initialValue"`
	// Shows suggestions as the input is entered
	// Only for input prompts.
	Suggestions CustomCommandSuggestions `yaml:"suggestions"`

	// The message of the confirmation prompt.
	// Only for confirm prompts.
	Body string `yaml:"body" jsonschema:"example=Are you sure you want to push to the remote?"`

	// Menu options.
	// Only for menu prompts.
	Options []CustomCommandMenuOption `yaml:"options"`

	// The command to run to generate menu options
	// Only for menuFromCommand prompts.
	Command string `yaml:"command" jsonschema:"example=git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"`
	// The regexp to run specifying groups which are going to be kept from the command's output.
	// Only for menuFromCommand prompts.
	Filter string `yaml:"filter" jsonschema:"example=.*{{.SelectedRemote.Name }}/(?P<branch>.*)"`
	// How to format matched groups from the filter to construct a menu item's value.
	// Only for menuFromCommand prompts.
	ValueFormat string `yaml:"valueFormat" jsonschema:"example={{ .branch }}"`
	// Like valueFormat but for the labels. If `labelFormat` is not specified, `valueFormat` is shown instead.
	// Only for menuFromCommand prompts.
	LabelFormat string `yaml:"labelFormat" jsonschema:"example={{ .branch | green }}"`
}

type CustomCommandSuggestions struct {
	// Uses built-in logic to obtain the suggestions. One of 'authors' | 'branches' | 'files' | 'refs' | 'remotes' | 'remoteBranches' | 'tags'
	Preset string `yaml:"preset" jsonschema:"enum=authors,enum=branches,enum=files,enum=refs,enum=remotes,enum=remoteBranches,enum=tags"`
	// Command to run such that each line in the output becomes a suggestion. Mutually exclusive with 'preset' field.
	Command string `yaml:"command" jsonschema:"example=git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"`
}

type CustomCommandMenuOption struct {
	// The first part of the label
	Name string `yaml:"name"`
	// The second part of the label
	Description string `yaml:"description"`
	// The value that will be used in the command
	Value string `yaml:"value" jsonschema:"example=feature,minLength=1"`
}

func GetDefaultConfig() *UserConfig {
	return &UserConfig{
		Gui: GuiConfig{
			ScrollHeight:             2,
			ScrollPastBottom:         true,
			ScrollOffMargin:          2,
			ScrollOffBehavior:        "margin",
			MouseEvents:              true,
			SkipDiscardChangeWarning: false,
			SkipStashWarning:         false,
			SidePanelWidth:           0.3333,
			ExpandFocusedSidePanel:   false,
			MainPanelSplitMode:       "flexible",
			EnlargedSideViewLocation: "left",
			Language:                 "auto",
			TimeFormat:               "02 Jan 06",
			ShortTimeFormat:          time.Kitchen,
			Theme: ThemeConfig{
				ActiveBorderColor:          []string{"green", "bold"},
				SearchingActiveBorderColor: []string{"cyan", "bold"},
				InactiveBorderColor:        []string{"default"},
				OptionsTextColor:           []string{"blue"},
				SelectedLineBgColor:        []string{"blue"},
				CherryPickedCommitBgColor:  []string{"cyan"},
				CherryPickedCommitFgColor:  []string{"blue"},
				MarkedBaseCommitBgColor:    []string{"yellow"},
				MarkedBaseCommitFgColor:    []string{"blue"},
				UnstagedChangesColor:       []string{"red"},
				DefaultFgColor:             []string{"default"},
			},
			CommitLength:              CommitLengthConfig{Show: true},
			SkipNoStagedFilesWarning:  false,
			ShowListFooter:            true,
			ShowCommandLog:            true,
			ShowBottomLine:            true,
			ShowPanelJumps:            true,
			ShowFileTree:              true,
			ShowRandomTip:             true,
			ShowIcons:                 false,
			NerdFontsVersion:          "",
			ShowFileIcons:             true,
			ShowBranchCommitHash:      false,
			CommandLogSize:            8,
			SplitDiff:                 "auto",
			SkipRewordInEditorWarning: false,
			Border:                    "rounded",
			AnimateExplosion:          true,
			PortraitMode:              "auto",
			FilterMode:                "substring",
		},
		Git: GitConfig{
			Paging: PagingConfig{
				ColorArg:            "always",
				Pager:               "",
				UseConfig:           false,
				ExternalDiffCommand: "",
			},
			Commit: CommitConfig{
				SignOff:               false,
				AutoWrapCommitMessage: true,
				AutoWrapWidth:         72,
			},
			Merging: MergingConfig{
				ManualCommit: false,
				Args:         "",
			},
			Log: LogConfig{
				Order:          "topo-order",
				ShowGraph:      "always",
				ShowWholeGraph: false,
			},
			SkipHookPrefix:               "WIP",
			MainBranches:                 []string{"master", "main"},
			AutoFetch:                    true,
			AutoRefresh:                  true,
			FetchAll:                     true,
			BranchLogCmd:                 "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --",
			AllBranchesLogCmd:            "git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium",
			DisableForcePushing:          false,
			CommitPrefixes:               map[string]CommitPrefixConfig(nil),
			ParseEmoji:                   false,
			TruncateCopiedCommitHashesTo: 12,
		},
		Refresher: RefresherConfig{
			RefreshInterval: 10,
			FetchInterval:   60,
		},
		Update: UpdateConfig{
			Method: "prompt",
			Days:   14,
		},
		ConfirmOnQuit:        false,
		QuitOnTopLevelReturn: false,
		Keybinding: KeybindingConfig{
			Universal: KeybindingUniversalConfig{
				Quit:                         "q",
				QuitAlt1:                     "<c-c>",
				Return:                       "<esc>",
				QuitWithoutChangingDirectory: "Q",
				TogglePanel:                  "<tab>",
				PrevItem:                     "<up>",
				NextItem:                     "<down>",
				PrevItemAlt:                  "k",
				NextItemAlt:                  "j",
				PrevPage:                     ",",
				NextPage:                     ".",
				ScrollLeft:                   "H",
				ScrollRight:                  "L",
				GotoTop:                      "<",
				GotoBottom:                   ">",
				ToggleRangeSelect:            "v",
				RangeSelectDown:              "<s-down>",
				RangeSelectUp:                "<s-up>",
				PrevBlock:                    "<left>",
				NextBlock:                    "<right>",
				PrevBlockAlt:                 "h",
				NextBlockAlt:                 "l",
				PrevBlockAlt2:                "<backtab>",
				NextBlockAlt2:                "<tab>",
				JumpToBlock:                  []string{"1", "2", "3", "4", "5"},
				NextMatch:                    "n",
				PrevMatch:                    "N",
				StartSearch:                  "/",
				OptionMenu:                   "<disabled>",
				OptionMenuAlt1:               "?",
				Select:                       "<space>",
				GoInto:                       "<enter>",
				Confirm:                      "<enter>",
				ConfirmInEditor:              "<a-enter>",
				Remove:                       "d",
				New:                          "n",
				Edit:                         "e",
				OpenFile:                     "o",
				OpenRecentRepos:              "<c-r>",
				ScrollUpMain:                 "<pgup>",
				ScrollDownMain:               "<pgdown>",
				ScrollUpMainAlt1:             "K",
				ScrollDownMainAlt1:           "J",
				ScrollUpMainAlt2:             "<c-u>",
				ScrollDownMainAlt2:           "<c-d>",
				ExecuteCustomCommand:         ":",
				CreateRebaseOptionsMenu:      "m",
				Push:                         "P",
				Pull:                         "p",
				Refresh:                      "R",
				CreatePatchOptionsMenu:       "<c-p>",
				NextTab:                      "]",
				PrevTab:                      "[",
				NextScreenMode:               "+",
				PrevScreenMode:               "_",
				Undo:                         "z",
				Redo:                         "<c-z>",
				FilteringMenu:                "<c-s>",
				DiffingMenu:                  "W",
				DiffingMenuAlt:               "<c-e>",
				CopyToClipboard:              "<c-o>",
				SubmitEditorText:             "<enter>",
				ExtrasMenu:                   "@",
				ToggleWhitespaceInDiffView:   "<c-w>",
				IncreaseContextInDiffView:    "}",
				DecreaseContextInDiffView:    "{",
				OpenDiffTool:                 "<c-t>",
			},
			Status: KeybindingStatusConfig{
				CheckForUpdate:      "u",
				RecentRepos:         "<enter>",
				AllBranchesLogGraph: "a",
			},
			Files: KeybindingFilesConfig{
				CommitChanges:            "c",
				CommitChangesWithoutHook: "w",
				AmendLastCommit:          "A",
				CommitChangesWithEditor:  "C",
				FindBaseCommitForFixup:   "<c-f>",
				IgnoreFile:               "i",
				RefreshFiles:             "r",
				StashAllChanges:          "s",
				ViewStashOptions:         "S",
				ToggleStagedAll:          "a",
				ViewResetOptions:         "D",
				Fetch:                    "f",
				ToggleTreeView:           "`",
				OpenMergeTool:            "M",
				OpenStatusFilter:         "<c-b>",
				ConfirmDiscard:           "x",
				CopyFileInfoToClipboard:  "y",
			},
			Branches: KeybindingBranchesConfig{
				CopyPullRequestURL:     "<c-y>",
				CreatePullRequest:      "o",
				ViewPullRequestOptions: "O",
				CheckoutBranchByName:   "c",
				ForceCheckoutBranch:    "F",
				RebaseBranch:           "r",
				RenameBranch:           "R",
				MergeIntoCurrentBranch: "M",
				ViewGitFlowOptions:     "i",
				FastForward:            "f",
				CreateTag:              "T",
				PushTag:                "P",
				SetUpstream:            "u",
				FetchRemote:            "f",
				SortOrder:              "s",
			},
			Worktrees: KeybindingWorktreesConfig{
				ViewWorktreeOptions: "w",
			},
			Commits: KeybindingCommitsConfig{
				SquashDown:                     "s",
				RenameCommit:                   "r",
				RenameCommitWithEditor:         "R",
				ViewResetOptions:               "g",
				MarkCommitAsFixup:              "f",
				CreateFixupCommit:              "F",
				SquashAboveCommits:             "S",
				MoveDownCommit:                 "<c-j>",
				MoveUpCommit:                   "<c-k>",
				AmendToCommit:                  "A",
				ResetCommitAuthor:              "a",
				PickCommit:                     "p",
				RevertCommit:                   "t",
				CherryPickCopy:                 "C",
				PasteCommits:                   "V",
				MarkCommitAsBaseForRebase:      "B",
				CreateTag:                      "T",
				CheckoutCommit:                 "<space>",
				ResetCherryPick:                "<c-R>",
				CopyCommitAttributeToClipboard: "y",
				OpenLogMenu:                    "<c-l>",
				OpenInBrowser:                  "o",
				ViewBisectOptions:              "b",
				StartInteractiveRebase:         "i",
			},
			AmendAttribute: KeybindingAmendAttributeConfig{
				ResetAuthor: "a",
				SetAuthor:   "A",
				AddCoAuthor: "c",
			},
			Stash: KeybindingStashConfig{
				PopStash:    "g",
				RenameStash: "r",
			},
			CommitFiles: KeybindingCommitFilesConfig{
				CheckoutCommitFile: "c",
			},
			Main: KeybindingMainConfig{
				ToggleSelectHunk: "a",
				PickBothHunks:    "b",
				EditSelectHunk:   "E",
			},
			Submodules: KeybindingSubmodulesConfig{
				Init:     "i",
				Update:   "u",
				BulkMenu: "b",
			},
			CommitMessage: KeybindingCommitMessageConfig{
				CommitMenu: "<c-o>",
			},
		},
		OS:                           OSConfig{},
		DisableStartupPopups:         false,
		CustomCommands:               []CustomCommand(nil),
		Services:                     map[string]string(nil),
		NotARepository:               "prompt",
		PromptToReturnFromSubprocess: true,
	}
}
