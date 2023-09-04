# Custom Pagers

Lazygit supports custom pagers, [configured](/docs/Config.md) in the config.yml file (which can be opened by pressing `e` in the Status panel).

Support does not extend to Windows users, because we're making use of a package which doesn't have Windows support.

## Default:

```yaml
git:
  paging:
    colorArg: always
    useConfig: false
```

the `colorArg` key is for whether you want the `--color=always` arg in your `git diff` command. Some pagers want it set to `always`, others want it set to `never`.

## Delta:

```yaml
git:
  paging:
    colorArg: always
    pager: delta --dark --paging=never
```

![](https://i.imgur.com/QJpQkF3.png)

## Diff-so-fancy

```yaml
git:
  paging:
    colorArg: always
    pager: diff-so-fancy
```

![](https://i.imgur.com/rjH1TpT.png)

## ydiff

```yaml
gui:
  sidePanelWidth: 0.2 # gives you more space to show things side-by-side
git:
  paging:
    colorArg: never
    pager: ydiff -p cat -s --wrap --width={{columnWidth}}
```

![](https://i.imgur.com/vaa8z0H.png)

Be careful with this one, I think the homebrew and pip versions are behind master. I needed to directly download the ydiff script to get the no-pager functionality working.

## Using git config

```yaml
git:
  paging:
    colorArg: always
    useConfig: true
```

If you set `useConfig: true`, lazygit will use whatever pager is specified in `$GIT_PAGER`, `$PAGER`, or your *git config*. If the pager ends with something like ` | less` we will strip that part out, because less doesn't play nice with our rendering approach. If the custom pager uses less under the hood, that will also break rendering (hence the `--paging=never` flag for the `delta` pager).

## Using external diff commands

Some diff tools can't work as a simple pager like the ones above do, because they need access to the entire diff, so just post-processing git's diff is not enough for them. The most notable example is probably [difftastic](https://difftastic.wilfred.me.uk).

These can be used in lazygit by using the `externalDiffCommand` config; in the case of difftastic, that could be

```yaml
git:
  paging:
    externalDiffCommand: difft --color=always
```

The `colorArg`, `pager`, and `useConfig` options are not used in this case.

You can add whatever extra arguments you prefer for your difftool; for instance

```yaml
git:
  paging:
    externalDiffCommand: difft --color=always --display=inline --syntax-highlight=off
```
