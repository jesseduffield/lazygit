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
