package config

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeybindingPlatform(t *testing.T) {
	scenarios := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Not set falls back to the host OS",
			envValue: "",
			expected: runtime.GOOS,
		},
		{
			name:     "darwin is honored",
			envValue: "darwin",
			expected: "darwin",
		},
		{
			name:     "linux is honored",
			envValue: "linux",
			expected: "linux",
		},
		{
			name:     "windows is honored",
			envValue: "windows",
			expected: "windows",
		},
		{
			name:     "An unrecognized value falls back to the host OS",
			envValue: "mac",
			expected: runtime.GOOS,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			t.Setenv("LAZYGIT_KEYBINDING_PLATFORM", s.envValue)
			assert.Equal(t, s.expected, KeybindingPlatform())
		})
	}
}

func TestMigrationOfRenamedKeys(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No rename needed",
			input: `foo:
  bar: 5
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Rename one",
			input: `gui:
  skipUnstageLineWarning: true
`,
			expected: `gui:
  skipDiscardChangeWarning: true
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Renamed 'gui.skipUnstageLineWarning' to 'skipDiscardChangeWarning'"},
		},
		{
			name: "Rename several",
			input: "gui:\n" +
				"  windowSize: half\n" +
				"  skipUnstageLineWarning: true\n" +
				"keybinding:\n" +
				"  universal:\n" +
				"    executeCustomCommand: a\n" +
				"    cyclePagers: b\n" +
				"    cyclePagersReverse: c\n",
			expected: "gui:\n" +
				"  screenMode: half\n" +
				"  skipDiscardChangeWarning: true\n" +
				"keybinding:\n" +
				"  universal:\n" +
				"    executeShellCommand: a\n" +
				"    cycleDiffRenderers: b\n" +
				"    cycleDiffRenderersReverse: c\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed 'gui.skipUnstageLineWarning' to 'skipDiscardChangeWarning'",
				"Renamed 'keybinding.universal.executeCustomCommand' to 'executeShellCommand'",
				"Renamed 'keybinding.universal.cyclePagers' to 'cycleDiffRenderers'",
				"Renamed 'keybinding.universal.cyclePagersReverse' to 'cycleDiffRenderersReverse'",
				"Renamed 'gui.windowSize' to 'screenMode'",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestMigrationOfMovedKeys(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No move needed",
			input: `foo:
  bar: 5
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Move worktree keybinding into the universal section",
			input: `keybinding:
  universal:
    quit: q
  worktrees:
    viewWorktreeOptions: w
`,
			expected: `keybinding:
  universal:
    quit: q
    newWorktree: w
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Moved 'keybinding.worktrees.viewWorktreeOptions' to 'keybinding.universal.newWorktree'"},
		},
		{
			name: "Create the universal section if it doesn't exist",
			input: `keybinding:
  worktrees:
    viewWorktreeOptions: w
`,
			expected: `keybinding:
  universal:
    newWorktree: w
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Moved 'keybinding.worktrees.viewWorktreeOptions' to 'keybinding.universal.newWorktree'"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestMigrateNullKeybindingsToDisabled(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No change needed",
			input: `keybinding:
  universal:
    quit: q
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Change one",
			input: `keybinding:
  universal:
    quit: null
`,
			expected: `keybinding:
  universal:
    quit: <disabled>
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.quit'"},
		},
		{
			name: "Change several",
			input: `keybinding:
  universal:
    quit: null
    return: <esc>
    new: null
`,
			expected: `keybinding:
  universal:
    quit: <disabled>
    return: <esc>
    new: <disabled>
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.quit'",
				"Changed 'null' to '<disabled>' for keybinding 'keybinding.universal.new'",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestCommitPrefixMigrations(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Single CommitPrefix Rename",
			input: `git:
  commitPrefix:
     pattern: "^\\w+-\\w+.*"
     replace: '[JIRA $0] '
`,
			expected: `git:
  commitPrefix:
    - pattern: "^\\w+-\\w+.*"
      replace: '[JIRA $0] '
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'git.commitPrefix' to an array of strings"},
		},
		{
			name: "Complicated CommitPrefixes Rename",
			input: `git:
  commitPrefixes:
    foo:
      pattern: "^\\w+-\\w+.*"
      replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      pattern: "^foo.bar*"
      replace: '[FUN $0] '
`,
			expected: `git:
  commitPrefixes:
    foo:
      - pattern: "^\\w+-\\w+.*"
        replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      - pattern: "^foo.bar*"
        replace: '[FUN $0] '
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'git.commitPrefixes' elements to arrays of strings"},
		},
		{
			name:              "Incomplete Configuration",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No changes made when already migrated",
			input: `
git:
   commitPrefix:
    - pattern: "Hello World"
      replace: "Goodbye"
   commitPrefixes:
    foo:
      - pattern: "^\\w+-\\w+.*"
        replace: '[JIRA $0] '`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestCustomCommandsOutputMigration(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Empty String",
			input:             "",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Convert subprocess to output=terminal",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: terminal
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'subprocess: true' to 'output: terminal' in custom command"},
		},
		{
			name: "Convert stream to output=log",
			input: `customCommands:
  - command: echo 'hello'
    stream: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: log
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'stream: true' to 'output: log' in custom command"},
		},
		{
			name: "Convert showOutput to output=popup",
			input: `customCommands:
  - command: echo 'hello'
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: popup
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Changed 'showOutput: true' to 'output: popup' in custom command"},
		},
		{
			name: "Subprocess wins over the other two",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: true
    stream: true
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: terminal
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'subprocess: true' to 'output: terminal' in custom command",
				"Deleted redundant 'stream: true' property in custom command",
				"Deleted redundant 'showOutput: true' property in custom command",
			},
		},
		{
			name: "Stream wins over showOutput",
			input: `customCommands:
  - command: echo 'hello'
    stream: true
    showOutput: true
  `,
			expected: `customCommands:
  - command: echo 'hello'
    output: log
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Changed 'stream: true' to 'output: log' in custom command",
				"Deleted redundant 'showOutput: true' property in custom command",
			},
		},
		{
			name: "Explicitly setting to false doesn't create an output=none key",
			input: `customCommands:
  - command: echo 'hello'
    subprocess: false
    stream: false
    showOutput: false
  `,
			expected: `customCommands:
  - command: echo 'hello'
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Deleted redundant 'subprocess: false' in custom command",
				"Deleted redundant 'stream: false' property in custom command",
				"Deleted redundant 'showOutput: false' property in custom command",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestAllBranchesLogCmdMigrations(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		{
			name:              "Incomplete Configuration Passes uneventfully",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Single Cmd with no Cmds",
			input: `git:
  allBranchesLogCmd: git log --graph --oneline
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Created git.allBranchesLogCmds array containing value of git.allBranchesLogCmd",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Cmd with one existing Cmds",
			input: `git:
  allBranchesLogCmd: git log --graph --oneline
  allBranchesLogCmds:
    - git log --graph --oneline --pretty
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
    - git log --graph --oneline --pretty
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Prepended git.allBranchesLogCmd value to git.allBranchesLogCmds array",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Only Cmds set have no changes",
			input: `git:
  allBranchesLogCmds:
    - git log
`,
			expected:        "",
			expectedChanges: []string{},
		},
		{
			name: "Removes Empty Cmd When at end of yaml",
			input: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  allBranchesLogCmd:
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Removed obsolete git.allBranchesLogCmd"},
		},
		{
			name: "Migrates when sequence defined inline",
			input: `git:
  allBranchesLogCmds: [foo, bar]
  allBranchesLogCmd: baz
`,
			expected: `git:
  allBranchesLogCmds: [baz, foo, bar]
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Prepended git.allBranchesLogCmd value to git.allBranchesLogCmds array",
				"Removed obsolete git.allBranchesLogCmd",
			},
		},
		{
			name: "Removes Empty Cmd With Keys Afterwards",
			input: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  allBranchesLogCmd:
  foo: bar
`,
			expected: `git:
  allBranchesLogCmds:
    - git log --graph --oneline
  foo: bar
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Removed obsolete git.allBranchesLogCmd"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}

func TestPagerMigration(t *testing.T) {
	scenarios := []struct {
		name              string
		input             string
		expected          string
		expectedDidChange bool
		expectedChanges   []string
	}{
		// Migrate 'paging' to 'pagers' array
		{
			name:              "Incomplete Configuration Passes uneventfully",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No paging section",
			input: "git:\n" +
				"  autoFetch: true\n",
			expected: "git:\n" +
				"  autoFetch: true\n",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "paging is not an object",
			input: "git:\n" +
				"  paging: 5\n",
			expected: "git:\n" +
				"  paging: 5\n",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "pagers is not an array",
			input: "git:\n" +
				"  pagers: 5\n",
			expected: "git:\n" +
				"  pagers: 5\n",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "paging and pagers coexist",
			input: "git:\n" +
				"  paging:\n" +
				"    pager: delta --dark --paging=never\n" +
				"  pagers:\n" +
				"    - pager: diff-so-fancy\n",
			expected: "git:\n" +
				"  paging:\n" +
				"    pager: delta --dark --paging=never\n" +
				"  diffRenderers:\n" +
				"    - command: diff-so-fancy\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed git.pagers to git.diffRenderers",
				"Renamed 'pager' to 'command' in git pager",
			},
		},
		{
			name: "paging and diffRenderers coexist",
			input: "git:\n" +
				"  paging:\n" +
				"    pager: delta --dark --paging=never\n" +
				"  diffRenderers:\n" +
				"    - command: diff-so-fancy\n",
			expected: "git:\n" +
				"  paging:\n" +
				"    pager: delta --dark --paging=never\n" +
				"  diffRenderers:\n" +
				"    - command: diff-so-fancy\n",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "pagers and diffRenderers coexist",
			input: "git:\n" +
				"  pagers:\n" +
				"    - pager: delta --dark --paging=never\n" +
				"  diffRenderers:\n" +
				"    - command: diff-so-fancy\n",
			expected: "git:\n" +
				"  pagers:\n" +
				"    - pager: delta --dark --paging=never\n" +
				"  diffRenderers:\n" +
				"    - command: diff-so-fancy\n",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "paging is moved to diffRenderers array preserving fields and order",
			input: "git:\n" +
				"  paging:\n" +
				"    name: delta\n" +
				"    colorArg: never\n" +
				"    pager: delta --dark --paging=never\n" +
				"  autoFetch: true\n",
			expected: "git:\n" +
				"  diffRenderers:\n" +
				"    - name: delta\n" +
				"      colorArg: never\n" +
				"      command: delta --dark --paging=never\n" +
				"  autoFetch: true\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Moved git.paging object to git.pagers array",
				"Renamed git.pagers to git.diffRenderers",
				"Renamed 'pager' to 'command' in git pager",
			},
		},
		{
			name: "paging is moved to diffRenderers array even if empty",
			input: "git:\n" +
				"  paging: {}\n",
			expected: "git:\n" +
				"  diffRenderers:\n" +
				"    - type: rawGit\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Moved git.paging object to git.pagers array",
				"Renamed git.pagers to git.diffRenderers",
				"Changed git pager without a command to 'type: rawGit'",
			},
		},

		// Migrate 'pagers' array to 'diffRenderers' array
		{
			name: "empty pagers array is renamed",
			input: "git:\n" +
				"  pagers: []\n",
			expected: "git:\n" +
				"  diffRenderers: []\n",
			expectedDidChange: true,
			expectedChanges:   []string{"Renamed git.pagers to git.diffRenderers"},
		},
		{
			name: "pagers array entries are adapted",
			input: "git:\n" +
				"  pagers:\n" +
				"    - name: delta\n" +
				"      colorArg: never\n" +
				"      pager: delta --dark --paging=never\n" +
				"    - name: difft\n" +
				"      colorArg: never\n" +
				"      externalDiffCommand: difft --color=always\n" +
				"    - name: git-config\n" +
				"      colorArg: never\n" +
				"      useExternalDiffGitConfig: TRUE\n" +
				"    - name: git\n" +
				"      colorArg: never\n" +
				"  autoFetch: true\n",
			expected: "git:\n" +
				"  diffRenderers:\n" +
				"    - name: delta\n" +
				"      colorArg: never\n" +
				"      command: delta --dark --paging=never\n" +
				"    - name: difft\n" +
				"      colorArg: never\n" +
				"      command: difft --color=always\n" +
				"      type: extDiff\n" +
				"    - name: git-config\n" +
				"      colorArg: never\n" +
				"      type: extDiff\n" +
				"    - name: git\n" +
				"      colorArg: never\n" +
				"      type: rawGit\n" +
				"  autoFetch: true\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed git.pagers to git.diffRenderers",
				"Renamed 'pager' to 'command' in git pager",
				"Changed 'externalDiffCommand' to 'command' with 'type: extDiff' in git pager",
				"Changed 'useExternalDiffGitConfig: true' to 'type: extDiff' in git pager",
				"Changed git pager without a command to 'type: rawGit'",
			},
		},
		{
			name: "zero-valued mechanism fields do not take precedence and are removed",
			input: "git:\n" +
				"  pagers:\n" +
				"    - pager: delta --dark --paging=never\n" +
				"      externalDiffCommand: null\n" +
				"      useExternalDiffGitConfig: false\n" +
				"    - pager: \"\"\n" +
				"      externalDiffCommand: difft --color=always\n" +
				"      useExternalDiffGitConfig: false\n" +
				"    - pager: null\n" +
				"      externalDiffCommand: \"\"\n" +
				"      useExternalDiffGitConfig: YES\n" +
				"    - pager: \"\"\n" +
				"      externalDiffCommand:\n" +
				"      useExternalDiffGitConfig: false\n" +
				"    - useExternalDiffGitConfig: null\n",
			expected: "git:\n" +
				"  diffRenderers:\n" +
				"    - command: delta --dark --paging=never\n" +
				"    - command: difft --color=always\n" +
				"      type: extDiff\n" +
				"    - type: extDiff\n" +
				"    - type: rawGit\n" +
				"    - type: rawGit\n",
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed git.pagers to git.diffRenderers",
				"Renamed 'pager' to 'command' in git pager",
				"Removed empty 'externalDiffCommand' from git pager",
				"Removed 'useExternalDiffGitConfig: false' from git pager",
				"Changed 'externalDiffCommand' to 'command' with 'type: extDiff' in git pager",
				"Removed empty 'pager' from git pager",
				"Changed 'useExternalDiffGitConfig: true' to 'type: extDiff' in git pager",
				"Changed git pager without a command to 'type: rawGit'",
				"Removed empty 'useExternalDiffGitConfig' from git pager",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			changes := NewChangesSet()
			actual, didChange, err := computeMigratedConfig("path doesn't matter", []byte(s.input), changes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedDidChange, didChange)
			if didChange {
				assert.Equal(t, s.expected, string(actual))
			}
			assert.Equal(t, s.expectedChanges, changes.ToSliceFromOldest())
		})
	}
}
