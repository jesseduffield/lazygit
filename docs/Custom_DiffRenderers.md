# Custom Diff Renderers

Custom diff renderers are useful for showing a better rendering of a diff than git's builtin raw diff, and using one is strongly recommended (I personally prefer delta myself, but that's a matter of personal preference). There are three types of diff renderers that lazygit supports:

- **stdin filters**, e.g. [delta](#delta) and [diff-so-fancy](#diff-so-fancy). They take git's raw output as stdin and produce something nicer on stdout, and they are hooked up using git's GIT_PAGER mechanism. (These used to be called "custom pagers" in earlier lazygit versions.)
- **external diff programs**, e.g. difftastic; these are called using git's `--ext-diff` flag, and they take over diff generation from git completely rather than post-processing git's output.
- **git's raw output using custom arguments**; mainly useful for `--color-words` (or `--word-diff` if you are color blind).

Diff renderers are configured with the `diffRenderers` array in the `git` section of lazygit's config file; it is an array because you can have multiple entries that you can cycle through with the `|` key. This can be useful if you usually prefer a particular diff renderer, but want to use a different one for certain kinds of diffs.

Fields that are shared by all renderer types:

- **type** The type of diff renderer; choices are `stdinFilter`, `extDiff`, or `rawGit`. `stdinFilter` is the default, because it's the most common one; so you can omit this if you use delta.
- **name** A name that is shown in the status bar toast when cycling renderers; defaults to the first word of the renderer command, but can be useful e.g. to distinguish "delta" from "delta side-by-side" if you have entries for both.

Fields only for `stdinFilter`:

- **command** The command line to use for `GIT_PAGER`.

- **colorArg** whether you want the `--color=always` arg in your `git diff` command. Some diff renderers want it set to `always`, others want it set to `never`. The default is `always`, since that's what most renderers need.

Fields only for `extDiff`:

- **command** The command line to use for the `diff.external` git config. If left empty, it uses the global value of git's `diff.external` config; this can be useful if you also want to use it for diffs on the command line, and it also has the advantage that you can configure it per file type in `.gitattributes`; see https://git-scm.com/docs/gitattributes#_defining_an_external_diff_driver.

  You can include the `{{diffContext}}` template variable to pass lazygit's current diff context size (the value controlled by the `{`/`}` keybindings) to the diff tool.

Fields only for `rawGit`:

- **args** The additional arguments to use in the `git diff` or `git show` call (e.g. `--color-words`)

Here's an example for a multi-renderer setup:

```yaml
git:
  diffRenderers:
    - command: delta --dark --paging=never
    - command: ydiff -p cat
      colorArg: never
    - type: extDiff
      command: difft --color=always --context={{diffContext}}
    - type: rawGit
      args: --color-words
      name: color-words
    - type: rawGit # git's default diff
      name: default
```

## Delta:

```yaml
git:
  diffRenderers:
    - command: delta --dark --paging=never
```

![](https://i.imgur.com/QJpQkF3.png)

A cool feature of delta is --hyperlinks, which renders clickable links for the line numbers in the left margin, and lazygit supports these. To use them, set the `command:` field to `delta --dark --paging=never --line-numbers --hyperlinks --hyperlinks-file-link-format="lazygit-edit://{path}:{line}"`; this allows you to click on an underlined line number in the diff to jump right to that same line in your editor.

Note that delta's `--navigate` option doesn't work in lazygit, for technical reasons.

## Diff-so-fancy

```yaml
git:
  diffRenderers:
    - command: diff-so-fancy
```

![](https://i.imgur.com/rjH1TpT.png)

## ydiff

```yaml
gui:
  sidePanelWidth: 0.2 # gives you more space to show things side-by-side
git:
  diffRenderers:
    - colorArg: never
      command: ydiff -p cat
```

![](https://i.imgur.com/vaa8z0H.png)
