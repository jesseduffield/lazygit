# User Config

Default path for the config file:

- Linux: `~/.config/lazygit/config.yml`
- MacOS: `~/Library/Application Support/lazygit/config.yml`
- Windows: `%APPDATA%\lazygit\config.yml`

For old installations (slightly embarrassing: I didn't realise at the time that you didn't need to supply a vendor name to the path so I just used my name):

- Linux: `~/.config/jesseduffield/lazygit/config.yml`
- MacOS: `~/Library/Application Support/jesseduffield/lazygit/config.yml`
- Windows: `%APPDATA%\jesseduffield\lazygit\config.yml`

If you want to change the config directory:

- MacOS: `export XDG_CONFIG_HOME="$HOME/.config"`

## Default

```yaml
gui:
  # stuff relating to the UI
  windowSize: 'normal' # one of 'normal' | 'half' | 'full' default is 'normal'
  scrollHeight: 2 # how many lines you scroll by
  scrollPastBottom: true # enable scrolling past the bottom
  sidePanelWidth: 0.3333 # number from 0 to 1
  expandFocusedSidePanel: false
  mainPanelSplitMode: 'flexible' # one of 'horizontal' | 'flexible' | 'vertical'
  language: 'auto' # one of 'auto' | 'en' | 'zh' | 'pl' | 'nl' | 'ja' | 'ko'
  timeFormat: '02 Jan 06 15:04 MST' # https://pkg.go.dev/time#Time.Format
  theme:
    activeBorderColor:
      - green
      - bold
    inactiveBorderColor:
      - white
    optionsTextColor:
      - blue
    selectedLineBgColor:
      - blue # set to `default` to have no background colour
    selectedRangeBgColor:
      - blue
    cherryPickedCommitBgColor:
      - cyan
    cherryPickedCommitFgColor:
      - blue
    unstagedChangesColor:
      - red
    defaultFgColor:
      - default
  commitLength:
    show: true
  mouseEvents: true
  skipUnstageLineWarning: false
  skipStashWarning: false
  showFileTree: true # for rendering changes files in a tree format
  showListFooter: true # for seeing the '5 of 20' message in list panels
  showRandomTip: true
  showBottomLine: true # for hiding the bottom information line (unless it has important information to tell you)
  showCommandLog: true
  showIcons: false
  commandLogSize: 8
  splitDiff: 'auto' # one of 'auto' | 'always'
  skipRewordInEditorWarning: false # for skipping the confirmation before launching the reword editor
git:
  paging:
    colorArg: always
    useConfig: false
  commit:
    signOff: false
    verbose: default # one of 'default' | 'always' | 'never'
  merging:
    # only applicable to unix users
    manualCommit: false
    # extra args passed to `git merge`, e.g. --no-ff
    args: ''
  log:
    # one of date-order, author-date-order, topo-order or default.
    # topo-order makes it easier to read the git log graph, but commits may not
    # appear chronologically. See https://git-scm.com/docs/git-log#_commit_ordering
    order: 'topo-order'
    # one of always, never, when-maximised
    # this determines whether the git graph is rendered in the commits panel
    showGraph: 'when-maximised'
    # displays the whole git graph by default in the commits panel (equivalent to passing the `--all` argument to `git log`)
    showWholeGraph: false
  skipHookPrefix: WIP
  autoFetch: true
  autoRefresh: true
  branchLogCmd: 'git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --'
  allBranchesLogCmd: 'git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium'
  overrideGpg: false # prevents lazygit from spawning a separate process when using GPG
  disableForcePushing: false
  parseEmoji: false
  diffContextSize: 3 # how many lines of context are shown around a change in diffs
os:
  editCommand: '' # see 'Configuring File Editing' section
  editCommandTemplate: ''
  openCommand: ''
refresher:
  refreshInterval: 10 # File/submodule refresh interval in seconds. Auto-refresh can be disabled via option 'git.autoRefresh'.
  fetchInterval: 60 # Re-fetch interval in seconds. Auto-fetch can be disabled via option 'git.autoFetch'.
update:
  method: prompt # can be: prompt | background | never
  days: 14 # how often an update is checked for
confirmOnQuit: false
# determines whether hitting 'esc' will quit the application when there is nothing to cancel/close
quitOnTopLevelReturn: false
disableStartupPopups: false
notARepository: 'prompt' # one of: 'prompt' | 'create' | 'skip' | 'quit'
promptToReturnFromSubprocess: true # display confirmation when subprocess terminates
keybinding:
  universal:
    quit: 'q'
    quit-alt1: '<c-c>' # alternative/alias of quit
    return: '<esc>' # return to previous menu, will quit if there's nowhere to return
    # When set to a printable character, this will work for returning from non-prompt panels
    return-alt1: null
    quitWithoutChangingDirectory: 'Q'
    togglePanel: '<tab>' # goto the next panel
    prevItem: '<up>' # go one line up
    nextItem: '<down>' # go one line down
    prevItem-alt: 'k' # go one line up
    nextItem-alt: 'j' # go one line down
    prevPage: ',' # go to next page in list
    nextPage: '.' # go to previous page in list
    gotoTop: '<' # go to top of list
    gotoBottom: '>' # go to bottom of list
    scrollLeft: 'H' # scroll left within list view
    scrollRight: 'L' # scroll right within list view
    prevBlock: '<left>' # goto the previous block / panel
    nextBlock: '<right>' # goto the next block / panel
    prevBlock-alt: 'h' # goto the previous block / panel
    nextBlock-alt: 'l' # goto the next block / panel
    jumpToBlock: ['1', '2', '3', '4', '5'] # goto the Nth block / panel
    nextMatch: 'n'
    prevMatch: 'N'
    optionMenu: 'x' # show help menu
    optionMenu-alt1: '?' # show help menu
    select: '<space>'
    goInto: '<enter>'
    openRecentRepos: '<c-r>'
    confirm: '<enter>'
    confirm-alt1: 'y'
    remove: 'd'
    new: 'n'
    edit: 'e'
    openFile: 'o'
    scrollUpMain: '<pgup>' # main panel scroll up
    scrollDownMain: '<pgdown>' # main panel scroll down
    scrollUpMain-alt1: 'K' # main panel scroll up
    scrollDownMain-alt1: 'J' # main panel scroll down
    scrollUpMain-alt2: '<c-u>' # main panel scroll up
    scrollDownMain-alt2: '<c-d>' # main panel scroll down
    executeCustomCommand: ':'
    createRebaseOptionsMenu: 'm'
    pushFiles: 'P'
    pullFiles: 'p'
    refresh: 'R'
    createPatchOptionsMenu: '<c-p>'
    nextTab: ']'
    prevTab: '['
    nextScreenMode: '+'
    prevScreenMode: '_'
    undo: 'z'
    redo: '<c-z>'
    filteringMenu: '<c-s>'
    diffingMenu: 'W'
    diffingMenu-alt: '<c-e>' # deprecated
    copyToClipboard: '<c-o>'
    submitEditorText: '<enter>'
    appendNewline: '<a-enter>'
    extrasMenu: '@'
    toggleWhitespaceInDiffView: '<c-w>'
    increaseContextInDiffView: '}'
    decreaseContextInDiffView: '{'
  status:
    checkForUpdate: 'u'
    recentRepos: '<enter>'
  files:
    commitChanges: 'c'
    commitChangesWithoutHook: 'w' # commit changes without pre-commit hook
    amendLastCommit: 'A'
    commitChangesWithEditor: 'C'
    ignoreFile: 'i'
    refreshFiles: 'r'
    stashAllChanges: 's'
    viewStashOptions: 'S'
    toggleStagedAll: 'a' # stage/unstage all
    viewResetOptions: 'D'
    fetch: 'f'
    toggleTreeView: '`'
    openMergeTool: 'M'
    openStatusFilter: '<c-b>'
  branches:
    createPullRequest: 'o'
    viewPullRequestOptions: 'O'
    checkoutBranchByName: 'c'
    forceCheckoutBranch: 'F'
    rebaseBranch: 'r'
    renameBranch: 'R'
    mergeIntoCurrentBranch: 'M'
    viewGitFlowOptions: 'i'
    fastForward: 'f' # fast-forward this branch from its upstream
    pushTag: 'P'
    setUpstream: 'u' # set as upstream of checked-out branch
    fetchRemote: 'f'
  commits:
    squashDown: 's'
    renameCommit: 'r'
    renameCommitWithEditor: 'R'
    viewResetOptions: 'g'
    markCommitAsFixup: 'f'
    createFixupCommit: 'F' # create fixup commit for this commit
    squashAboveCommits: 'S'
    moveDownCommit: '<c-j>' # move commit down one
    moveUpCommit: '<c-k>' # move commit up one
    amendToCommit: 'A'
    pickCommit: 'p' # pick commit (when mid-rebase)
    revertCommit: 't'
    cherryPickCopy: 'c'
    cherryPickCopyRange: 'C'
    pasteCommits: 'v'
    tagCommit: 'T'
    checkoutCommit: '<space>'
    resetCherryPick: '<c-R>'
    copyCommitMessageToClipboard: '<c-y>'
    openLogMenu: '<c-l>'
    viewBisectOptions: 'b'
  stash:
    popStash: 'g'
    renameStash: 'r'
  commitFiles:
    checkoutCommitFile: 'c'
  main:
    toggleDragSelect: 'v'
    toggleDragSelect-alt: 'V'
    toggleSelectHunk: 'a'
    pickBothHunks: 'b'
  submodules:
    init: 'i'
    update: 'u'
    bulkMenu: 'b'
```

## Platform Defaults

### Windows

```yaml
os:
  openCommand: 'start "" {{filename}}'
```

### Linux

```yaml
os:
  openCommand: 'xdg-open {{filename}} >/dev/null'
```

### OSX

```yaml
os:
  openCommand: 'open {{filename}}'
```

### Configuring File Editing

Lazygit will edit a file with the first set editor in the following:

1. config.yaml

```yaml
os:
  editCommand: 'vim' # as an example
```

2. \$(git config core.editor)
3. \$GIT_EDITOR
4. \$VISUAL
5. \$EDITOR
6. \$(which vi)

Lazygit will log an error if none of these options are set.

You can specify the current line number when you're in the patch explorer.

```yaml
os:
  editCommand: 'vim'
  editCommandTemplate: '{{editor}} +{{line}} -- {{filename}}'
```

or

```yaml
os:
  editCommand: 'code'
  editCommandTemplate: '{{editor}} --goto -- {{filename}}:{{line}}'
```

`{{editor}}` in `editCommandTemplate` is replaced with the value of `editCommand`.

### Overriding default config file location

To override the default config directory, use `CONFIG_DIR="$HOME/.config/lazygit"`. This directory contains the config file in addition to some other files lazygit uses to keep track of state across sessions.

To override the individual config file used, use the `--use-config-file` arg or the `LG_CONFIG_FILE` env var.

If you want to merge a specific config file into a more general config file, perhaps for the sake of setting some theme-specific options, you can supply a list of comma-separated config file paths, like so:

```sh
lazygit --use-config-file="$HOME/.base_lg_conf,$HOME/.light_theme_lg_conf"
or
LG_CONFIG_FILE="$HOME/.base_lg_conf,$HOME/.light_theme_lg_conf" lazygit
```

### Recommended Config Values

for users of VSCode

```yaml
os:
  openCommand: 'code -rg {{filename}}'
```

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

## Highlighting the selected line

If you don't like the default behaviour of highlighting the selected line with a blue background, you can use the `selectedLineBgColor` and `selectedRangeBgColor` keys to customise the behaviour. If you just want to embolden the selected line (this was the original default), you can do the following:

```yaml
gui:
  theme:
    selectedLineBgColor:
      - default
    selectedRangeBgColor:
      - default
```

You can also use the reverse attribute like so:

```yaml
gui:
  theme:
    selectedLineBgColor:
      - reverse
    selectedRangeBgColor:
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
  showIcons: true
```

## Keybindings

For all possible keybinding options, check [Custom_Keybindings.md](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md)

You can disable certain key bindings by specifying `null`.

```yaml
keybinding:
  universal:
    edit: null # disable 'edit file'
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

Some git provider setups (e.g. on-premises GitLab) can have distinct URLs for git-related calls and
the web interface/API itself. To work with those, Lazygit needs to know where it needs to create
the pull request. You can do so on your `config.yml` file using the following syntax:

```yaml
services:
  '<gitDomain>': '<provider>:<webDomain>'
```

Where:

- `gitDomain` stands for the domain used by git itself (i.e. the one present on clone URLs), e.g. `git.work.com`
- `provider` is one of `github`, `bitbucket`, `bitbucketServer`, `azuredevops` or `gitlab`
- `webDomain` is the URL where your git service exposes a web interface and APIs, e.g. `gitservice.work.com`

## Predefined commit message prefix

In situations where certain naming pattern is used for branches and commits, pattern can be used to populate
commit message with prefix that is parsed from the branch name.

Example:

- Branch name: feature/AB-123
- Commit message: [AB-123] Adding feature

```yaml
git:
  commitPrefixes:
    my_project: # This is repository folder name
      pattern: "^\\w+\\/(\\w+-\\w+).*"
      replace: '[$1] '
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

By default, when launching lazygit from a directory that is not a repository,
you will be prompted to choose if you would like to initialize a repo. You can
override this behaviour in the config with one of the following:

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
