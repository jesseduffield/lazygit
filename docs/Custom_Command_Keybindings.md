# Custom Command Keybindings

You can add custom command keybindings in your config.yml (accessible by pressing 'o' on the status panel from within lazygit) like so:

```yml
customCommands:
  - key: '<c-r>'
    command: 'hub browse -- "commit/{{.SelectedLocalCommit.Sha}}"'
    context: 'commits'
  - key: 'a'
    command: "git {{if .SelectedFile.HasUnstagedChanges}} add {{else}} reset {{end}} {{.SelectedFile.Name | quote}}"
    context: 'files'
    description: 'toggle file staged'
  - key: 'C'
    command: "git commit"
    context: 'global'
    subprocess: true
  - key: 'n'
    prompts:
      - type: 'menu'
        title: 'What kind of branch is it?'
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
        initialValue: ''
    command: "git flow {{index .PromptResponses 0}} start {{index .PromptResponses 1}}"
    context: 'localBranches'
    loadingText: 'creating branch'
  - key : 'r'
    description: 'Checkout a remote branch as FETCH_HEAD'
    command: "git fetch {{index .PromptResponses 0}} {{index .PromptResponses 1}} && git checkout FETCH_HEAD"
    context: 'remotes'
    prompts:
      - type: 'input'
        title: 'Remote:'
        initialValue: "{{index .SelectedRemote.Name }}"
      - type: 'menuFromCommand'
        title: 'Remote branch:'
        command: 'git branch  -r --list {{index .PromptResponses 0}}/*'
        filter: '.*{{index .PromptResponses 0}}/(?P<branch>.*)'
        valueFormat: '{{ .branch }}'
        labelFormat: '{{ .branch | green }}'
  - key: '<f1>'
    command: 'git reset --soft {{.CheckedOutBranch.UpstreamRemote}}'
    context: 'files'
    prompts:
      - type: 'confirm'
        title: "Confirm:"
        body: "Are you sure you want to reset HEAD to {{.CheckedOutBranch.UpstreamRemote}}?"
```

Looking at the command assigned to the 'n' key, here's what the result looks like:

![](../../assets/custom-command-keybindings.gif)

Custom command keybindings will appear alongside inbuilt keybindings when you view the options menu by pressing 'x':

![](https://i.imgur.com/QB21FPx.png)

For a given custom command, here are the allowed fields:
| _field_ | _description_ | required |
|-----------------|----------------------|-|
| key | the key to trigger the command. Use a single letter or one of the values from [here](https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md) | yes |
| command | the command to run | yes |
| context | the context in which to listen for the key (see below) | yes |
| subprocess | whether you want the command to run in a subprocess (necessary if you want to view the output of the command or provide user input) | no |
| prompts | a list of prompts that will request user input before running the final command | no |
| loadingText | text to display while waiting for command to finish | no |
| description | text to display in the keybindings menu that appears when you press 'x' | no |
| stream | whether you want to stream the command's output to the Command Log panel | no |
| showOutput | whether you want to show the command's output in a gui prompt | no |

### Contexts

The permitted contexts are:

| _context_      | _description_                                                                                            |
| -------------- | -------------------------------------------------------------------------------------------------------- |
| status         | the 'Status' tab                                                                                         |
| files          | the 'Files' tab                                                                                          |
| localBranches  | the 'Local Branches' tab                                                                                 |
| remotes        | the 'Remotes' tab                                                                                        |
| remoteBranches | the context you get when pressing enter on a remote in the remotes tab                                   |
| tags           | the 'Tags' tab                                                                                           |
| commits        | the 'Commits' tab                                                                                        |
| reflogCommits  | the 'Reflog' tab                                                                                         |
| subCommits     | the context you see when pressing enter on a branch                                                      |
| commitFiles    | the context you see when pressing enter on a commit or stash entry (warning, might be renamed in future) |
| stash          | the 'Stash' tab                                                                                          |
| global         | this keybinding will take affect everywhere                                                              |

### Prompts

The permitted prompt fields are:

| _field_           | _description_                                                                                  | _required_ |
| ------------      | -----------------------------------------------------------------------------------------------| ---------- |
| type              | one of 'input', 'menu', or 'confirm'                                                           | yes        |
| title             | the title to display in the popup panel                                                        | no         |
| initialValue      | (only applicable to 'input' prompts) the initial value to appear in the text box               | no         |
| body              | (only applicable to 'confirm' prompts) the immutable body text to appear in the text box       | no         |
| options           | (only applicable to 'menu' prompts) the options to display in the menu                         | no         |
| command           | (only applicable to 'menuFromCommand' prompts) the command to run to generate                  | yes        |
|                   | menu options                                                                                   |            |
| filter            | (only applicable to 'menuFromCommand' prompts) the regexp to run specifying groups which are going to be kept from the command's output      | yes        |
| valueFormat       | (only applicable to 'menuFromCommand' prompts) how to format matched groups from the filter to construct a menu item's value (What gets appended to prompt responses when the item is selected). You can use named groups, or `{{ .group_GROUPID }}`. PS: named groups keep first match only | yes        |
| labelFormat       | (only applicable to 'menuFromCommand' prompts) how to format matched groups from the filter to construct the item's label (What's shown on screen). You can use named groups, or `{{ .group_GROUPID }}`. You can also color each match with `{{ .group_GROUPID \| colorname }}` (Color names from [here](https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md)). If `labelFormat` is not specified, `valueFormat` is shown instead. PS: named groups keep first match only | no         |

The permitted option fields are:
| _field_ | _description_ | _required_ |
|-----------------|----------------------|-|
| name | the string which will appear first on the line | no |
| description | the string which will appear second on the line | no |
| value | the value that will be stored in `.PromptResponses` if the option is selected | yes |

If an option has no name the value will be displayed to the user in place of the name, so you're allowed to only include the value like so:

```yml
    prompts:
      - type: 'menu'
        title: 'What kind of branch is it?'
        options:
          - value: 'feature'
          - value: 'hotfix'
          - value: 'release'
```

### Placeholder values

Your commands can contain placeholder strings using Go's [template syntax](https://jan.newmarch.name/golang/template/chapter-template.html). The template syntax is pretty powerful, letting you do things like conditionals if you want, but for the most part you'll simply want to be accessing the fields on the following objects:

```
SelectedLocalCommit
SelectedReflogCommit
SelectedSubCommit
SelectedFile
SelectedPath
SelectedLocalBranch
SelectedRemoteBranch
SelectedRemote
SelectedTag
SelectedStashEntry
SelectedCommitFile
CheckedOutBranch
```

To see what fields are available on e.g. the `SelectedFile`, see [here](https://github.com/jesseduffield/lazygit/blob/master/pkg/commands/models/file.go) (all the modelling lives in the same directory). Note that the custom commands feature does not guarantee backwards compatibility (until we hit lazygit version 1.0 of course) which means a field you're accessing on an object may no longer be available from one release to the next. Typically however, all you'll need is `{{.SelectedFile.Name}}`, `{{.SelectedLocalCommit.Sha}}` and `{{.SelectedLocalBranch.Name}}`. In the future we will likely introduce a tighter interface that exposes a limited set of fields for each model.

### Keybinding collisions

If your custom keybinding collides with an inbuilt keybinding that is defined for the same context, only the custom keybinding will be executed. This also applies to the global context. However, one caveat is that if you have a custom keybinding defined on the global context for some key, and there is an in-built keybinding defined for the same key and for a specific context (say the 'files' context), then the in-built keybinding will take precedence. See how to change in-built keybindings [here](https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md#keybindings)

### Debugging

If you want to verify that your command actually does what you expect, you can wrap it in an 'echo' call and set `showOutput: true` so that it doesn't actually execute the command but you can see how the placeholders were resolved. Alternatively you can run lazygit in debug mode with `lazygit --debug` and in another terminal window run `lazygit --logs` to see which commands are actually run

### More Examples

See the [wiki](https://github.com/jesseduffield/lazygit/wiki/Custom-Commands-Compendium) page for more examples, and feel free to add your own custom commands to this page so others can benefit!
