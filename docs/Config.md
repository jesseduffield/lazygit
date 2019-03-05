# User Config:

## Default:

```
  gui:
    # stuff relating to the UI
    scrollHeight: 2 # how many lines you scroll by
    scrollPastBottom: true # enable scrolling past the bottom
    theme:
      activeBorderColor:
        - white
        - bold
      inactiveBorderColor:
        - white
      optionsTextColor:
        - blue
    commitLength:
      show: true
  git:
    merging:
      # only applicable to unix users
      manualCommit: false
  update:
    method: prompt # can be: prompt | background | never
    days: 14 # how often an update is checked for
  reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
  confirmOnQuit: false
```

## Platform Defaults:

### Windows:

```
  os:
    openCommand: 'cmd /c "start "" {{filename}}"'
```

### Linux:

```
  os:
    openCommand: 'sh -c "xdg-open {{filename}} >/dev/null"'
```

### OSX:

```
  os:
    openCommand: 'open {{filename}}'
```

### Recommended Config Values:

for users of VSCode

```
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

## Example Coloring:

![border example](/docs/resources/colored-border-example.png)
