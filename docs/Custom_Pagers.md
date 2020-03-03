# Custom Pagers

Lazygit supports custom pagers, [configured](/docs/Config.md) in the config.yml file (which can be opened by pressing 'o' in the Status panel).

Support does not extend to windows users, because we're making use of a package which doesn't have windows support.

## Default:

```yaml
git:
  paging:
    colorArg: always
    useConfig: false
```

the `colorArg` key is for whether you want the `--color=always` arg in your `git diff` command. Some pagers want it set to always, others want it set to 'never'.

## Delta:

```yaml
git:
  paging:
    colorArg: always
    pager: delta --dark --paging=never --24-bit-color=never
```

![](https://i.imgur.com/A4C3xiZ.png)

## Diff-so-fancy

```yaml
git:
  paging:
    colorArg: always
    pager: diff-so-fancy
```

![](https://i.imgur.com/rjH1TpT.png)

## Using git config

```yaml
git:
  paging:
    colorArg: always
    useConfig: true
```

If you set `useConfig: true`, lazygit will use whatever pager is specified in $GIT_PAGER, $PAGER, or your git config. If the pager ends with something like ' | less' we will strip that part out, because less doesn't play nice with our rendering approach. If the custom pager uses less under the hood, that will also break rendering (hence the `--paging=never` flag for the `delta` pager).
