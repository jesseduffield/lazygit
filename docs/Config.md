# User Config:

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
      createRebaseOptionsMenu: 'm'
      pushFiles: 'P'
      pullFiles: 'p'
      refresh: 'R'
      createPatchOptionsMenu: '<c-p>'
    status:
      checkForUpdate: 'u'
      recentRepos: 's'
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
      executeCustomCommand: 'X'
    branches:
      createPullRequest: 'o'
      checkoutBranchesByName: 'c'
      forceCheckoutBranch: 'F'
      rebaseBranch: 'r'
      mergeIntoCurrentBranch: 'M'
      FastForward: 'f' # fast-forward this branch from its upstream
      pushTag: 'P'
      nextBranchTab: ']'
      prevBranchTab: '['
      setUpstream: 'u' # set as upstream of checked-out branch
    commits:
      squashDown: 's'
      renameCommit: 'r'
      renameCommitWithEditor: 'R'
      resetToThisCommit: 'g'
      fixupCommit: 'f'
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
      viewCommitFiles: '<enter>'
      tagCommit: 'T'
    stash:
      popStash: 'g'
    commitFiles:
      checkoutCommitFile: 'c'
    main:
      toggleDragSelect: 'v'
      toggleDragSelect-alt: 'V'
      toggleSelectHunk: 'a'
      PickBothHunks: 'b'
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

## Example Keybindings For Colemak Users:
```yaml
  keybinding:
    universal:
      prevItem-alt: 'u' # go one line up
      nextItem-alt: 'e' # go one line down
      prevBlock-alt: 'n' # goto the previous block / panel
      nextBlock-alt: 'i' # goto the next block / panel
      new: 'k'
      edit: 'o'
      openFile: 'O'
      scrollUpMain-alt1: 'U' # main panel scrool up
      scrollDownMain-alt1: 'E' # main panel scrool down
      scrollDownMain-alt2: '<c-e>' # main panel scrool down
    status:
      checkForUpdate: '<c-u>'
    files:
      ignoreFile: 'I'
    commits:
      moveDownCommit: '<c-e>' # move commit down one
      moveUpCommit: '<c-u>' # move commit up one
```
