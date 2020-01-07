# User Config:

Default path for the config file: `~/.config/jesseduffield/lazygit/config.yml`

## Default:

```yaml
  gui:
    # stuff relating to the UI
    scrollHeight: 2 # how many lines you scroll by
    scrollPastBottom: true # enable scrolling past the bottom
    theme:
      lightTheme: false # For terminals with a light background
      activeBorderColor:
        - white
        - bold
      inactiveBorderColor:
        - white
      optionsTextColor:
        - blue
    commitLength:
      show: true
    mouseEvents: true
  git:
    merging:
      # only applicable to unix users
      manualCommit: false
    skipHookPrefix: WIP
    autoFetch: true
  update:
    method: prompt # can be: prompt | background | never
    days: 14 # how often an update is checked for
  reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
  confirmOnQuit: false
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
      prevBlock: '<left>' # goto the previous block / panel
      nextBlock: '<right>' # goto the next block / panel
      prevBlock-alt: 'h' # goto the previous block / panel
      nextBlock-alt: 'l' # goto the next block / panel
      optionMenu: 'x' # show help menu
      optionMenu-alt1: '?' # show help menu
      select: '<space>'
      goInto: '<enter>'
      remove: 'd'
      new: 'n'
      edit: 'e'
      openFile: 'o'
      scrollUpMain: '<pgup>' # main panel scrool up
      scrollDownMain: '<pgdown>' # main panel scrool down
      scrollUpMain-alt1: 'K' # main panel scrool up
      scrollDownMain-alt1: 'J' # main panel scrool down
      scrollUpMain-alt2: '<c-u>' # main panel scrool up
      scrollDownMain-alt2: '<c-d>' # main panel scrool down
      executeCustomCommand: 'X'
      createRebaseOptionsMenu: 'm'
      pushFiles: 'P'
      pullFiles: 'p'
      refresh: 'R'
      createPatchOptionsMenu: '<c-p>'
      nextBranchTab: ']'
      prevBranchTab: '['
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
    branches:
      createPullRequest: 'o'
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
      toggleDiffCommit: 'h'
      checkoutCommit: '<space>'
    stash:
      popStash: 'g'
    commitFiles:
      checkoutCommitFile: 'c'
    main:
      toggleDragSelect: 'v'
      toggleDragSelect-alt: 'V'
      toggleSelectHunk: 'a'
      pickBothHunks: 'b'
      undo: 'z'
```

## Platform Defaults:

### Windows:

```yaml
  os:
    openCommand: 'cmd /c "start "" {{filename}}"'
```

### Linux:

```yaml
  os:
    openCommand: 'sh -c "xdg-open {{filename}} >/dev/null"'
```

### OSX:

```yaml
  os:
    openCommand: 'open {{filename}}'
```

### Recommended Config Values:

for users of VSCode

```yaml
  os:
    openCommand: 'code -r {{filename}}'
```

## Color Attributes:

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

## Light terminal theme:

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
```

## Example Coloring:

![border example](/docs/resources/colored-border-example.png)

## Keybindings:
For all possible keybinding options, check [Custom_Keybinding.md](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybinding.md) <++>


#### Example Keybindings For Colemak Users:
```yaml
  keybinding:
    universal:
      prevItem-alt: 'u'
      nextItem-alt: 'e'
      prevBlock-alt: 'n'
      nextBlock-alt: 'i'
      new: 'k'
      edit: 'o'
      openFile: 'O'
      scrollUpMain-alt1: 'U'
      scrollDownMain-alt1: 'E'
      scrollUpMain-alt2: '<c-u>'
      scrollDownMain-alt2: '<c-e>'
    files:
      ignoreFile: 'I'
    commits:
      moveDownCommit: '<c-e>'
      moveUpCommit: '<c-u>'
```

