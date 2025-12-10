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
	// Config relating to things outside of Lazygit like how files are opened, copying to clipboard, etc
	OS OSConfig `yaml:"os,omitempty"`
	// If true, don't display introductory popups upon opening Lazygit.
	DisableStartupPopups bool `yaml:"disableStartupPopups"`
	// User-configured commands that can be invoked from within Lazygit
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Command_Keybindings.md
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
	// Keybindings
	Keybinding KeybindingConfig `yaml:"keybinding"`
}

type RefresherConfig struct {
	// File/submodule refresh interval in seconds.
	// Auto-refresh can be disabled via option 'git.autoRefresh'.
	RefreshInterval int `yaml:"refreshInterval" jsonschema:"minimum=0"`
	// Re-fetch interval in seconds.
	// Auto-fetch can be disabled via option 'git.autoFetch'.
	FetchInterval int `yaml:"fetchInterval" jsonschema:"minimum=0"`
}

func (c *RefresherConfig) RefreshIntervalDuration() time.Duration {
	return time.Second * time.Duration(c.RefreshInterval)
}

func (c *RefresherConfig) FetchIntervalDuration() time.Duration {
	return time.Second * time.Duration(c.FetchInterval)
}

type GuiConfig struct {
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-author-color
	AuthorColors map[string]string `yaml:"authorColors"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-branch-color
	// Deprecated: use branchColorPatterns instead
	BranchColors map[string]string `yaml:"branchColors" jsonschema:"deprecated"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-branch-color
	BranchColorPatterns map[string]string `yaml:"branchColorPatterns"`
	// Custom icons for filenames and file extensions
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-files-icon--color
	CustomIcons CustomIconsConfig `yaml:"customIcons"`
	// The number of lines you scroll by when scrolling the main window
	ScrollHeight int `yaml:"scrollHeight" jsonschema:"minimum=1"`
	// If true, allow scrolling past the bottom of the content in the main window
	ScrollPastBottom bool `yaml:"scrollPastBottom"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#scroll-off-margin
	ScrollOffMargin int `yaml:"scrollOffMargin"`
	// One of: 'margin' (default) | 'jump'
	ScrollOffBehavior string `yaml:"scrollOffBehavior"`
	// The number of spaces per tab; used for everything that's shown in the main view, but probably mostly relevant for diffs.
	// Note that when using a pager, the pager has its own tab width setting, so you need to pass it separately in the pager command.
	TabWidth int `yaml:"tabWidth" jsonschema:"minimum=1"`
	// If true, capture mouse events.
	// When mouse events are captured, it's a little harder to select text: e.g. requiring you to hold the option key when on macOS.
	MouseEvents bool `yaml:"mouseEvents"`
	// If true, do not show a warning when amending a commit.
	SkipAmendWarning bool `yaml:"skipAmendWarning"`
	// If true, do not show a warning when discarding changes in the staging view.
	SkipDiscardChangeWarning bool `yaml:"skipDiscardChangeWarning"`
	// If true, do not show warning when applying/popping the stash
	SkipStashWarning bool `yaml:"skipStashWarning"`
	// If true, do not show a warning when attempting to commit without any staged files; instead stage all unstaged files.
	SkipNoStagedFilesWarning bool `yaml:"skipNoStagedFilesWarning"`
	// If true, do not show a warning when rewording a commit via an external editor
	SkipRewordInEditorWarning bool `yaml:"skipRewordInEditorWarning"`
	// If true, switch to a different worktree without confirmation when checking out a branch that is checked out in that worktree
	SkipSwitchWorktreeOnCheckoutWarning bool `yaml:"skipSwitchWorktreeOnCheckoutWarning"`
	// Fraction of the total screen width to use for the left side section. You may want to pick a small number (e.g. 0.2) if you're using a narrow screen, so that you can see more of the main section.
	// Number from 0 to 1.0.
	SidePanelWidth float64 `yaml:"sidePanelWidth" jsonschema:"maximum=1,minimum=0"`
	// If true, increase the height of the focused side window; creating an accordion effect.
	ExpandFocusedSidePanel bool `yaml:"expandFocusedSidePanel"`
	// The weight of the expanded side panel, relative to the other panels. 2 means twice as tall as the other panels. Only relevant if `expandFocusedSidePanel` is true.
	ExpandedSidePanelWeight int `yaml:"expandedSidePanelWeight"`
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
	// If true, wrap lines in the staging view to the width of the view. This makes it much easier to work with diffs that have long lines, e.g. paragraphs of markdown text.
	WrapLinesInStagingView bool `yaml:"wrapLinesInStagingView"`
	// If true, hunk selection mode will be enabled by default when entering the staging view.
	UseHunkModeInStagingView bool `yaml:"useHunkModeInStagingView"`
	// One of 'auto' (default) | 'en' | 'zh-CN' | 'zh-TW' | 'pl' | 'nl' | 'ja' | 'ko' | 'ru' | 'pt'
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
	// This can be toggled from within Lazygit with the '`' key, but that will not change the default.
	ShowFileTree bool `yaml:"showFileTree"`
	// If true, add a "/" root item in the file tree representing the root of the repository. It is only added when necessary, i.e. when there is more than one item at top level.
	ShowRootItemInFileTree bool `yaml:"showRootItemInFileTree"`
	// If true, show the number of lines changed per file in the Files view
	ShowNumstatInFilesView bool `yaml:"showNumstatInFilesView"`
	// If true, show a random tip in the command log when Lazygit starts
	ShowRandomTip bool `yaml:"showRandomTip"`
	// If true, show the command log
	ShowCommandLog bool `yaml:"showCommandLog"`
	// If true, show the bottom line that contains keybinding info and useful buttons. If false, this line will be hidden except to display a loader for an in-progress action.
	ShowBottomLine bool `yaml:"showBottomLine"`
	// If true, show jump-to-window keybindings in window titles.
	ShowPanelJumps bool `yaml:"showPanelJumps"`
	// Deprecated: use nerdFontsVersion instead
	ShowIcons bool `yaml:"showIcons" jsonschema:"deprecated"`
	// Nerd fonts version to use.
	// One of: '2' | '3' | empty string (default)
	// If empty, do not show icons.
	NerdFontsVersion string `yaml:"nerdFontsVersion" jsonschema:"enum=2,enum=3,enum="`
	// If true (default), file icons are shown in the file views. Only relevant if NerdFontsVersion is not empty.
	ShowFileIcons bool `yaml:"showFileIcons"`
	// Length of author name in (non-expanded) commits view. 2 means show initials only.
	CommitAuthorShortLength int `yaml:"commitAuthorShortLength"`
	// Length of author name in expanded commits view. 2 means show initials only.
	CommitAuthorLongLength int `yaml:"commitAuthorLongLength"`
	// Length of commit hash in commits view. 0 shows '*' if NF icons aren't on.
	CommitHashLength int `yaml:"commitHashLength" jsonschema:"minimum=0"`
	// If true, show commit hashes alongside branch names in the branches view.
	ShowBranchCommitHash bool `yaml:"showBranchCommitHash"`
	// Whether to show the divergence from the base branch in the branches view.
	// One of: 'none' | 'onlyArrow'  | 'arrowAndNumber'
	ShowDivergenceFromBaseBranch string `yaml:"showDivergenceFromBaseBranch" jsonschema:"enum=none,enum=onlyArrow,enum=arrowAndNumber"`
	// Height of the command log view
	CommandLogSize int `yaml:"commandLogSize" jsonschema:"minimum=0"`
	// Whether to split the main window when viewing file changes.
	// One of: 'auto' | 'always'
	// If 'auto', only split the main window when a file has both staged and unstaged changes
	SplitDiff string `yaml:"splitDiff" jsonschema:"enum=auto,enum=always"`
	// Default size for focused window. Can be changed from within Lazygit with '+' and '_' (but this won't change the default).
	// One of: 'normal' (default) | 'half' | 'full'
	ScreenMode string `yaml:"screenMode" jsonschema:"enum=normal,enum=half,enum=full"`
	// Window border style.
	// One of 'rounded' (default) | 'single' | 'double' | 'hidden' | 'bold'
	Border string `yaml:"border" jsonschema:"enum=single,enum=double,enum=rounded,enum=hidden,enum=bold"`
	// If true, show a seriously epic explosion animation when nuking the working tree.
	AnimateExplosion bool `yaml:"animateExplosion"`
	// Whether to stack UI components on top of each other.
	// One of 'auto' (default) | 'always' | 'never'
	PortraitMode string `yaml:"portraitMode"`
	// How things are filtered when typing '/'.
	// One of 'substring' (default) | 'fuzzy'
	FilterMode string `yaml:"filterMode" jsonschema:"enum=substring,enum=fuzzy"`
	// Config relating to the spinner.
	Spinner SpinnerConfig `yaml:"spinner"`
	// Status panel view.
	// One of 'dashboard' (default) | 'allBranchesLog'
	StatusPanelView string `yaml:"statusPanelView" jsonschema:"enum=dashboard,enum=allBranchesLog"`
	// If true, jump to the Files panel after popping a stash
	SwitchToFilesAfterStashPop bool `yaml:"switchToFilesAfterStashPop"`
	// If true, jump to the Files panel after applying a stash
	SwitchToFilesAfterStashApply bool `yaml:"switchToFilesAfterStashApply"`
	// If true, when using the panel jump keys (default 1 through 5) and target panel is already active, go to next tab instead
	SwitchTabsWithPanelJumpKeys bool `yaml:"switchTabsWithPanelJumpKeys"`
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
	// Background color of selected line when view doesn't have focus.
	InactiveViewSelectedLineBgColor []string `yaml:"inactiveViewSelectedLineBgColor" jsonschema:"minItems=1,uniqueItems=true"`
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

type SpinnerConfig struct {
	// The frames of the spinner animation.
	Frames []string `yaml:"frames"`
	// The "speed" of the spinner in milliseconds.
	Rate int `yaml:"rate" jsonschema:"minimum=1"`
}

type GitConfig struct {
	// Array of pagers. Each entry has the following format:
	// [dev] The following documentation is duplicated from the PagingConfig struct below.
	//
	//   # Value of the --color arg in the git diff command. Some pagers want
	//   # this to be set to 'always' and some want it set to 'never'
	//   colorArg: "always"
	//
	//   # e.g.
	//   # diff-so-fancy
	//   # delta --dark --paging=never
	//   # ydiff -p cat -s --wrap --width={{columnWidth}}
	//   pager: ""
	//
	//   # e.g. 'difft --color=always'
	//   externalDiffCommand: ""
	//
	//   # If true, Lazygit will use git's `diff.external` config for paging.
	//   # The advantage over `externalDiffCommand` is that this can be
	//   # configured per file type in .gitattributes; see
	//   # https://git-scm.com/docs/gitattributes#_defining_an_external_diff_driver.
	//   useExternalDiffGitConfig: false
	//
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Pagers.md for more information.
	Pagers []PagingConfig `yaml:"pagers"`
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
	// If not "none", lazygit will automatically fast-forward local branches to match their upstream after fetching. Applies to branches that are not the currently checked out branch, and only to those that are strictly behind their upstream (as opposed to diverged).
	// Possible values: 'none' | 'onlyMainBranches' | 'allBranches'
	AutoForwardBranches string `yaml:"autoForwardBranches" jsonschema:"enum=none,enum=onlyMainBranches,enum=allBranches"`
	// If true, pass the --all arg to git fetch
	FetchAll bool `yaml:"fetchAll"`
	// If true, lazygit will automatically stage files that used to have merge conflicts but no longer do; and it will also ask you if you want to continue a merge or rebase if you've resolved all conflicts. If false, it won't do either of these things.
	AutoStageResolvedConflicts bool `yaml:"autoStageResolvedConflicts"`
	// Command used when displaying the current branch git log in the main window
	BranchLogCmd string `yaml:"branchLogCmd"`
	// Commands used to display git log of all branches in the main window, they will be cycled in order of appearance (array of strings)
	AllBranchesLogCmds []string `yaml:"allBranchesLogCmds"`
	// If true, git diffs are rendered with the `--ignore-all-space` flag, which ignores whitespace changes. Can be toggled from within Lazygit with `<c-w>`.
	IgnoreWhitespaceInDiffView bool `yaml:"ignoreWhitespaceInDiffView"`
	// The number of lines of context to show around each diff hunk. Can be changed from within Lazygit with the `{` and `}` keys.
	DiffContextSize uint64 `yaml:"diffContextSize"`
	// The threshold for considering a file to be renamed, in percent. Can be changed from within Lazygit with the `(` and `)` keys.
	RenameSimilarityThreshold int `yaml:"renameSimilarityThreshold" jsonschema:"minimum=0,maximum=100"`
	// If true, do not spawn a separate process when using GPG
	OverrideGpg bool `yaml:"overrideGpg"`
	// If true, do not allow force pushes
	DisableForcePushing bool `yaml:"disableForcePushing"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-commit-message-prefix
	CommitPrefix []CommitPrefixConfig `yaml:"commitPrefix"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-commit-message-prefix
	CommitPrefixes map[string][]CommitPrefixConfig `yaml:"commitPrefixes"`
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-branch-name-prefix
	BranchPrefix string `yaml:"branchPrefix"`
	// If true, parse emoji strings in commit messages e.g. render :rocket: as 🚀
	// (This should really be under 'gui', not 'git')
	ParseEmoji bool `yaml:"parseEmoji"`
	// Config for showing the log in the commits view
	Log LogConfig `yaml:"log"`
	// How branches are sorted in the local branches view.
	// One of: 'date' (default) | 'recency' | 'alphabetical'
	// Can be changed from within Lazygit with the Sort Order menu (`s`) in the branches panel.
	LocalBranchSortOrder string `yaml:"localBranchSortOrder" jsonschema:"enum=date,enum=recency,enum=alphabetical"`
	// How branches are sorted in the remote branches view.
	// One of: 'date' (default) | 'alphabetical'
	// Can be changed from within Lazygit with the Sort Order menu (`s`) in the remote branches panel.
	RemoteBranchSortOrder string `yaml:"remoteBranchSortOrder" jsonschema:"enum=date,enum=alphabetical"`
	// When copying commit hashes to the clipboard, truncate them to this length. Set to 40 to disable truncation.
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

// [dev] This documentation is duplicated in the GitConfig struct. If you make changes here, make them there too.
type PagingConfig struct {
	// Value of the --color arg in the git diff command. Some pagers want this to be set to 'always' and some want it set to 'never'
	ColorArg string `yaml:"colorArg" jsonschema:"enum=always,enum=never"`
	// e.g.
	// diff-so-fancy
	// delta --dark --paging=never
	// ydiff -p cat -s --wrap --width={{columnWidth}}
	Pager PagerType `yaml:"pager"`
	// e.g. 'difft --color=always'
	ExternalDiffCommand string `yaml:"externalDiffCommand"`
	// If true, Lazygit will use git's `diff.external` config for paging. The advantage over `externalDiffCommand` is that this can be configured per file type in .gitattributes; see https://git-scm.com/docs/gitattributes#_defining_an_external_diff_driver.
	UseExternalDiffGitConfig bool `yaml:"useExternalDiffGitConfig"`
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
	// The commit message to use for a squash merge commit. Can contain "{{selectedRef}}" and "{{currentBranch}}" placeholders.
	SquashMergeMessage string `yaml:"squashMergeMessage"`
}

type LogConfig struct {
	// One of: 'date-order' | 'author-date-order' | 'topo-order' | 'default'
	// 'topo-order' makes it easier to read the git log graph, but commits may not appear chronologically. See https://git-scm.com/docs/
	//
	// Can be changed from within Lazygit with `Log menu -> Commit sort order` (`<c-l>` in the commits window by default).
	Order string `yaml:"order" jsonschema:"enum=date-order,enum=author-date-order,enum=topo-order,enum=default"`
	// This determines whether the git graph is rendered in the commits panel
	// One of 'always' | 'never' | 'when-maximised'
	//
	// Can be toggled from within lazygit with `Log menu -> Show git graph` (`<c-l>` in the commits window by default).
	ShowGraph string `yaml:"showGraph" jsonschema:"enum=always,enum=never,enum=when-maximised"`
	// displays the whole git graph by default in the commits view (equivalent to passing the `--all` argument to `git log`)
	ShowWholeGraph bool `yaml:"showWholeGraph"`

	// Custom format for the commit list side panel.
	// Supports placeholders like:
	//   %h  short hash   %s  subject
	//   %an author name  %as short author
	//   %g  graph line   %d  decorations (tags/branches/markers)
	//   %a  action       %cd date
	// Default: "%h %a %as %g %d %s".
	CustomPaneLogFormat string `yaml:"customPaneLogFormat"`
}

type CommitPrefixConfig struct {
	// pattern to match on. E.g. for 'feature/AB-123' to match on the AB-123 use "^\\w+\\/(\\w+-\\w+).*"
	Pattern string `yaml:"pattern" jsonschema:"example=^\\w+\\/(\\w+-\\w+).*"`
	// Replace directive. E.g. for 'feature/AB-123' to start the commit message with 'AB-123 ' use "[$1] "
	Replace string `yaml:"replace" jsonschema:"example=[$1]"`
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
	Quit                              string   `yaml:"quit"`
	QuitAlt1                          string   `yaml:"quit-alt1"`
	SuspendApp                        string   `yaml:"suspendApp"`
	Return                            string   `yaml:"return"`
	QuitWithoutChangingDirectory      string   `yaml:"quitWithoutChangingDirectory"`
	TogglePanel                       string   `yaml:"togglePanel"`
	PrevItem                          string   `yaml:"prevItem"`
	NextItem                          string   `yaml:"nextItem"`
	PrevItemAlt                       string   `yaml:"prevItem-alt"`
	NextItemAlt                       string   `yaml:"nextItem-alt"`
	PrevPage                          string   `yaml:"prevPage"`
	NextPage                          string   `yaml:"nextPage"`
	ScrollLeft                        string   `yaml:"scrollLeft"`
	ScrollRight                       string   `yaml:"scrollRight"`
	GotoTop                           string   `yaml:"gotoTop"`
	GotoBottom                        string   `yaml:"gotoBottom"`
	GotoTopAlt                        string   `yaml:"gotoTop-alt"`
	GotoBottomAlt                     string   `yaml:"gotoBottom-alt"`
	ToggleRangeSelect                 string   `yaml:"toggleRangeSelect"`
	RangeSelectDown                   string   `yaml:"rangeSelectDown"`
	RangeSelectUp                     string   `yaml:"rangeSelectUp"`
	PrevBlock                         string   `yaml:"prevBlock"`
	NextBlock                         string   `yaml:"nextBlock"`
	PrevBlockAlt                      string   `yaml:"prevBlock-alt"`
	NextBlockAlt                      string   `yaml:"nextBlock-alt"`
	NextBlockAlt2                     string   `yaml:"nextBlock-alt2"`
	PrevBlockAlt2                     string   `yaml:"prevBlock-alt2"`
	JumpToBlock                       []string `yaml:"jumpToBlock"`
	FocusMainView                     string   `yaml:"focusMainView"`
	NextMatch                         string   `yaml:"nextMatch"`
	PrevMatch                         string   `yaml:"prevMatch"`
	StartSearch                       string   `yaml:"startSearch"`
	OptionMenu                        string   `yaml:"optionMenu"`
	OptionMenuAlt1                    string   `yaml:"optionMenu-alt1"`
	Select                            string   `yaml:"select"`
	GoInto                            string   `yaml:"goInto"`
	Confirm                           string   `yaml:"confirm"`
	ConfirmMenu                       string   `yaml:"confirmMenu"`
	ConfirmSuggestion                 string   `yaml:"confirmSuggestion"`
	ConfirmInEditor                   string   `yaml:"confirmInEditor"`
	ConfirmInEditorAlt                string   `yaml:"confirmInEditor-alt"`
	Remove                            string   `yaml:"remove"`
	New                               string   `yaml:"new"`
	Edit                              string   `yaml:"edit"`
	OpenFile                          string   `yaml:"openFile"`
	ScrollUpMain                      string   `yaml:"scrollUpMain"`
	ScrollDownMain                    string   `yaml:"scrollDownMain"`
	ScrollUpMainAlt1                  string   `yaml:"scrollUpMain-alt1"`
	ScrollDownMainAlt1                string   `yaml:"scrollDownMain-alt1"`
	ScrollUpMainAlt2                  string   `yaml:"scrollUpMain-alt2"`
	ScrollDownMainAlt2                string   `yaml:"scrollDownMain-alt2"`
	ExecuteShellCommand               string   `yaml:"executeShellCommand"`
	CreateRebaseOptionsMenu           string   `yaml:"createRebaseOptionsMenu"`
	Push                              string   `yaml:"pushFiles"` // 'Files' appended for legacy reasons
	Pull                              string   `yaml:"pullFiles"` // 'Files' appended for legacy reasons
	Refresh                           string   `yaml:"refresh"`
	CreatePatchOptionsMenu            string   `yaml:"createPatchOptionsMenu"`
	NextTab                           string   `yaml:"nextTab"`
	PrevTab                           string   `yaml:"prevTab"`
	NextScreenMode                    string   `yaml:"nextScreenMode"`
	PrevScreenMode                    string   `yaml:"prevScreenMode"`
	CyclePagers                       string   `yaml:"cyclePagers"`
	Undo                              string   `yaml:"undo"`
	Redo                              string   `yaml:"redo"`
	FilteringMenu                     string   `yaml:"filteringMenu"`
	DiffingMenu                       string   `yaml:"diffingMenu"`
	DiffingMenuAlt                    string   `yaml:"diffingMenu-alt"`
	CopyToClipboard                   string   `yaml:"copyToClipboard"`
	OpenRecentRepos                   string   `yaml:"openRecentRepos"`
	SubmitEditorText                  string   `yaml:"submitEditorText"`
	ExtrasMenu                        string   `yaml:"extrasMenu"`
	ToggleWhitespaceInDiffView        string   `yaml:"toggleWhitespaceInDiffView"`
	IncreaseContextInDiffView         string   `yaml:"increaseContextInDiffView"`
	DecreaseContextInDiffView         string   `yaml:"decreaseContextInDiffView"`
	IncreaseRenameSimilarityThreshold string   `yaml:"increaseRenameSimilarityThreshold"`
	DecreaseRenameSimilarityThreshold string   `yaml:"decreaseRenameSimilarityThreshold"`
	OpenDiffTool                      string   `yaml:"openDiffTool"`
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
	OpenMergeOptions         string `yaml:"openMergeOptions"`
	OpenStatusFilter         string `yaml:"openStatusFilter"`
	CopyFileInfoToClipboard  string `yaml:"copyFileInfoToClipboard"`
	CollapseAll              string `yaml:"collapseAll"`
	ExpandAll                string `yaml:"expandAll"`
}

type KeybindingBranchesConfig struct {
	CreatePullRequest      string `yaml:"createPullRequest"`
	ViewPullRequestOptions string `yaml:"viewPullRequestOptions"`
	CopyPullRequestURL     string `yaml:"copyPullRequestURL"`
	CheckoutBranchByName   string `yaml:"checkoutBranchByName"`
	ForceCheckoutBranch    string `yaml:"forceCheckoutBranch"`
	CheckoutPreviousBranch string `yaml:"checkoutPreviousBranch"`
	RebaseBranch           string `yaml:"rebaseBranch"`
	RenameBranch           string `yaml:"renameBranch"`
	MergeIntoCurrentBranch string `yaml:"mergeIntoCurrentBranch"`
	MoveCommitsToNewBranch string `yaml:"moveCommitsToNewBranch"`
	ViewGitFlowOptions     string `yaml:"viewGitFlowOptions"`
	FastForward            string `yaml:"fastForward"`
	CreateTag              string `yaml:"createTag"`
	PushTag                string `yaml:"pushTag"`
	SetUpstream            string `yaml:"setUpstream"`
	FetchRemote            string `yaml:"fetchRemote"`
	AddForkRemote          string `yaml:"addForkRemote"`
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
	SelectCommitsOfCurrentBranch   string `yaml:"selectCommitsOfCurrentBranch"`
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

	// Command for editing a file at a given line number. Should contain "{{filename}}", and may optionally contain "{{line}}".
	EditAtLine string `yaml:"editAtLine,omitempty"`

	// Same as EditAtLine, except that the command needs to wait until the window is closed.
	EditAtLineAndWait string `yaml:"editAtLineAndWait,omitempty"`

	// Whether lazygit suspends until an edit process returns
	// [dev] Pointer to bool so that we can distinguish unset (nil) from false.
	// [dev] We're naming this `editInTerminal` for backwards compatibility
	SuspendOnEdit *bool `yaml:"editInTerminal,omitempty"`

	// For opening a directory in an editor
	OpenDirInEditor string `yaml:"openDirInEditor,omitempty"`

	// A built-in preset that sets all of the above settings. Supported presets are defined in the getPreset function in editor_presets.go.
	EditPreset string `yaml:"editPreset,omitempty" jsonschema:"example=vim,example=nvim,example=emacs,example=nano,example=vscode,example=sublime,example=kakoune,example=helix,example=xcode,example=zed,example=acme"`

	// Command for opening a file, as if the file is double-clicked. Should contain "{{filename}}", but doesn't support "{{line}}".
	Open string `yaml:"open,omitempty"`

	// Command for opening a link. Should contain "{{link}}".
	OpenLink string `yaml:"openLink,omitempty"`

	// CopyToClipboardCmd is the command for copying to clipboard.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-command-for-copying-to-and-pasting-from-clipboard
	CopyToClipboardCmd string `yaml:"copyToClipboardCmd,omitempty"`

	// ReadFromClipboardCmd is the command for reading the clipboard.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-command-for-copying-to-and-pasting-from-clipboard
	ReadFromClipboardCmd string `yaml:"readFromClipboardCmd,omitempty"`

	// A shell startup file containing shell aliases or shell functions. This will be sourced before running any shell commands, so that shell functions are available in the `:` command prompt or even in custom commands.
	// See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#using-aliases-or-functions-in-shell-commands
	ShellFunctionsFile string `yaml:"shellFunctionsFile"`
}

type CustomCommandAfterHook struct {
	CheckForConflicts bool `yaml:"checkForConflicts"`
}

type CustomCommand struct {
	// The key to trigger the command. Use a single letter or one of the values from https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md
	Key string `yaml:"key"`
	// Instead of defining a single custom command, create a menu of custom commands. Useful for grouping related commands together under a single keybinding, and for keeping them out of the global keybindings menu.
	// When using this, all other fields except Key and Description are ignored and must be empty.
	CommandMenu []CustomCommand `yaml:"commandMenu"`
	// The context in which to listen for the key. Valid values are: status, files, worktrees, localBranches, remotes, remoteBranches, tags, commits, reflogCommits, subCommits, commitFiles, stash, and global. Multiple contexts separated by comma are allowed; most useful for "commits, subCommits" or "files, commitFiles".
	Context string `yaml:"context" jsonschema:"example=status,example=files,example=worktrees,example=localBranches,example=remotes,example=remoteBranches,example=tags,example=commits,example=reflogCommits,example=subCommits,example=commitFiles,example=stash,example=global"`
	// The command to run (using Go template syntax for placeholder values)
	Command string `yaml:"command" jsonschema:"example=git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"`
	// A list of prompts that will request user input before running the final command
	Prompts []CustomCommandPrompt `yaml:"prompts"`
	// Text to display while waiting for command to finish
	LoadingText string `yaml:"loadingText" jsonschema:"example=Loading..."`
	// Label for the custom command when displayed in the keybindings menu
	Description string `yaml:"description"`
	// Where the output of the command should go. 'none' discards it, 'terminal' suspends lazygit and runs the command in the terminal (useful for commands that require user input), 'log' streams it to the command log, 'logWithPty' is like 'log' but runs the command in a pseudo terminal (can be useful for commands that produce colored output when the output is a terminal), and 'popup' shows it in a popup.
	Output string `yaml:"output" jsonschema:"enum=none,enum=terminal,enum=log,enum=logWithPty,enum=popup"`
	// The title to display in the popup panel if output is set to 'popup'. If left unset, the command will be used as the title.
	OutputTitle string `yaml:"outputTitle"`
	// Actions to take after the command has completed
	// [dev] Pointer so that we can tell whether it appears in the config file
	After *CustomCommandAfterHook `yaml:"after"`
}

func (c *CustomCommand) GetDescription() string {
	if c.Description != "" {
		return c.Description
	}

	return c.Command
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

type CustomIconsConfig struct {
	// Map of filenames to icon properties (icon and color)
	Filenames map[string]IconProperties `yaml:"filenames"`
	// Map of file extensions (including the dot) to icon properties (icon and color)
	Extensions map[string]IconProperties `yaml:"extensions"`
}

type IconProperties struct {
	Icon  string `yaml:"icon"`
	Color string `yaml:"color"`
}

func GetDefaultConfig() *UserConfig {
	return &UserConfig{
		Gui: GuiConfig{
			ScrollHeight:             2,
			ScrollPastBottom:         true,
			ScrollOffMargin:          2,
			ScrollOffBehavior:        "margin",
			TabWidth:                 4,
			MouseEvents:              true,
			SkipAmendWarning:         false,
			SkipDiscardChangeWarning: false,
			SkipStashWarning:         false,
			SidePanelWidth:           0.3333,
			ExpandFocusedSidePanel:   false,
			ExpandedSidePanelWeight:  2,
			MainPanelSplitMode:       "flexible",
			EnlargedSideViewLocation: "left",
			WrapLinesInStagingView:   true,
			UseHunkModeInStagingView: true,
			Language:                 "auto",
			TimeFormat:               "02 Jan 06",
			ShortTimeFormat:          time.Kitchen,
			Theme: ThemeConfig{
				ActiveBorderColor:               []string{"green", "bold"},
				SearchingActiveBorderColor:      []string{"cyan", "bold"},
				InactiveBorderColor:             []string{"default"},
				OptionsTextColor:                []string{"blue"},
				SelectedLineBgColor:             []string{"blue"},
				InactiveViewSelectedLineBgColor: []string{"bold"},
				CherryPickedCommitBgColor:       []string{"cyan"},
				CherryPickedCommitFgColor:       []string{"blue"},
				MarkedBaseCommitBgColor:         []string{"yellow"},
				MarkedBaseCommitFgColor:         []string{"blue"},
				UnstagedChangesColor:            []string{"red"},
				DefaultFgColor:                  []string{"default"},
			},
			CommitLength:                        CommitLengthConfig{Show: true},
			SkipNoStagedFilesWarning:            false,
			ShowListFooter:                      true,
			ShowCommandLog:                      true,
			ShowBottomLine:                      true,
			ShowPanelJumps:                      true,
			ShowFileTree:                        true,
			ShowRootItemInFileTree:              true,
			ShowNumstatInFilesView:              false,
			ShowRandomTip:                       true,
			ShowIcons:                           false,
			NerdFontsVersion:                    "",
			ShowFileIcons:                       true,
			CommitAuthorShortLength:             2,
			CommitAuthorLongLength:              17,
			CommitHashLength:                    8,
			ShowBranchCommitHash:                false,
			ShowDivergenceFromBaseBranch:        "none",
			CommandLogSize:                      8,
			SplitDiff:                           "auto",
			SkipRewordInEditorWarning:           false,
			SkipSwitchWorktreeOnCheckoutWarning: false,
			ScreenMode:                          "normal",
			Border:                              "rounded",
			AnimateExplosion:                    true,
			PortraitMode:                        "auto",
			FilterMode:                          "substring",
			Spinner: SpinnerConfig{
				Frames: []string{"|", "/", "-", "\\"},
				Rate:   50,
			},
			StatusPanelView:              "dashboard",
			SwitchToFilesAfterStashPop:   true,
			SwitchToFilesAfterStashApply: true,
			SwitchTabsWithPanelJumpKeys:  false,
		},
		Git: GitConfig{
			Commit: CommitConfig{
				SignOff:               false,
				AutoWrapCommitMessage: true,
				AutoWrapWidth:         72,
			},
			Merging: MergingConfig{
				ManualCommit:       false,
				Args:               "",
				SquashMergeMessage: "Squash merge {{selectedRef}} into {{currentBranch}}",
			},
			Log: LogConfig{
				Order:               "topo-order",
				ShowGraph:           "always",
				ShowWholeGraph:      false,
				CustomPaneLogFormat: "%h %a %as %g %d %s",
			},
			LocalBranchSortOrder:         "date",
			RemoteBranchSortOrder:        "date",
			SkipHookPrefix:               "WIP",
			MainBranches:                 []string{"master", "main"},
			AutoFetch:                    true,
			AutoRefresh:                  true,
			AutoForwardBranches:          "onlyMainBranches",
			FetchAll:                     true,
			AutoStageResolvedConflicts:   true,
			BranchLogCmd:                 "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --",
			AllBranchesLogCmds:           []string{"git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium"},
			IgnoreWhitespaceInDiffView:   false,
			DiffContextSize:              3,
			RenameSimilarityThreshold:    50,
			DisableForcePushing:          false,
			CommitPrefixes:               map[string][]CommitPrefixConfig(nil),
			BranchPrefix:                 "",
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
		ConfirmOnQuit:                false,
		QuitOnTopLevelReturn:         false,
		OS:                           OSConfig{},
		DisableStartupPopups:         false,
		CustomCommands:               []CustomCommand(nil),
		Services:                     map[string]string(nil),
		NotARepository:               "prompt",
		PromptToReturnFromSubprocess: true,
		Keybinding: KeybindingConfig{
			Universal: KeybindingUniversalConfig{
				Quit:                              "q",
				QuitAlt1:                          "<c-c>",
				SuspendApp:                        "<c-z>",
				Return:                            "<esc>",
				QuitWithoutChangingDirectory:      "Q",
				TogglePanel:                       "<tab>",
				PrevItem:                          "<up>",
				NextItem:                          "<down>",
				PrevItemAlt:                       "k",
				NextItemAlt:                       "j",
				PrevPage:                          ",",
				NextPage:                          ".",
				ScrollLeft:                        "H",
				ScrollRight:                       "L",
				GotoTop:                           "<",
				GotoBottom:                        ">",
				GotoTopAlt:                        "<home>",
				GotoBottomAlt:                     "<end>",
				ToggleRangeSelect:                 "v",
				RangeSelectDown:                   "<s-down>",
				RangeSelectUp:                     "<s-up>",
				PrevBlock:                         "<left>",
				NextBlock:                         "<right>",
				PrevBlockAlt:                      "h",
				NextBlockAlt:                      "l",
				PrevBlockAlt2:                     "<backtab>",
				NextBlockAlt2:                     "<tab>",
				JumpToBlock:                       []string{"1", "2", "3", "4", "5"},
				FocusMainView:                     "0",
				NextMatch:                         "n",
				PrevMatch:                         "N",
				StartSearch:                       "/",
				OptionMenu:                        "<disabled>",
				OptionMenuAlt1:                    "?",
				Select:                            "<space>",
				GoInto:                            "<enter>",
				Confirm:                           "<enter>",
				ConfirmMenu:                       "<enter>",
				ConfirmSuggestion:                 "<enter>",
				ConfirmInEditor:                   "<a-enter>",
				ConfirmInEditorAlt:                "<c-s>",
				Remove:                            "d",
				New:                               "n",
				Edit:                              "e",
				OpenFile:                          "o",
				OpenRecentRepos:                   "<c-r>",
				ScrollUpMain:                      "<pgup>",
				ScrollDownMain:                    "<pgdown>",
				ScrollUpMainAlt1:                  "K",
				ScrollDownMainAlt1:                "J",
				ScrollUpMainAlt2:                  "<c-u>",
				ScrollDownMainAlt2:                "<c-d>",
				ExecuteShellCommand:               ":",
				CreateRebaseOptionsMenu:           "m",
				Push:                              "P",
				Pull:                              "p",
				Refresh:                           "R",
				CreatePatchOptionsMenu:            "<c-p>",
				NextTab:                           "]",
				PrevTab:                           "[",
				NextScreenMode:                    "+",
				PrevScreenMode:                    "_",
				CyclePagers:                       "|",
				Undo:                              "z",
				Redo:                              "Z",
				FilteringMenu:                     "<c-s>",
				DiffingMenu:                       "W",
				DiffingMenuAlt:                    "<c-e>",
				CopyToClipboard:                   "<c-o>",
				SubmitEditorText:                  "<enter>",
				ExtrasMenu:                        "@",
				ToggleWhitespaceInDiffView:        "<c-w>",
				IncreaseContextInDiffView:         "}",
				DecreaseContextInDiffView:         "{",
				IncreaseRenameSimilarityThreshold: ")",
				DecreaseRenameSimilarityThreshold: "(",
				OpenDiffTool:                      "<c-t>",
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
				OpenMergeOptions:         "M",
				OpenStatusFilter:         "<c-b>",
				ConfirmDiscard:           "x",
				CopyFileInfoToClipboard:  "y",
				CollapseAll:              "-",
				ExpandAll:                "=",
			},
			Branches: KeybindingBranchesConfig{
				CopyPullRequestURL:     "<c-y>",
				CreatePullRequest:      "o",
				ViewPullRequestOptions: "O",
				CheckoutBranchByName:   "c",
				ForceCheckoutBranch:    "F",
				CheckoutPreviousBranch: "-",
				RebaseBranch:           "r",
				RenameBranch:           "R",
				MergeIntoCurrentBranch: "M",
				MoveCommitsToNewBranch: "N",
				ViewGitFlowOptions:     "i",
				FastForward:            "f",
				CreateTag:              "T",
				PushTag:                "P",
				SetUpstream:            "u",
				FetchRemote:            "f",
				AddForkRemote:          "F",
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
				SelectCommitsOfCurrentBranch:   "*",
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
	}
}
