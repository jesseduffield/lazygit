# User Config

Default path for the config file:

- Linux: `~/.config/lazygit/config.yml`
- MacOS: `~/Library/Application\ Support/lazygit/config.yml`
- Windows: `%LOCALAPPDATA%\lazygit\config.yml` (default location, but it will also be found in `%APPDATA%\lazygit\config.yml`

For old installations (slightly embarrassing: I didn't realise at the time that you didn't need to supply a vendor name to the path so I just used my name):

- Linux: `~/.config/jesseduffield/lazygit/config.yml`
- MacOS: `~/Library/Application\ Support/jesseduffield/lazygit/config.yml`
- Windows: `%APPDATA%\jesseduffield\lazygit\config.yml`

If you want to change the config directory:

- MacOS: `export XDG_CONFIG_HOME="$HOME/.config"`

JSON schema is available for `config.yml` so that IntelliSense in Visual Studio Code (completion and error checking) is automatically enabled when the [YAML Red Hat][yaml] extension is installed. However, note that automatic schema detection only works if your config file is in one of the standard paths mentioned above. If you override the path to the file, you can still make IntelliSense work by adding

```yaml
# yaml-language-server: $schema=https://json.schemastore.org/lazygit.json
```

to the top of your config file or via [Visual Studio Code settings.json config][settings].

[yaml]: https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml
[settings]: https://github.com/redhat-developer/vscode-yaml#associating-a-schema-to-a-glob-pattern-via-yamlschemas

## Default

```yaml
gui:
  # stuff relating to the UI
  windowSize: 'normal' # one of 'normal' | 'half' | 'full' default is 'normal'
  scrollHeight: 2 # how many lines you scroll by
  scrollPastBottom: true # enable scrolling past the bottom
  scrollOffMargin: 2 # how many lines to keep before/after the cursor when it reaches the top/bottom of the view; see 'Scroll-off Margin' section below
  scrollOffBehavior: 'margin' # one of 'margin' | 'jump'; see 'Scroll-off Margin' section below
  sidePanelWidth: 0.3333 # number from 0 to 1
  expandFocusedSidePanel: false
  mainPanelSplitMode: 'flexible' # one of 'horizontal' | 'flexible' | 'vertical'
  enlargedSideViewLocation: 'left' # one of 'left' | 'top'
  language: 'auto' # one of 'auto' | 'en' | 'zh-CN' | 'zh-TW' | 'pl' | 'nl' | 'ja' | 'ko' | 'ru'
  timeFormat: '02 Jan 06' # https://pkg.go.dev/time#Time.Format
  shortTimeFormat: '3:04PM'
  theme:
    activeBorderColor:
      - green
      - bold
    inactiveBorderColor:
      - white
    searchingActiveBorderColor:
      - cyan
      - bold
    optionsTextColor:
      - blue
    selectedLineBgColor:
      - blue # set to `default` to have no background colour
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
  skipDiscardChangeWarning: false
  skipStashWarning: false
  showFileTree: true # for rendering changes files in a tree format
  showListFooter: true # for seeing the '5 of 20' message in list panels
  showRandomTip: true
  showBranchCommitHash: false # show commit hashes alongside branch names
  showBottomLine: true # for hiding the bottom information line (unless it has important information to tell you)
  showPanelJumps: true # for showing the jump-to-panel keybindings as panel subtitles
  showCommandLog: true
  showIcons: false # deprecated: use nerdFontsVersion instead
  nerdFontsVersion: "" # nerd fonts version to use ("2" or "3"); empty means don't show nerd font icons
  showFileIcons: true # for hiding file icons in the file views
  commandLogSize: 8
  splitDiff: 'auto' # one of 'auto' | 'always'
  skipRewordInEditorWarning: false # for skipping the confirmation before launching the reword editor
  border: 'rounded' # one of 'single' | 'double' | 'rounded' | 'hidden'
  animateExplosion: true # shows an explosion animation when nuking the working tree
  portraitMode: 'auto' # one of 'auto' | 'never' | 'always'
  filterMode: 'substring' # one of 'substring' | 'fuzzy'; see 'Filtering' section below
git:
  paging:
    colorArg: always
    useConfig: false
  commit:
    signOff: false
    autoWrapCommitMessage: true # automatic WYSIWYG wrapping of the commit message as you type
    autoWrapWidth: 72 # if autoWrapCommitMessage is true, the width to wrap to
  merging:
    # only applicable to unix users
    manualCommit: false
    # extra args passed to `git merge`, e.g. --no-ff
    args: ''
  log:
    # one of date-order, author-date-order, topo-order or default.
    # topo-order makes it easier to read the git log graph, but commits may not
    # appear chronologically. See https://git-scm.com/docs/git-log#_commit_ordering
    #
    # Deprecated: Configure this with `Log menu -> Commit sort order` (<c-l> in the commits window by default).
    order: 'topo-order'
    # one of always, never, when-maximised
    # this determines whether the git graph is rendered in the commits panel
    #
    # Deprecated: Configure this with `Log menu -> Show git graph` (<c-l> in the commits window by default).
    showGraph: 'always'
    # displays the whole git graph by default in the commits panel (equivalent to passing the `--all` argument to `git log`)
    showWholeGraph: false
  skipHookPrefix: WIP
  # The main branches. We colour commits green if they belong to one of these branches,
  # so that you can easily see which commits are unique to your branch (coloured in yellow)
  mainBranches: [master, main]
  autoFetch: true
  autoRefresh: true
  fetchAll: true # Pass --all flag when running git fetch. Set to false to fetch only origin (or the current branch's upstream remote if there is one)
  branchLogCmd: 'git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --'
  allBranchesLogCmd: 'git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium'
  overrideGpg: false # prevents lazygit from spawning a separate process when using GPG
  disableForcePushing: false
  parseEmoji: false
  truncateCopiedCommitHashesTo: 12 # When copying commit hashes to the clipboard, truncate them to this length. Set to 40 to disable truncation.
os:
  copyToClipboardCmd: '' # See 'Custom Command for Copying to Clipboard' section
  editPreset: '' # see 'Configuring File Editing' section
  edit: ''
  editAtLine: ''
  editAtLineAndWait: ''
  open: ''
  openLink: ''
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
    optionMenu: <disabled> # show help menu
    optionMenu-alt1: '?' # show help menu
    select: '<space>'
    goInto: '<enter>'
    openRecentRepos: '<c-r>'
    confirm: '<enter>'
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
    extrasMenu: '@'
    toggleWhitespaceInDiffView: '<c-w>'
    increaseContextInDiffView: '}'
    decreaseContextInDiffView: '{'
    toggleRangeSelect: 'v'
    rangeSelectUp: '<s-up>'
    rangeSelectDown: '<s-down>'
  status:
    checkForUpdate: 'u'
    recentRepos: '<enter>'
  files:
    commitChanges: 'c'
    commitChangesWithoutHook: 'w' # commit changes without pre-commit hook
    amendLastCommit: 'A'
    commitChangesWithEditor: 'C'
    findBaseCommitForFixup: '<c-f>'
    confirmDiscard: 'x'
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
    createTag: 'T'
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
    amendAttributeMenu: 'a'
    pickCommit: 'p' # pick commit (when mid-rebase)
    revertCommit: 't'
    cherryPickCopy: 'C'
    pasteCommits: 'V'
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
    toggleSelectHunk: 'a'
    pickBothHunks: 'b'
  submodules:
    init: 'i'
    update: 'u'
    bulkMenu: 'b'
  commitMessage:
    commitMenu: '<c-o>'
  amendAttribute:
    addCoAuthor: 'c'
    resetAuthor: 'a'
    setAuthor: 'A'
```

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

## Custom Command for Copying to Clipboard
```yaml
os:
  copyToClipboardCmd: ''
```
Specify an external command to invoke when copying to clipboard is requested. `{{text}` will be replaced by text to be copied. Default is to copy to system clipboard.

If you are working on a terminal that supports OSC52, the following command will let you take advantage of it:
```
os:
  copyToClipboardCmd: printf "\033]52;c;$(printf {{text}} | base64)\a" > /dev/tty
```


## Configuring File Editing

There are two commands for opening files, `o` for "open" and `e` for "edit". `o` acts as if the file was double-clicked in the Finder/Explorer, so it also works for non-text files, whereas `e` opens the file in an editor. `e` can also jump to the right line in the file if you invoke it from the staging panel, for example.

To tell lazygit which editor to use for the `e` command, the easiest way to do that is to provide an editPreset config, e.g.

```yaml
os:
  editPreset: 'vscode'
```

Supported presets are `vim`, `nvim`, `nvim-remote`, `lvim`, `emacs`, `nano`, `micro`, `vscode`, `sublime`, `bbedit`, `kakoune`, `helix`, and `xcode`. In many cases lazygit will be able to guess the right preset from your $(git config core.editor), or an environment variable such as $VISUAL or $EDITOR.

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
