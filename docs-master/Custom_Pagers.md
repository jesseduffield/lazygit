# Custom Pagers

Lazygit supports custom pagers, [configured](/docs/Config.md) in the config.yml file (which can be opened by pressing `e` in the Status panel).

Support does not extend to Windows users, because we're making use of a package which doesn't have Windows support. However, see [below](#emulating-custom-pagers-on-windows) for a workaround.

Multiple pagers are supported; you can cycle through them with the `|` key. This can be useful if you usually prefer a particular pager, but want to use a different one for certain kinds of diffs.

Pagers are configured with the `pagers` array in the git section; here's an example for a multi-pager setup:

```yaml
git:
  pagers:
    - pager: delta --dark --paging=never
    - pager: ydiff -p cat -s --wrap --width={{columnWidth}}
      colorArg: never
    - externalDiffCommand: difft --color=always
```

The `colorArg` key is for whether you want the `--color=always` arg in your `git diff` command. Some pagers want it set to `always`, others want it set to `never`. The default is `always`, since that's what most pagers need.

## Delta:

```yaml
git:
  pagers:
    - pager: delta --dark --paging=never
```

![](https://i.imgur.com/QJpQkF3.png)

A cool feature of delta is --hyperlinks, which renders clickable links for the line numbers in the left margin, and lazygit supports these. To use them, set the `pager:` config to `delta --dark --paging=never --line-numbers --hyperlinks --hyperlinks-file-link-format="lazygit-edit://{path}:{line}"`; this allows you to click on an underlined line number in the diff to jump right to that same line in your editor.

## Diff-so-fancy

```yaml
git:
  pagers:
    - pager: diff-so-fancy
```

![](https://i.imgur.com/rjH1TpT.png)

## ydiff

```yaml
gui:
  sidePanelWidth: 0.2 # gives you more space to show things side-by-side
git:
  pagers:
    - colorArg: never
      pager: ydiff -p cat -s --wrap --width={{columnWidth}}
```

![](https://i.imgur.com/vaa8z0H.png)

Be careful with this one, I think the homebrew and pip versions are behind master. I needed to directly download the ydiff script to get the no-pager functionality working.

## Using external diff commands

Some diff tools can't work as a simple pager like the ones above do, because they need access to the entire diff, so just post-processing git's diff is not enough for them. The most notable example is probably [difftastic](https://difftastic.wilfred.me.uk).

These can be used in lazygit by using the `externalDiffCommand` config; in the case of difftastic, that could be

```yaml
git:
  pagers:
    - externalDiffCommand: difft --color=always
```

The `colorArg` and `pager` options are not used in this case.

You can add whatever extra arguments you prefer for your difftool; for instance

```yaml
git:
  pagers:
    - externalDiffCommand: difft --color=always --display=inline --syntax-highlight=off
```

Instead of setting this command in lazygit's `externalDiffCommand` config, you can also tell lazygit to use the external diff command that is configured in git itself (`diff.external`), by using

```yaml
git:
  pagers:
    - useExternalDiffGitConfig: true
```

This can be useful if you also want to use it for diffs on the command line, and it also has the advantage that you can configure it per file type in `.gitattributes`; see https://git-scm.com/docs/gitattributes#_defining_an_external_diff_driver.

## Emulating custom pagers on Windows

There is a trick to emulate custom pagers on Windows using a Powershell script configured as an external diff command. It's not perfect, but certainly better than nothing. To do this, save the following script as `lazygit-pager.ps1` at a convenient place on your disk:

```pwsh
#!/usr/bin/env pwsh

$old = $args[1].Replace('\', '/')
$new = $args[4].Replace('\', '/')
$path = $args[0]
git diff --no-index --no-ext-diff $old $new
  | %{ $_.Replace($old, $path).Replace($new, $path) }
  | delta --width=$env:LAZYGIT_COLUMNS
```

Use the pager of your choice with the arguments you like in the last line of the script. Personally I wouldn't want to use lazygit anymore without delta's `--hyperlinks --hyperlinks-file-link-format="lazygit-edit://{path}:{line}"` args, see [above](#delta).

In your lazygit config, use

```yml
git:
  pagers:
    - externalDiffCommand: "C:/wherever/lazygit-pager.ps1"
```

The main limitation of this approach compared to a "real" pager is that renames are not displayed correctly; they are shown as if they were modifications of the old file. (This affects only the hunk headers; the diff itself is always correct.)
