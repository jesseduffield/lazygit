# lazygit Configuration Guide

## Table of Contents
- [Basic Configuration](#basic-configuration)
- [Configuration File Locations](#configuration-file-locations)
- [Configuration Examples](#configuration-examples)
- [Theme and Appearance](#theme-and-appearance)
- [Custom Commands](#custom-commands)
- [Git Integration](#git-integration)
- [Keybindings](#keybindings)
- [Editor Configuration](#editor-configuration)
- [Performance and Behavior](#performance-and-behavior)
- [Troubleshooting](#troubleshooting)

## Basic Configuration

Config file locations:
```bash
# Linux/MacOS
~/.config/lazygit/config.yml

# Windows
%LOCALAPPDATA%\lazygit\config.yml
```

Minimal config:
```yaml
gui:
  theme:
    activeBorderColor:
      - green
      - bold
  sidePanelWidth: 0.3
git:
  paging:
    colorArg: always
```

## Configuration File Locations

### Global Configuration

| OS | Primary Location | Legacy Location |
|---|---|---|
| Linux | `~/.config/lazygit/config.yml` | `~/.config/jesseduffield/lazygit/config.yml` |
| macOS | `~/Library/Application Support/lazygit/config.yml` | `~/Library/Application Support/jesseduffield/lazygit/config.yml` |
| Windows | `%LOCALAPPDATA%\lazygit\config.yml` | `%APPDATA%\jesseduffield\lazygit\config.yml` |

### Repository-Specific Configuration

Override global settings with repository-specific configs:

1. **In-repo config:** `<repo>/.git/lazygit.yml`
2. **Parent directory configs:** `.lazygit.yml` in any parent directory

Example hierarchy:
```
~/projects/
  .lazygit.yml          # Applies to all projects
  work/
    .lazygit.yml        # Applies to all work projects
    client-a/
      .git/lazygit.yml  # Specific to this repo
```

## Configuration Examples

### Dark Theme

```yaml
gui:
  theme:
    activeBorderColor:
      - cyan
      - bold
    inactiveBorderColor:
      - white
    selectedLineBgColor:
      - blue
    defaultFgColor:
      - white
  showNumstatInFilesView: true
  sidePanelWidth: 0.2
```

### Skip Confirmations

```yaml
gui:
  skipStashWarning: true
  skipDiscardChangeWarning: false
  skipRewordInEditorWarning: true
  showRandomTip: true
  showFileTree: true
  showRootItemInFileTree: true

git:
  autoFetch: true
  autoRefresh: true
  autoStageResolvedConflicts: true
```

### Git Log Format

```yaml
git:
  branchLogCmd: "git log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative {{branchName}} --"
  
  allBranchesLogCmds:
    - "git log --graph --all --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"
    - "git log --graph --all --oneline --decorate"
    - "git log --graph --all --format='%h %s' --simplify-by-decoration"
```

### Nerd Font Icons

```yaml
gui:
  nerdFontsVersion: "3"  # or "2"
  showFileIcons: true
  
  customIcons:
    filenames:
      "README.md": 
        icon: "\uf48a"
        color: "#0366d6"
      "package.json":
        icon: "\ue718"
        color: "#8bc34a"
    extensions:
      ".go":
        icon: "\ue626"
        color: "#00add8"
      ".js":
        icon: "\ue74e"
        color: "#f7df1e"
```

## Theme and Appearance

### Author Colors

```yaml
gui:
  authorColors:
    'John Doe': 'cyan'
    'Jane Smith': '#ff00ff'
    '*': 'white'
```

### Branch Colors

```yaml
gui:
  branchColorPatterns:
    '^feature/': '#00ff00'
    '^bugfix/': '#ff0000'
    '^release/': '#ffff00'
    'main|master': '#00ffff'
```

### Selected Line

```yaml
gui:
  theme:
    # Bold text only
    selectedLineBgColor:
      - default
      - bold
    
    # Reverse colors
    selectedLineBgColor:
      - reverse
```

## Custom Commands

### Conventional Commits

```yaml
customCommands:
  - key: "C"
    context: "files"
    description: "Create conventional commit"
    prompts:
      - type: "menu"
        title: "Commit Type"
        key: "Type"
        options:
          - value: "feat"
          - value: "fix"
          - value: "docs"
          - value: "style"
          - value: "refactor"
      - type: "input"
        title: "Scope (optional)"
        key: "Scope"
      - type: "input"
        title: "Summary"
        key: "Summary"
    command: "git commit -m '{{.Form.Type}}{{if .Form.Scope}}({{.Form.Scope}}){{end}}: {{.Form.Summary}}'"
```

### Commit and Push

```yaml
customCommands:
  - key: "P"
    context: "global"
    description: "Commit and push"
    command: "git commit && git push"
    subprocess: true
```

### External Diff Tool

```yaml
customCommands:
  - key: "D"
    context: "files"
    description: "Diff with external tool"
    command: "git difftool {{.SelectedFile.Name}}"
    subprocess: true
```

## Git Integration

### Pager Configuration

```yaml
git:
  paging:
    # diff-so-fancy
    pager: "diff-so-fancy"
    colorArg: always
    
    # delta
    pager: "delta --dark --paging=never --line-numbers"
    colorArg: always
    
    # ydiff
    pager: "ydiff -p cat -s --wrap --width={{columnWidth}}"
    colorArg: always
```

### Commit Prefixes

```yaml
git:
  commitPrefix:
    - pattern: "^(\\w+)/(\\w+-\\d+).*"
      replace: "[$2] "
  
  commitPrefixes:
    my-project:
      - pattern: "^feature/(\\w+-\\d+).*"
        replace: "feat($1): "
      - pattern: "^bugfix/(\\w+-\\d+).*"  
        replace: "fix($1): "
```

### Branch Configuration

```yaml
git:
  mainBranches:
    - master
    - main
    - develop
  
  autoForwardBranches: allBranches  # or 'onlyMainBranches', 'none'
  
  merging:
    args: "--no-ff"
    squashMergeMessage: "Squash merge {{selectedRef}} into {{currentBranch}}"
```

## Keybindings

### Alternative Keybindings

```yaml
keybinding:
  universal:
    prevItem-alt: 'k'
    nextItem-alt: 'j'
    prevBlock-alt: 'h'
    nextBlock-alt: 'l'
    scrollUpMain-alt1: '<c-u>'
    scrollDownMain-alt1: '<c-d>'
    gotoTop: 'gg'
    gotoBottom: 'G'
```

### Disable Bindings

```yaml
keybinding:
  universal:
    edit: '<disabled>'
  files:
    ignoreFile: '<disabled>'
```

### Custom Submit/Cancel

```yaml
keybinding:
  universal:
    confirm: '<c-y>'
    return: '<c-n>'
```

## Editor Configuration

### VSCode

```yaml
os:
  editPreset: 'vscode'
```

### Custom Editor

```yaml
os:
  edit: 'emacsclient -n {{filename}}'
  editAtLine: 'emacsclient -n +{{line}} {{filename}}'
  editInTerminal: false
```

### Neovim Remote

```yaml
os:
  editPreset: 'nvim-remote'
```

## Performance and Behavior

### Large Repositories

```yaml
gui:
  refreshInterval: 30
  commitHashLength: 6
  showFileTree: false
  showNumstatInFilesView: false

git:
  autoFetch: false
  autoRefresh: false
```

### Mouse Support

```yaml
gui:
  mouseEvents: false
```

### Startup

```yaml
disableStartupPopups: true
notARepository: 'prompt'  # or 'create', 'skip', 'quit'
promptToReturnFromSubprocess: false
```

## Troubleshooting

Debug mode:
```bash
lazygit --debug
lazygit --logs
```

### Colors in tmux

```yaml
gui:
  theme:
    defaultFgColor:
      - 'color250'
```

### External Commands

```yaml
os:
  shellFunctionsFile: ~/.my_aliases.sh
```

### Config Not Loading

Check location:
```bash
lazygit --print-config-dir
```

Validate YAML:
```bash
python -c "import yaml; yaml.safe_load(open('config.yml'))"
```

### Dubious Ownership

```bash
git config --global --add safe.directory /path/to/repo
```

### Config Merging

```bash
lazygit --use-config-file="~/.base_config.yml,~/.theme_config.yml"
LG_CONFIG_FILE="~/.base_config.yml,~/.theme_config.yml" lazygit
```

### IntelliSense

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/jesseduffield/lazygit/master/schema/config.json
```