# User Config:

## Default:

```
  gui:
    # stuff relating to the UI
    scrollHeight: 2 # how many lines you scroll by
    theme:
      activeBorderColor:
        - white
        - bold
      inactiveBorderColor:
        - white
      optionsTextColor:
        - blue
  git:
    # stuff relating to git
  os:
    # stuff relating to the OS
  update:
    method: prompt # can be: prompt | background | never
    days: 14 # how often an update is checked for
  reporting: 'undetermined' # one of: 'on' | 'off' | 'undetermined'
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
