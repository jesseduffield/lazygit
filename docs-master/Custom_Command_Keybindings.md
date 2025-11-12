# Custom Command Keybindings

You can add custom command keybindings in your config.yml (accessible by pressing 'e' on the status panel from within lazygit) like so:

```yml
customCommands:
  - key: '<c-r>'
    context: 'commits'
    command: 'hub browse -- "commit/{{.SelectedLocalCommit.Hash}}"'
  - key: 'a'
    context: 'files'
    command: "git {{if .SelectedFile.HasUnstagedChanges}} add {{else}} reset {{end}} {{.SelectedFile.Name | quote}}"
    description: 'Toggle file staged'
  - key: 'C'
    context: 'global'
    command: "git commit"
    output: terminal
  - key: 'n'
    context: 'localBranches'
    prompts:
      - type: 'menu'
        title: 'What kind of branch is it?'
        key: 'BranchType'
        options:
          - name: 'feature'
            description: 'a feature branch'
            value: 'feature'
          - name: 'hotfix'
            description: 'a hotfix branch'
            value: 'hotfix'
          - name: 'release'
            description: 'a release branch'
            value: 'release'
      - type: 'input'
        title: 'What is the new branch name?'
        key: 'BranchName'
        initialValue: ''
    command: "git flow {{.Form.BranchType}} start {{.Form.BranchName}}"
    loadingText: 'Creating branch'
```

Looking at the command assigned to the 'n' key, here's what the result looks like:

![](../../assets/custom-command-keybindings.gif)

Custom command keybindings will appear alongside inbuilt keybindings when you view the keybindings menu by pressing '?':

![](https://i.imgur.com/QB21FPx.png)

For a given custom command, here are the allowed fields:
| _field_ | _description_ | required |
|-----------------|----------------------|-|
| key | The key to trigger the command. Use a single letter or one of the values from [here](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md). Custom commands without a key specified can be triggered by selecting them from the keybindings (`?`) menu | no |
| command | The command to run (using Go template syntax for placeholder values) | yes |
| context | The context in which to listen for the key (see [below](#contexts)) | yes |
| prompts | A list of prompts that will request user input before running the final command | no |
| loadingText | Text to display while waiting for command to finish | no |
| description | Label for the custom command when displayed in the keybindings menu | no |
| output | Where the output of the command should go. 'none' discards it, 'terminal' suspends lazygit and runs the command in the terminal (useful for commands that require user input), 'log' streams it to the command log, 'logWithPty' is like 'log' but runs the command in a pseudo terminal (can be useful for commands that produce colored output when the output is a terminal), and 'popup' shows it in a popup. | no |
| outputTitle | The title to display in the popup panel if output is set to 'popup'. If left unset, the command will be used as the title. | no |
| after | Actions to take after the command has completed | no |

Here are the options for the `after` key:
| _field_ | _description_ | required |
|-----------------|----------------------|-|
| checkForConflicts | true/false. If true, check for merge conflicts | no |

## Contexts

The permitted contexts are:

| _context_      | _description_                                                                                            |
| -------------- | -------------------------------------------------------------------------------------------------------- |
| status         | The 'Status' tab                                                                                         |
| files          | The 'Files' tab                                                                                          |
| worktrees      | The 'Worktrees' tab                                                                                      |
| submodules     | The 'Submodules' tab                                                                                     |
| localBranches  | The 'Local Branches' tab                                                                                 |
| remotes        | The 'Remotes' tab                                                                                        |
| remoteBranches | The context you get when pressing enter on a remote in the remotes tab                                   |
| tags           | The 'Tags' tab                                                                                           |
| commits        | The 'Commits' tab                                                                                        |
| reflogCommits  | The 'Reflog' tab                                                                                         |
| subCommits     | The context you see when pressing enter on a branch                                                      |
| commitFiles    | The context you see when pressing enter on a commit or stash entry (warning, might be renamed in future) |
| stash          | The 'Stash' tab                                                                                          |
| global         | This keybinding will take affect everywhere                                                              |

> **Bonus**
>
> You can use a comma-separated string, such as `context: 'commits, subCommits'`, to make it effective in multiple contexts.


## Prompts

### Common fields

These fields are applicable to all prompts.

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| type              | One of 'input', 'confirm', 'menu', 'menuFromCommand'                                                           | yes        |
| title             | The title to display in the popup panel                                                        | no         |
| key | Used to reference the entered value from within the custom command. E.g. a prompt with `key: 'Branch'` can be referred to as `{{.Form.Branch}}` in the command | yes |

### Input

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| initialValue      | The initial value to appear in the text box               | no         |
| suggestions       | Shows suggestions as the input is entered. See below for details                                                          | no         |

The permitted suggestions fields are:
| _field_ | _description_ | _required_ |
|-----------------|----------------------|-|
| preset | Uses built-in logic to obtain the suggestions. One of 'authors', 'branches', 'files', 'refs', 'remotes', 'remoteBranches', 'tags' | no |
| command | Command to run such that each line in the output becomes a suggestion. Mutually exclusive with 'preset' field. | no |

Here's an example of passing a preset:

```yml
customCommands:
  - key: 'a'
    command: 'echo {{.Form.Branch | quote}}'
    context: 'commits'
    prompts:
      - type: 'input'
        title: 'Which branch?'
        key: 'Branch'
        suggestions:
          preset: 'branches' # use built-in logic for obtaining branches
```

Here's an example of passing a command directly:

```yml
customCommands:
  - key: 'a'
    command: 'echo {{.Form.Branch | quote}}'
    context: 'commits'
    prompts:
      - type: 'input'
        title: 'Which branch?'
        key: 'Branch'
        suggestions:
          command: "git branch --format='%(refname:short)'"
```


Here's an example of passing an initial value for the input:

```yml
customCommands:
  - key: 'a'
    command: 'echo {{.Form.Remote | quote}}'
    context: 'commits'
    prompts:
    - type: 'input'
      title: 'Remote:'
      key: 'Remote'
      initialValue: "{{.SelectedRemote.Name}}"
```

### Confirm

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| body              | The immutable body text to appear in the text box       | no         |

Example:

```yml
customCommands:
  - key: 'a'
    command: 'echo "pushing to remote"'
    context: 'commits'
    prompts:
    - type: 'confirm'
      title: 'Push to remote'
      body: 'Are you sure you want to push to the remote?'
```

### Menu

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| options           | The options to display in the menu                         | yes         |

The permitted option fields are:
| _field_ | _description_ | _required_ |
|-----------------|----------------------|-|
| name | The first part of the label | no |
| description | The second part of the label | no |
| value | the value that will be used in the command | yes |

If an option has no name the value will be displayed to the user in place of the name, so you're allowed to only include the value like so:

```yml
customCommands:
  - key: 'a'
    command: 'echo {{.Form.BranchType | quote}}'
    context: 'commits'
    prompts:
      - type: 'menu'
        title: 'What kind of branch is it?'
        key: 'BranchType'
        options:
          - value: 'feature'
          - value: 'hotfix'
          - value: 'release'
```

Here's an example of supplying more detail for each option:

```yml
customCommands:
  - key: 'a'
    command: 'echo {{.Form.BranchType | quote}}'
    context: 'commits'
    prompts:
      - type: 'menu'
        title: 'What kind of branch is it?'
        key: 'BranchType'
        options:
          - value: 'feature'
            name: 'feature branch'
            description: 'branch based off develop'
          - value: 'hotfix'
            name: 'hotfix branch'
            description: 'branch based off main for fast bug fixes'
          - value: 'release'
            name: 'release branch'
            description: 'branch for a release'
```

### Menu-from-command

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| command           | The command to run to generate menu options                  | yes        |
| filter            | The regexp to run specifying groups which are going to be kept from the command's output      | no        |
| valueFormat       | How to format matched groups from the filter to construct a menu item's value | no        |
| labelFormat       | Like valueFormat but for the labels. If `labelFormat` is not specified, `valueFormat` is shown instead. | no         |

Here's an example using named groups in the regex. Notice how we can pipe the label to a colour function for coloured output (available colours [here](https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md))

```yml
  - key : 'a'
    description: 'Checkout a remote branch as FETCH_HEAD'
    command: "git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"
    context: 'remotes'
    prompts:
      - type: 'menuFromCommand'
        title: 'Remote branch:'
        key: 'Branch'
        command: 'git branch  -r --list {{.SelectedRemote.Name }}/*'
        filter: '.*{{.SelectedRemote.Name }}/(?P<branch>.*)'
        valueFormat: '{{ .branch }}'
        labelFormat: '{{ .branch | green }}'
```

Here's an example using unnamed groups:

```yml
  - key : 'a'
    description: 'Checkout a remote branch as FETCH_HEAD'
    command: "git fetch {{.Form.Remote}} {{.Form.Branch}} && git checkout FETCH_HEAD"
    context: 'remotes'
    prompts:
      - type: 'menuFromCommand'
        title: 'Remote branch:'
        key: 'Branch'
        command: 'git branch  -r --list {{.SelectedRemote.Name }}/*'
        filter: '.*{{.SelectedRemote.Name }}/(.*)'
        valueFormat: '{{ .group_1 }}'
        labelFormat: '{{ .group_1 | green }}'
```

Here's an example using a command but not specifying anything else: so each line from the command becomes the value and label of the menu items

```yml
  - key : 'a'
    description: 'Checkout a remote branch as FETCH_HEAD'
    command: "open {{.Form.File | quote}}"
    context: 'global'
    prompts:
      - type: 'menuFromCommand'
        title: 'File:'
        key: 'File'
        command: 'ls'
```

## Placeholder values

Your commands can contain placeholder strings using Go's [template syntax](https://jan.newmarch.name/golang/template/chapter-template.html). The template syntax is pretty powerful, letting you do things like conditionals if you want, but for the most part you'll simply want to be accessing the fields on the following objects:

```
SelectedCommit
SelectedCommitRange
SelectedFile
SelectedPath
SelectedSubmodule
SelectedLocalBranch
SelectedRemoteBranch
SelectedRemote
SelectedTag
SelectedStashEntry
SelectedCommitFile
SelectedWorktree
CheckedOutBranch
```

(For legacy reasons, `SelectedLocalCommit`, `SelectedReflogCommit`, and `SelectedSubCommit` are also available, but they are deprecated.)


To see what fields are available on e.g. the `SelectedFile`, see [here](https://github.com/jesseduffield/lazygit/blob/master/pkg/gui/services/custom_commands/models.go) (all the modelling lives in the same file).

We don't support accessing all elements of a range selection yet. We might add this in the future, but as a special case you can access the range of selected commits by using `SelectedCommitRange`, which has two properties `.To` and `.From` which are the hashes of the bottom and top selected commits, respectively. This is useful for passing them to a git command that operates on a range of commits. For example, to create patches for all selected commits, you might use
```yml
  command: "git format-patch {{.SelectedCommitRange.From}}^..{{.SelectedCommitRange.To}}"
```

We support the following functions:

### Quoting

Quote wraps a string in quotes with necessary escaping for the current platform.

```
git {{.SelectedFile.Name | quote}}
```

### Running a command

Runs a command and returns the output. If the command outputs more than a single line, it will produce an error.

```
initialValue: "username/{{ runCommand "date +\"%Y/%-m\"" }}/"
```

## Keybinding collisions

If your custom keybinding collides with an inbuilt keybinding that is defined for the same context, only the custom keybinding will be executed. This also applies to the global context. However, one caveat is that if you have a custom keybinding defined on the global context for some key, and there is an in-built keybinding defined for the same key and for a specific context (say the 'files' context), then the in-built keybinding will take precedence. See how to change in-built keybindings [here](https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#keybindings)

## Menus of custom commands

For custom commands that are not used very frequently it may be preferable to hide them in a menu; you can assign a key to open the menu, and the commands will appear inside. This has the advantage that you don't have to come up with individual unique keybindings for all those commands that you don't use often; the keybindings for the commands in the menu only need to be unique within the menu. Here is an example:

```yml
customCommands:
- key: X
  description: "Copy/paste commits across repos"
  commandMenu:
  - key: c
    command: 'git format-patch --stdout {{.SelectedCommitRange.From}}^..{{.SelectedCommitRange.To}} | pbcopy'
    context: commits, subCommits
    description: "Copy selected commits to clipboard"
  - key: v
    command: 'pbpaste | git am'
    context: "commits"
    description: "Paste selected commits from clipboard"
```

If you use the commandMenu property, none of the other properties except key and description can be used.

## Debugging

If you want to verify that your command actually does what you expect, you can wrap it in an 'echo' call and set `output: popup` so that it doesn't actually execute the command but you can see how the placeholders were resolved.

## More Examples

See the [wiki](https://github.com/jesseduffield/lazygit/wiki/Custom-Commands-Compendium) page for more examples, and feel free to add your own custom commands to this page so others can benefit!
