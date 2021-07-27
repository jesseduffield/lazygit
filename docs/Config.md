# User Config

Default path for the config file:

- Linux: `~/.config/lazygit/config.yml`
- MacOS: `~/Library/Application Support/lazygit/config.yml`
- Windows: `%APPDATA%\lazygit\config.yml`

For old installations (slightly embarrassing: I didn't realise at the time that you didn't need to supply a vendor name to the path so I just used my name):

- Linux: `~/.config/jesseduffield/lazygit/config.yml`
- MacOS: `~/Library/Application Support/jesseduffield/lazygit/config.yml`
- Windows: `%APPDATA%\jesseduffield\lazygit\config.yml`

## Default

```yaml
gui:
  # stuff relating to the UI
  scrollHeight: 2 # how many lines you scroll by
  scrollPastBottom: true # enable scrolling past the bottom
  sidePanelWidth: 0.3333 # number from 0 to 1
  expandFocusedSidePanel: false
  mainPanelSplitMode: 'flexible' # one of 'horizontal' | 'flexible' | 'vertical'
  theme:
    lightTheme: false # For terminals with a light background
    activeBorderColor:
      - white
      - bold
    inactiveBorderColor:
      - green
    optionsTextColor:
      - blue
    selectedLineBgColor:
      - default
    selectedRangeBgColor:
      - blue
  commitLength:
    show: true
  mouseEvents: true
  skipUnstageLineWarning: false
  skipStashWarning: true
  showFileTree: false # for rendering changes files in a tree format
  showListFooter: true # for seeing the '5 of 20' message in list panels
  showRandomTip: true
  showCommandLog: true
  commandLogSize: 8
git:
  paging:
    colorArg: always
    useConfig: false
  merging:
    # only applicable to unix users
    manualCommit: false
    # extra args passed to `git merge`, e.g. --no-ff
    args: ''
  pull:
    mode: 'auto' # one of 'auto' | 'merge' | 'rebase' | 'ff-only', auto reads from git configuration
  skipHookPrefix: WIP
  autoFetch: true
  branchLogCmd: 'git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --'
  allBranchesLogCmd: 'git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium'
  overrideGpg: false # prevents lazygit from spawning a separate process when using GPG
  disableForcePushing: false
os:
  editCommand: '' # see 'Configuring File Editing' section
  openCommand: ''
refresher:
  refreshInterval: 10 # file/submodule refresh interval in seconds
  fetchInterval: 60 # re-fetch interval in seconds
update:
  method: prompt # can be: prompt | background | never
  days: 14 # how often an update is checked for
reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
confirmOnQuit: false
# determines whether hitting 'esc' will quit the application when there is nothing to cancel/close
quitOnTopLevelReturn: false
disableStartupPopups: false
notARepository: 'prompt' # one of: 'prompt' | 'create' | 'skip'
keybinding:
  universal:
    quit: 'q'
    quit-alt1: '<c-c>' # alternative/alias of quit
    return: '<esc>' # return to previous menu, will quit if there's nowhere to return
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
    prevBlock: '<left>' # goto the previous block / panel
    nextBlock: '<right>' # goto the next block / panel
    prevBlock-alt: 'h' # goto the previous block / panel
    nextBlock-alt: 'l' # goto the next block / panel
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
    appendNewline: '<tab>'
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
  branches:
    createPullRequest: 'o'
    viewPullRequestOptions: 'O'
    checkoutBranchByName: 'c'
    forceCheckoutBranch: 'F'
    rebaseBranch: 'r'
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
  stash:
    popStash: 'g'
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
  openCommand: 'cmd /c "start "" {{filename}}"'
```

### Linux

```yaml
os:
  openCommand: 'sh -c "xdg-open {{filename}} >/dev/null"'
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

### Recommended Config Values

for users of VSCode

```yaml
os:
  openCommand: 'code -rg {{filename}}'
```

## Color Attributes

For color attributes you can choose an array of attributes (with max one color attribute)
The available attributes are:

- default
- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white
- bold
- reverse # useful for high-contrast
- underline

## Light terminal theme

If you have issues with a light terminal theme where you can't read / see the text add these settings

```yaml
gui:
  theme:
    lightTheme: true
    activeBorderColor:
      - black
      - bold
    inactiveBorderColor:
      - black
    selectedLineBgColor:
      - default
```

## Struggling to see selected line

If you struggle to see the selected line I recommend using the reverse attribute on selected lines like so:

```yaml
gui:
  theme:
    selectedLineBgColor:
      - reverse
    selectedRangeBgColor:
      - reverse
```

The following has also worked for a couple of people:

```yaml
gui:
  theme:
    activeBorderColor:
      - white
      - bold
    inactiveBorderColor:
      - white
    selectedLineBgColor:
      - reverse
      - blue
```

Alternatively you may have bold fonts disabled in your terminal, in which case enabling bold fonts should solve the problem.

If you're still having trouble please raise an issue.

## Example Coloring

![border example](../../assets/colored-border-example.png)

## Keybindings

For all possible keybinding options, check [Custom_Keybindings.md](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md)

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
- `provider` is one of `github`, `bitbucket` or `gitlab`
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
