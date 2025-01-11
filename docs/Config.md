# User Config

Default path for the global config file:

- Linux: `~/.config/lazygit/config.yml`
- MacOS: `~/Library/Application\ Support/lazygit/config.yml`
- Windows: `%LOCALAPPDATA%\lazygit\config.yml` (default location, but it will also be found in `%APPDATA%\lazygit\config.yml`

For old installations (slightly embarrassing: I didn't realise at the time that you didn't need to supply a vendor name to the path so I just used my name):

- Linux: `~/.config/jesseduffield/lazygit/config.yml`
- MacOS: `~/Library/Application\ Support/jesseduffield/lazygit/config.yml`
- Windows: `%APPDATA%\jesseduffield\lazygit\config.yml`

If you want to change the config directory:

- MacOS: `export XDG_CONFIG_HOME="$HOME/.config"`

In addition to the global config file you can create repo-specific config files in `<repo>/.git/lazygit.yml`. Settings in these files override settings in the global config file. In addition, files called `.lazygit.yml` in any of the parent directories of a repo will also be loaded; this can be useful if you have settings that you want to apply to a group of repositories.

JSON schema is available for `config.yml` so that IntelliSense in Visual Studio Code (completion and error checking) is automatically enabled when the [YAML Red Hat][yaml] extension is installed. However, note that automatic schema detection only works if your config file is in one of the standard paths mentioned above. If you override the path to the file, you can still make IntelliSense work by adding

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/jesseduffield/lazygit/master/schema/config.json
```

to the top of your config file or via [Visual Studio Code settings.json config][settings].

[yaml]: https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml
[settings]: https://github.com/redhat-developer/vscode-yaml#associating-a-schema-to-a-glob-pattern-via-yamlschemas

## Default

<!-- START CONFIG YAML: AUTOMATICALLY GENERATED with `go generate ./..., DO NOT UPDATE MANUALLY -->
```yaml
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

  # If true, capture mouse events.
  # When mouse events are captured, it's a little harder to select text: e.g. requiring you to hold the option key when on macOS.
  mouseEvents: true

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
  # twice as tall as the other panels. Only relevant if `expandFocusedSidePanel` is true.
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
  # This can be toggled from within Lazygit with the '~' key, but that will not change the default.
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
  # One of 'rounded' (default) | 'single' | 'double' | 'hidden'
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

    # If true, Lazygit will use whatever pager is specified in `$GIT_PAGER`, `$PAGER`, or your *git config*. If the pager ends with something like ` | less` we will strip that part out, because less doesn't play nice with our rendering approach. If the custom pager uses less under the hood, that will also break rendering (hence the `--paging=never` flag for the `delta` pager).
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

    # Extra args passed to `git merge`, e.g. --no-ff
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
  # Deprecated: Use `allBranchesLogCmds` instead.
  allBranchesLogCmd: git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium

  # If true, do not spawn a separate process when using GPG
  overrideGpg: false

  # If true, do not allow force pushes
  disableForcePushing: false

  # See https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#predefined-commit-message-prefix
  commitPrefix:
    # pattern to match on. E.g. for 'feature/AB-123' to match on the AB-123 use "^\\w+\\/(\\w+-\\w+).*"
    pattern: ""

    # Replace directive. E.g. for 'feature/AB-123' to start the commit message with 'AB-123 ' use "[$1] "
    replace: ""

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
    # Deprecated: Configure this with `Log menu -> Commit sort order` (<c-l> in the commits window by default).
    order: topo-order

    # This determines whether the git graph is rendered in the commits panel
    # One of 'always' | 'never' | 'when-maximised'
    #
    # Deprecated: Configure this with `Log menu -> Show git graph` (<c-l> in the commits window by default).
    showGraph: always

    # displays the whole git graph by default in the commits view (equivalent to passing the `--all` argument to `git log`)
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
    redo: <c-z>
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
    toggleTreeView: '`'
    openMergeTool: M
    openStatusFilter: <c-b>
    copyFileInfoToClipboard: "y"
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
```
<!-- END CONFIG YAML -->

## Platform Defaults

### Windows

```yaml
os:
  open: 'start "" {{filename}}'
```

### Linux

```yaml
os:
  open: 'xdg-open {{filename}} >/dev/null'
```

### OSX

```yaml
os:
  open: 'open {{filename}}'
```

## Custom Command for Opening a Link
```yaml
os:
  openLink: 'bash -C /path/to/your/shell-script.sh {{link}}'
```
Specify the external command to invoke when opening URL links (i.e. creating MR/PR in GitLab, BitBucket or GitHub). `{{link}}` will be replaced by the URL to be opened. A simple shell script can be used to further mangle the passed URL.

## Custom Command for Copying to and Pasting from Clipboard
```yaml
os:
  copyToClipboardCmd: ''
```
Specify an external command to invoke when copying to clipboard is requested. `{{text}` will be replaced by text to be copied. Default is to copy to system clipboard.

If you are working on a terminal that supports OSC52, the following command will let you take advantage of it:
```yaml
os:
  copyToClipboardCmd: printf "\033]52;c;$(printf {{text}} | base64 -w 0)\a" > /dev/tty
```

For tmux you need to wrap it with the [tmux escape sequence](https://github.com/tmux/tmux/wiki/FAQ#what-is-the-passthrough-escape-sequence-and-how-do-i-use-it), and enable passthrough in tmux config with `set -g allow-passthrough on`:
```yaml
os:
  copyToClipboardCmd: printf "\033Ptmux;\033\033]52;c;$(printf {{text}} | base64 -w 0)\a\033\\" > /dev/tty
```

For the best of both worlds, we can let the command determine if we are running in a tmux session and send the correct sequence:
```yaml
os:
  copyToClipboardCmd: >
    if [[ "$TERM" =~ ^(screen|tmux) ]]; then
      printf "\033Ptmux;\033\033]52;c;$(printf {{text}} | base64 -w 0)\a\033\\" > /dev/tty
    else
      printf "\033]52;c;$(printf {{text}} | base64 -w 0)\a" > /dev/tty
    fi
```

A custom command for reading from the clipboard can be set using
```yaml
os:
  readFromClipboardCmd: ''
```
It is used, for example, when pasting a commit message into the commit message panel. The command is supposed to output the clipboard content to stdout.

## Configuring File Editing

There are two commands for opening files, `o` for "open" and `e` for "edit". `o` acts as if the file was double-clicked in the Finder/Explorer, so it also works for non-text files, whereas `e` opens the file in an editor. `e` can also jump to the right line in the file if you invoke it from the staging panel, for example.

To tell lazygit which editor to use for the `e` command, the easiest way to do that is to provide an editPreset config, e.g.

```yaml
os:
  editPreset: 'vscode'
```

Supported presets are `vim`, `nvim`, `nvim-remote`, `lvim`, `emacs`, `nano`, `micro`, `vscode`, `sublime`, `bbedit`, `kakoune`, `helix`, `xcode`, and `zed`. In many cases lazygit will be able to guess the right preset from your $(git config core.editor), or an environment variable such as $VISUAL or $EDITOR.

`nvim-remote` is an experimental preset for when you have invoked lazygit from within a neovim process, allowing lazygit to open the file from within the parent process rather than spawning a new one.

If for some reason you are not happy with the default commands from a preset, or there simply is no preset for your editor, you can customize the commands by setting the `edit`, `editAtLine`, and `editAtLineAndWait` options, e.g.:

```yaml
os:
  edit: 'myeditor {{filename}}'
  editAtLine: 'myeditor --line={{line}} {{filename}}'
  editAtLineAndWait: 'myeditor --block --line={{line}} {{filename}}'
  editInTerminal: true
  openDirInEditor: 'myeditor {{dir}}'
```

The `editInTerminal` option is used to decide whether lazygit needs to suspend itself to the background before calling the editor. It should really be named `suspend` because for some cases like when lazygit is opened from within a neovim session and you're using the `nvim-remote` preset, you're technically still in a terminal. Nonetheless we're sticking with the name `editInTerminal` for backwards compatibility.

Contributions of new editor presets are welcome; see the `getPreset` function in [`editor_presets.go`](https://github.com/jesseduffield/lazygit/blob/master/pkg/config/editor_presets.go).

## Overriding default config file location

To override the default config directory, use `CONFIG_DIR="$HOME/.config/lazygit"`. This directory contains the config file in addition to some other files lazygit uses to keep track of state across sessions.

To override the individual config file used, use the `--use-config-file` arg or the `LG_CONFIG_FILE` env var.

If you want to merge a specific config file into a more general config file, perhaps for the sake of setting some theme-specific options, you can supply a list of comma-separated config file paths, like so:

```sh
lazygit --use-config-file="$HOME/.base_lg_conf,$HOME/.light_theme_lg_conf"
or
LG_CONFIG_FILE="$HOME/.base_lg_conf,$HOME/.light_theme_lg_conf" lazygit
```

## Scroll-off Margin

When the selected line gets close to the bottom of the window and you hit down-arrow, there's a feature called "scroll-off margin" that lets the view scroll a little earlier so that you can see a bit of what's coming in the direction that you are moving. This is controlled by the `gui.scrollOffMargin` setting (default: 2), so it keeps 2 lines below the selection visible as you scroll down. It can be set to 0 to scroll only when the selection reaches the bottom of the window.

That's the behavior when `gui.scrollOffBehavior` is set to "margin" (the default). If you set `gui.scrollOffBehavior` to "jump", then upon reaching the last line of a view and hitting down-arrow the view will scroll by half a page so that the selection ends up in the middle of the view. This may feel a little jarring because the cursor jumps around when continuously moving down, but it has the advantage that the view doesn't scroll as often.

This setting applies both to all list views (e.g. commits and branches etc), and to the staging view.

## Filtering

We have two ways to filter things, substring matching (the default) and fuzzy searching. With substring matching, the text you enter gets searched for verbatim (usually case-insensitive, except when your filter string contains uppercase letters, in which case we search case-sensitively). You can search for multiple non-contiguous substrings by separating them with spaces; for example, "int test" will match "integration-testing". All substrings have to match, but not necessarily in the given order.

Fuzzy searching is smarter in that it allows every letter of the filter string to match anywhere in the text (only in order though), assigning a weight to the quality of the match and sorting by that order. This has the advantage that it allows typing "clt" to match "commit_loader_test" (letters at the beginning of subwords get more weight); but it has the disadvantage that it tends to return lots of irrelevant results, especially with short filter strings.

## Color Attributes

For color attributes you can choose an array of attributes (with max one color attribute)
The available attributes are:

**Colors**

- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white
- '#ff00ff'

**Modifiers**

- bold
- default
- reverse # useful for high-contrast
- underline
- strikethrough

## Highlighting the selected line

If you don't like the default behaviour of highlighting the selected line with a blue background, you can use the `selectedLineBgColor` key to customise the behaviour. If you just want to embolden the selected line (this was the original default), you can do the following:

```yaml
gui:
  theme:
    selectedLineBgColor:
      - default
```

You can also use the reverse attribute like so:

```yaml
gui:
  theme:
    selectedLineBgColor:
      - reverse
```

## Custom Author Color

Lazygit will assign a random color for every commit author in the commits pane by default.

You can customize the color in case you're not happy with the randomly assigned one:

```yaml
gui:
  authorColors:
    'John Smith': 'red' # use red for John Smith
    'Alan Smithee': '#00ff00' # use green for Alan Smithee
```

You can use wildcard to set a unified color in case your are lazy to customize the color for every author or you just want a single color for all/other authors:

```yaml
gui:
  authorColors:
    # use red for John Smith
    'John Smith': 'red'
    # use blue for other authors
    '*': '#0000ff'
```

## Custom Branch Color

You can customize the color of branches based on the branch prefix:

```yaml
gui:
  branchColors:
    'docs': '#11aaff' # use a light blue for branches beginning with 'docs/'
```

## Example Coloring

![border example](../../assets/colored-border-example.png)

## Display Nerd Fonts Icons

If you are using [Nerd Fonts](https://www.nerdfonts.com), you can display icons.

```yaml
gui:
  nerdFontsVersion: "3"
```

Supported versions are "2" and "3". The deprecated config `showIcons` sets the version to "2" for backwards compatibility.

## Keybindings

For all possible keybinding options, check [Custom_Keybindings.md](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md)

You can disable certain key bindings by specifying `<disabled>`.

```yaml
keybinding:
  universal:
    edit: <disabled> # disable 'edit file'
```

### Example Keybindings For Colemak Users

```yaml
keybinding:
  universal:
    prevItem-alt: 'u'
    nextItem-alt: 'e'
    prevBlock-alt: 'n'
    nextBlock-alt: 'i'
    nextMatch: '='
    prevMatch: '-'
    new: 'k'
    edit: 'o'
    openFile: 'O'
    scrollUpMain-alt1: 'U'
    scrollDownMain-alt1: 'E'
    scrollUpMain-alt2: '<c-u>'
    scrollDownMain-alt2: '<c-e>'
    undo: 'l'
    redo: '<c-r>'
    diffingMenu: 'M'
    filteringMenu: '<c-f>'
  files:
    ignoreFile: 'I'
  commits:
    moveDownCommit: '<c-e>'
    moveUpCommit: '<c-u>'
  branches:
    viewGitFlowOptions: 'I'
    setUpstream: 'U'
```

## Custom pull request URLs

Some git provider setups (e.g. on-premises GitLab) can have distinct URLs for git-related calls and the web interface/API itself. To work with those, Lazygit needs to know where it needs to create the pull request. You can do so on your `config.yml` file using the following syntax:

```yaml
services:
  '<gitDomain>': '<provider>:<webDomain>'
```

Where:

- `gitDomain` stands for the domain used by git itself (i.e. the one present on clone URLs), e.g. `git.work.com`
- `provider` is one of `github`, `bitbucket`, `bitbucketServer`, `azuredevops`, `gitlab` or `gitea`
- `webDomain` is the URL where your git service exposes a web interface and APIs, e.g. `gitservice.work.com`

## Predefined commit message prefix

In situations where certain naming pattern is used for branches and commits, pattern can be used to populate commit message with prefix that is parsed from the branch name.

Example:

- Branch name: feature/AB-123
- Commit message: [AB-123] Adding feature

```yaml
git:
  commitPrefix:
    pattern: "^\\w+\\/(\\w+-\\w+).*"
    replace: '[$1] '
```

If you want repository-specific prefixes, you can map them with `commitPrefixes`. If you have both `commitPrefixes` defined and an entry in `commitPrefixes` for the current repo, the `commitPrefixes` entry is given higher precedence. Repository folder names must be an exact match.

```yaml
git:
  commitPrefixes:
    my_project: # This is repository folder name
      pattern: "^\\w+\\/(\\w+-\\w+).*"
      replace: '[$1] '
```

> [!IMPORTANT]
> The way golang regex works is when you use `$n` in the replacement string, where `n` is a number, it puts the nth captured subgroup at that place. If `n` is out of range because there aren't that many capture groups in the regex, it puts an empty string there.
>
> So make sure you are capturing group or groups in your regex.
>
> For example `^[A-Z]+-\d+$` won't work on branch name like BRANCH-1111
> But `^([A-Z]+-\d+)$` will

## Predefined branch name prefix

In situations where certain naming pattern is used for branches, this can be used to populate new branch creation with a static prefix.

Example:

Some branches:
- jsmith/AB-123
- cwilson/AB-125

```yaml
git:
  branchPrefix: "firstlast/"
```

## Custom git log command

You can override the `git log` command that's used to render the log of the selected branch like so:

```
git:
  branchLogCmd: "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium --oneline {{branchName}} --"
```

Result:

![](https://i.imgur.com/Nibq35B.png)

## Launching not in a repository behaviour

By default, when launching lazygit from a directory that is not a repository, you will be prompted to choose if you would like to initialize a repo. You can override this behaviour in the config with one of the following:

```yaml
# for default prompting behaviour
notARepository: 'prompt'
```

```yaml
# to skip and initialize a new repo
notARepository: 'create'
```

```yaml
# to skip without creating a new repo
notARepository: 'skip'
```

```yaml
# to exit immediately if run outside of the Git repository
notARepository: 'quit'
```
