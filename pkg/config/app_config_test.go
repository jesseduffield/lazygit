package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationOfRenamedKeys(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No rename needed",
			input: `foo:
  bar: 5
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Rename one",
			input: `gui:
  skipUnstageLineWarning: true
`,
			expected: `gui:
  skipDiscardChangeWarning: true
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Renamed 'gui.skipUnstageLineWarning' to 'skipDiscardChangeWarning'"},
		},
		{
			name: "Rename several",
			input: `gui:
  windowSize: half
  skipUnstageLineWarning: true
keybinding:
  universal:
    executeCustomCommand: a
`,
			expected: `gui:
  screenMode: half
  skipDiscardChangeWarning: true
keybinding:
  universal:
    executeShellCommand: a
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed 'gui.skipUnstageLineWarning' to 'skipDiscardChangeWarning'",
				"Renamed 'keybinding.universal.executeCustomCommand' to 'executeShellCommand'",
				"Renamed 'gui.windowSize' to 'screenMode'",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestMigrateNullKeybindingsToDisabled(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No change needed",
			input: `keybinding:
  universal:
    quit: q
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Change one",
			input: `keybinding:
  universal:
    quit: null
`,
			expected: `keybinding:
  universal:
    quit: <disabled>
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.quit'"},
		},
		{
			name: "Change several",
			input: `keybinding:
  universal:
    quit: null
    return: <esc>
    new: null
`,
			expected: `keybinding:
  universal:
    quit: <disabled>
    return: <esc>
    new: <disabled>
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.quit'",
				"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.new'",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestCommitPrefixMigrations(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Single CommitPrefix Rename",
			input: `git:
  commitPrefix:
     pattern: "^\\w+-\\w+.*"
     replace: '[JIRA $0] '
`,
			expected: `git:
  commitPrefix:
    - pattern: "^\\w+-\\w+.*"
      replace: '[JIRA $0] '
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'git.commitPrefix' to an array of strings"},
		},
		{
			name: "Complicated CommitPrefixes Rename",
			input: `git:
  commitPrefixes:
    foo:
      pattern: "^\\w+-\\w+.*"
      replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      pattern: "^foo.bar*"
      replace: '[FUN $0] '
`,
			expected: `git:
  commitPrefixes:
    foo:
      - pattern: "^\\w+-\\w+.*"
        replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      - pattern: "^foo.bar*"
        replace: '[FUN $0] '
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'git.commitPrefixes' elements to arrays of strings"},
		},
		{
			name:              "Incomplete Configuration",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No changes made when already migrated",
			input: `
git:
   commitPrefix:
    - pattern: "Hello World"
      replace: "Goodbye"
   commitPrefixes:
    foo:
      - pattern: "^\\w+-\\w+.*"
        replace: '[JIRA $0] '`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestCustomCommandsOutputMigration(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Convert subprocess to output=terminal",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: terminal
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'subprocess: true' to 'output: terminal' in custom command"},
		},
		{
			name: "Convert stream to output=log",
			input: `customCommands:
  - command: echo 'hello'
    stream: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: log
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'stream: true' to 'output: log' in custom command"},
		},
		{
			name: "Convert showOutput to output=popup",
			input: `customCommands:
  - command: echo 'hello'
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: popup
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'showOutput: true' to 'output: popup' in custom command"},
		},
		{
			name: "Subprocess wins over the other two",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: true
    stream: true
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: terminal
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'subprocess: true' to 'output: terminal' in custom command",
				"Deleted redundant 'stream: true' property in custom command",
				"Deleted redundant 'showOutput: true' property in custom command",
			},
		},
		{
			name: "Stream wins over showOutput",
			input: `customCommands:
  - command: echo 'hello'
    stream: true
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: log
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'stream: true' to 'output: log' in custom command",
				"Deleted redundant 'showOutput: true' property in custom command",
			},
		},
		{
			name: "Explicitly setting to false doesn't create an output=none key",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: false
    stream: false
    showOutput: false
  `,
			expected: `customCommands:
  - command: echo 'hello'
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Deleted redundant 'subprocess: false' in custom command",
				"Deleted redundant 'stream: false' property in custom command",
				"Deleted redundant 'showOutput: false' property in custom command",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

var largeConfiguration = []byte(`
# Config relating to the Lazygit UI
gui:
  # The number of lines you scroll by when scrolling the main window
  scrollHeight: 2

  # If true, allow scrolling past the bottom of the content in the main window
  scrollPastBottom: true

  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#scroll-off-margin
  scrollOffMargin: 2

  # One of: 'margin' (default) | 'jump'
  scrollOffBehavior: margin

  # The number of spaces per tab; used for everything that's shown in the main view, but probably mostly relevant for diffs.
  # Note that when using a pager, the pager has its own tab width setting, so you need to pass it separately in the pager command.
  tabWidth: 4

  # If true, capture mouse events.
  # When mouse events are captured, it's a little harder to select text: e.g. requiring you to hold the option key when on macOS.
  mouseEvents: true

  # If true, do not show a warning when amending a commit.
  skipAmendWarning: false

  # If true, do not show a warning when discarding changes in the staging view.
  skipDiscardChangeWarning: false

  # If true, do not show warning when applying/popping the stash
  skipStashWarning: false

  # If true, do not show a warning when attempting to commit without any staged files; instead stage all unstaged files.
  skipNoStagedFilesWarning: false

  # If true, do not show a warning when rewording a commit via an external editor
  skipRewordInEditorWarning: false

  # Fraction of the total screen width to use for the left side section. You may want to pick a small number (e.g. 0.2) if you're using a narrow screen, so that you can see more of the main section.
  # Number from 0 to 1.0.
  sidePanelWidth: 0.3333

  # If true, increase the height of the focused side window; creating an accordion effect.
  expandFocusedSidePanel: false

  # The weight of the expanded side panel, relative to the other panels. 2 means
  # twice as tall as the other panels. Only relevant if expandFocusedSidePanel is true.
  expandedSidePanelWeight: 2

  # Sometimes the main window is split in two (e.g. when the selected file has both staged and unstaged changes). This setting controls how the two sections are split.
  # Options are:
  # - 'horizontal': split the window horizontally
  # - 'vertical': split the window vertically
  # - 'flexible': (default) split the window horizontally if the window is wide enough, otherwise split vertically
  mainPanelSplitMode: flexible

  # How the window is split when in half screen mode (i.e. after hitting '+' once).
  # Possible values:
  # - 'left': split the window horizontally (side panel on the left, main view on the right)
  # - 'top': split the window vertically (side panel on top, main view below)
  enlargedSideViewLocation: left

  # If true, wrap lines in the staging view to the width of the view. This
  # makes it much easier to work with diffs that have long lines, e.g.
  # paragraphs of markdown text.
  wrapLinesInStagingView: true

  # One of 'auto' (default) | 'en' | 'zh-CN' | 'zh-TW' | 'pl' | 'nl' | 'ja' | 'ko' | 'ru'
  language: auto

  # Format used when displaying time e.g. commit time.
  # Uses Go's time format syntax: https://pkg.go.dev/time#Time.Format
  timeFormat: 02 Jan 06

  # Format used when displaying time if the time is less than 24 hours ago.
  # Uses Go's time format syntax: https://pkg.go.dev/time#Time.Format
  shortTimeFormat: 3:04PM

  # Config relating to colors and styles.
  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#color-attributes
  theme:
    # Border color of focused window
    activeBorderColor:
      - green
      - bold

    # Border color of non-focused windows
    inactiveBorderColor:
      - default

    # Border color of focused window when searching in that window
    searchingActiveBorderColor:
      - cyan
      - bold

    # Color of keybindings help text in the bottom line
    optionsTextColor:
      - blue

    # Background color of selected line.
    # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#highlighting-the-selected-line
    selectedLineBgColor:
      - blue

    # Background color of selected line when view doesn't have focus.
    inactiveViewSelectedLineBgColor:
      - bold

    # Foreground color of copied commit
    cherryPickedCommitFgColor:
      - blue

    # Background color of copied commit
    cherryPickedCommitBgColor:
      - cyan

    # Foreground color of marked base commit (for rebase)
    markedBaseCommitFgColor:
      - blue

    # Background color of marked base commit (for rebase)
    markedBaseCommitBgColor:
      - yellow

    # Color for file with unstaged changes
    unstagedChangesColor:
      - red

    # Default text color
    defaultFgColor:
      - default

  # Config relating to the commit length indicator
  commitLength:
    # If true, show an indicator of commit message length
    show: true

  # If true, show the '5 of 20' footer at the bottom of list views
  showListFooter: true

  # If true, display the files in the file views as a tree. If false, display the files as a flat list.
  # This can be toggled from within Lazygit with the '' key, but that will not change the default.
  showFileTree: true

  # If true, show the number of lines changed per file in the Files view
  showNumstatInFilesView: false

  # If true, show a random tip in the command log when Lazygit starts
  showRandomTip: true

  # If true, show the command log
  showCommandLog: true

  # If true, show the bottom line that contains keybinding info and useful buttons. If false, this line will be hidden except to display a loader for an in-progress action.
  showBottomLine: true

  # If true, show jump-to-window keybindings in window titles.
  showPanelJumps: true

  # Deprecated: use nerdFontsVersion instead
  showIcons: false

  # Nerd fonts version to use.
  # One of: '2' | '3' | empty string (default)
  # If empty, do not show icons.
  nerdFontsVersion: ""

  # If true (default), file icons are shown in the file views. Only relevant if NerdFontsVersion is not empty.
  showFileIcons: true

  # Length of author name in (non-expanded) commits view. 2 means show initials only.
  commitAuthorShortLength: 2

  # Length of author name in expanded commits view. 2 means show initials only.
  commitAuthorLongLength: 17

  # Length of commit hash in commits view. 0 shows '*' if NF icons aren't on.
  commitHashLength: 8

  # If true, show commit hashes alongside branch names in the branches view.
  showBranchCommitHash: false

  # Whether to show the divergence from the base branch in the branches view.
  # One of: 'none' | 'onlyArrow'  | 'arrowAndNumber'
  showDivergenceFromBaseBranch: none

  # Height of the command log view
  commandLogSize: 8

  # Whether to split the main window when viewing file changes.
  # One of: 'auto' | 'always'
  # If 'auto', only split the main window when a file has both staged and unstaged changes
  splitDiff: auto

  # Default size for focused window. Can be changed from within Lazygit with '+' and '_' (but this won't change the default).
  # One of: 'normal' (default) | 'half' | 'full'
  screenMode: normal

  # Window border style.
  # One of 'rounded' (default) | 'single' | 'double' | 'hidden' | 'bold'
  border: rounded

  # If true, show a seriously epic explosion animation when nuking the working tree.
  animateExplosion: true

  # Whether to stack UI components on top of each other.
  # One of 'auto' (default) | 'always' | 'never'
  portraitMode: auto

  # How things are filtered when typing '/'.
  # One of 'substring' (default) | 'fuzzy'
  filterMode: substring

  # Config relating to the spinner.
  spinner:
    # The frames of the spinner animation.
    frames:
      - '|'
      - /
      - '-'
      - \

    # The "speed" of the spinner in milliseconds.
    rate: 50

  # Status panel view.
  # One of 'dashboard' (default) | 'allBranchesLog'
  statusPanelView: dashboard

  # If true, jump to the Files panel after popping a stash
  switchToFilesAfterStashPop: true

  # If true, jump to the Files panel after applying a stash
  switchToFilesAfterStashApply: true

  # If true, when using the panel jump keys (default 1 through 5) and target panel is already active, go to next tab instead
  switchTabsWithPanelJumpKeys: false

# Config relating to git
git:
  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Pagers.md
  paging:
    # Value of the --color arg in the git diff command. Some pagers want this to be set to 'always' and some want it set to 'never'
    colorArg: always

    # e.g.
    # diff-so-fancy
    # delta --dark --paging=never
    # ydiff -p cat -s --wrap --width={{columnWidth}}
    pager: ""

    useConfig: false

    # e.g. 'difft --color=always'
    externalDiffCommand: ""

  # Config relating to committing
  commit:
    # If true, pass '--signoff' flag when committing
    signOff: false

    # Automatic WYSIWYG wrapping of the commit message as you type
    autoWrapCommitMessage: true

    # If autoWrapCommitMessage is true, the width to wrap to
    autoWrapWidth: 72

  # Config relating to merging
  merging:
    # If true, run merges in a subprocess so that if a commit message is required, Lazygit will not hang
    # Only applicable to unix users.
    manualCommit: false

    # Extra args passed to , e.g. --no-ff
    args: ""

    # The commit message to use for a squash merge commit. Can contain "{{selectedRef}}" and "{{currentBranch}}" placeholders.
    squashMergeMessage: Squash merge {{selectedRef}} into {{currentBranch}}

  # list of branches that are considered 'main' branches, used when displaying commits
  mainBranches:
    - master
    - main

  # Prefix to use when skipping hooks. E.g. if set to 'WIP', then pre-commit hooks will be skipped when the commit message starts with 'WIP'
  skipHookPrefix: WIP

  # If true, periodically fetch from remote
  autoFetch: true

  # If true, periodically refresh files and submodules
  autoRefresh: true

  # If true, pass the --all arg to git fetch
  fetchAll: true

  # If true, lazygit will automatically stage files that used to have merge
  # conflicts but no longer do; and it will also ask you if you want to
  # continue a merge or rebase if you've resolved all conflicts. If false, it
  # won't do either of these things.
  autoStageResolvedConflicts: true

  # Command used when displaying the current branch git log in the main window
  branchLogCmd: git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --

  # Command used to display git log of all branches in the main window.
  # Deprecated: Use allBranchesLogCmds instead.
  allBranchesLogCmd: git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium

  # If true, do not spawn a separate process when using GPG
  overrideGpg: false

  # If true, do not allow force pushes
  disableForcePushing: false

  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-branch-name-prefix
  branchPrefix: ""

  # If true, parse emoji strings in commit messages e.g. render :rocket: as 🚀
  # (This should really be under 'gui', not 'git')
  parseEmoji: false

  # Config for showing the log in the commits view
  log:
    # One of: 'date-order' | 'author-date-order' | 'topo-order' | 'default'
    # 'topo-order' makes it easier to read the git log graph, but commits may not
    # appear chronologically. See https://git-scm.com/docs/
    #
    # Deprecated: Configure this with Log menu -> Commit sort order (<c-l> in the commits window by default).
    order: topo-order

    # This determines whether the git graph is rendered in the commits panel
    # One of 'always' | 'never' | 'when-maximised'
    #
    # Deprecated: Configure this with Log menu -> Show git graph (<c-l> in the commits window by default).
    showGraph: always

    # displays the whole git graph by default in the commits view (equivalent to passing the --all argument to git log)
    showWholeGraph: false

  # When copying commit hashes to the clipboard, truncate them to this
  # length. Set to 40 to disable truncation.
  truncateCopiedCommitHashesTo: 12

# Periodic update checks
update:
  # One of: 'prompt' (default) | 'background' | 'never'
  method: prompt

  # Period in days between update checks
  days: 14

# Background refreshes
refresher:
  # File/submodule refresh interval in seconds.
  # Auto-refresh can be disabled via option 'git.autoRefresh'.
  refreshInterval: 10

  # Re-fetch interval in seconds.
  # Auto-fetch can be disabled via option 'git.autoFetch'.
  fetchInterval: 60

# If true, show a confirmation popup before quitting Lazygit
confirmOnQuit: false

# If true, exit Lazygit when the user presses escape in a context where there is nothing to cancel/close
quitOnTopLevelReturn: false

# Config relating to things outside of Lazygit like how files are opened, copying to clipboard, etc
os:
  # Command for editing a file. Should contain "{{filename}}".
  edit: ""

  # Command for editing a file at a given line number. Should contain
  # "{{filename}}", and may optionally contain "{{line}}".
  editAtLine: ""

  # Same as EditAtLine, except that the command needs to wait until the
  # window is closed.
  editAtLineAndWait: ""

  # Whether lazygit suspends until an edit process returns
  editInTerminal: false

  # For opening a directory in an editor
  openDirInEditor: ""

  # A built-in preset that sets all of the above settings. Supported presets
  # are defined in the getPreset function in editor_presets.go.
  editPreset: ""

  # Command for opening a file, as if the file is double-clicked. Should
  # contain "{{filename}}", but doesn't support "{{line}}".
  open: ""

  # Command for opening a link. Should contain "{{link}}".
  openLink: ""

  # EditCommand is the command for editing a file.
  # Deprecated: use Edit instead. Note that semantics are different:
  # EditCommand is just the command itself, whereas Edit contains a
  # "{{filename}}" variable.
  editCommand: ""

  # EditCommandTemplate is the command template for editing a file
  # Deprecated: use EditAtLine instead.
  editCommandTemplate: ""

  # OpenCommand is the command for opening a file
  # Deprecated: use Open instead.
  openCommand: ""

  # OpenLinkCommand is the command for opening a link
  # Deprecated: use OpenLink instead.
  openLinkCommand: ""

  # CopyToClipboardCmd is the command for copying to clipboard.
  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-command-for-copying-to-and-pasting-from-clipboard
  copyToClipboardCmd: ""

  # ReadFromClipboardCmd is the command for reading the clipboard.
  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#custom-command-for-copying-to-and-pasting-from-clipboard
  readFromClipboardCmd: ""

# If true, don't display introductory popups upon opening Lazygit.
disableStartupPopups: false

# What to do when opening Lazygit outside of a git repo.
# - 'prompt': (default) ask whether to initialize a new repo or open in the most recent repo
# - 'create': initialize a new repo
# - 'skip': open most recent repo
# - 'quit': exit Lazygit
notARepository: prompt

# If true, display a confirmation when subprocess terminates. This allows you to view the output of the subprocess before returning to Lazygit.
promptToReturnFromSubprocess: true

# Keybindings
keybinding:
  universal:
    quit: q
    quit-alt1: <c-c>
    return: <esc>
    quitWithoutChangingDirectory: Q
    togglePanel: <tab>
    prevItem: <up>
    nextItem: <down>
    prevItem-alt: k
    nextItem-alt: j
    prevPage: ','
    nextPage: .
    scrollLeft: H
    scrollRight: L
    gotoTop: <
    gotoBottom: '>'
    toggleRangeSelect: v
    rangeSelectDown: <s-down>
    rangeSelectUp: <s-up>
    prevBlock: <left>
    nextBlock: <right>
    prevBlock-alt: h
    nextBlock-alt: l
    nextBlock-alt2: <tab>
    prevBlock-alt2: <backtab>
    jumpToBlock:
      - "1"
      - "2"
      - "3"
      - "4"
      - "5"
    nextMatch: "n"
    prevMatch: "N"
    startSearch: /
    optionMenu: <disabled>
    optionMenu-alt1: '?'
    select: <space>
    goInto: <enter>
    confirm: <enter>
    confirmInEditor: <a-enter>
    remove: d
    new: "n"
    edit: e
    openFile: o
    scrollUpMain: <pgup>
    scrollDownMain: <pgdown>
    scrollUpMain-alt1: K
    scrollDownMain-alt1: J
    scrollUpMain-alt2: <c-u>
    scrollDownMain-alt2: <c-d>
    executeShellCommand: ':'
    createRebaseOptionsMenu: m

    # 'Files' appended for legacy reasons
    pushFiles: P

    # 'Files' appended for legacy reasons
    pullFiles: p
    refresh: R
    createPatchOptionsMenu: <c-p>
    nextTab: ']'
    prevTab: '['
    nextScreenMode: +
    prevScreenMode: _
    undo: z
    redo: Z
    filteringMenu: <c-s>
    diffingMenu: W
    diffingMenu-alt: <c-e>
    copyToClipboard: <c-o>
    openRecentRepos: <c-r>
    submitEditorText: <enter>
    extrasMenu: '@'
    toggleWhitespaceInDiffView: <c-w>
    increaseContextInDiffView: '}'
    decreaseContextInDiffView: '{'
    increaseRenameSimilarityThreshold: )
    decreaseRenameSimilarityThreshold: (
    openDiffTool: <c-t>
  status:
    checkForUpdate: u
    recentRepos: <enter>
    allBranchesLogGraph: a
  files:
    commitChanges: c
    commitChangesWithoutHook: w
    amendLastCommit: A
    commitChangesWithEditor: C
    findBaseCommitForFixup: <c-f>
    confirmDiscard: x
    ignoreFile: i
    refreshFiles: r
    stashAllChanges: s
    viewStashOptions: S
    toggleStagedAll: a
    viewResetOptions: D
    fetch: f
    openMergeTool: M
    openStatusFilter: <c-b>
    copyFileInfoToClipboard: "y"
    collapseAll: '-'
    expandAll: =
  branches:
    createPullRequest: o
    viewPullRequestOptions: O
    copyPullRequestURL: <c-y>
    checkoutBranchByName: c
    forceCheckoutBranch: F
    rebaseBranch: r
    renameBranch: R
    mergeIntoCurrentBranch: M
    viewGitFlowOptions: i
    fastForward: f
    createTag: T
    pushTag: P
    setUpstream: u
    fetchRemote: f
    sortOrder: s
  worktrees:
    viewWorktreeOptions: w
  commits:
    squashDown: s
    renameCommit: r
    renameCommitWithEditor: R
    viewResetOptions: g
    markCommitAsFixup: f
    createFixupCommit: F
    squashAboveCommits: S
    moveDownCommit: <c-j>
    moveUpCommit: <c-k>
    amendToCommit: A
    resetCommitAuthor: a
    pickCommit: p
    revertCommit: t
    cherryPickCopy: C
    pasteCommits: V
    markCommitAsBaseForRebase: B
    tagCommit: T
    checkoutCommit: <space>
    resetCherryPick: <c-R>
    copyCommitAttributeToClipboard: "y"
    openLogMenu: <c-l>
    openInBrowser: o
    viewBisectOptions: b
    startInteractiveRebase: i
  amendAttribute:
    resetAuthor: a
    setAuthor: A
    addCoAuthor: c
  stash:
    popStash: g
    renameStash: r
  commitFiles:
    checkoutCommitFile: c
  main:
    toggleSelectHunk: a
    pickBothHunks: b
    editSelectHunk: E
  submodules:
    init: i
    update: u
    bulkMenu: b
  commitMessage:
    commitMenu: <c-o>
`)

func BenchmarkMigrationOnLargeConfiguration(b *testing.B) {
	for b.Loop() {
		changes := NewChangesSet()
		_, _, _ = computeMigratedConfig("path doesn't matter", largeConfiguration, changes)
	}
}

func TestAllBranchesLogCmdMigrations(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Incomplete Configuration Passes uneventfully",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Single Cmd with no Cmds",
			input: `git:
  allBranchesLogCmd: git log --graph --oneline
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Created git.allBranchesLogCmds array containing value of git.allBranchesLogCmd",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Cmd with one existing Cmds",
			input: `git:
  allBranchesLogCmd: git log --graph --oneline
  allBranchesLogCmds:
    - git log --graph --oneline --pretty
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
    - git log --graph --oneline --pretty
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Prepended git.allBranchesLogCmd value to git.allBranchesLogCmds array",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Only Cmds set have no changes",
			input: `git:
  allBranchesLogCmds:
    - git log
`,
			expected:        "",
			expectedChanges: []string{},
		},
		{
			name: "Removes Empty Cmd When at end of yaml",
			input: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  allBranchesLogCmd:
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Removed obsolete git.allBranchesLogCmd"},
		},
		{
			name: "Migrates when sequence defined inline",
			input: `git:
  allBranchesLogCmds: [foo, bar]
  allBranchesLogCmd: baz
`,
			expected: `git:
  allBranchesLogCmds: [baz, foo, bar]
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Prepended git.allBranchesLogCmd value to git.allBranchesLogCmds array",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Removes Empty Cmd With Keys Afterwards",
			input: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  allBranchesLogCmd:
  foo: bar
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  foo: bar
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Removed obsolete git.allBranchesLogCmd"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}
