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
			input: `gui:
  windowSize: half
  skipUnstageLineWarning: true
keybinding:
  universal:
    executeCustomCommand: a
`,
			expected: `gui:
  screenMode: half
  skipDiscardChangeWarning: true
keybinding:
  universal:
    executeShellCommand: a
`,
			expectedDidChange: true,
			expectedChanges: []string{
				"Renamed 'gui.skipUnstageLineWarning' to 'skipDiscardChangeWarning'",
				"Renamed 'keybinding.universal.executeCustomCommand' to 'executeShellCommand'",
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
		{
			name:              "Incomplete Configuration Passes uneventfully",
			input:             "git:",
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "No paging section",
			input: `git:
  autoFetch: true
`,
			expected: `git:
  autoFetch: true
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "Both paging and pagers exist",
			input: `git:
  paging:
    pager: delta --dark --paging=never
  pagers:
    - diff: diff-so-fancy
`,
			expected: `git:
  paging:
    pager: delta --dark --paging=never
  pagers:
    - diff: diff-so-fancy
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "paging is not an object",
			input: `git:
  paging: 5
`,
			expected: `git:
  paging: 5
`,
			expectedDidChange: false,
			expectedChanges:   []string{},
		},
		{
			name: "paging is moved to pagers array (keeping the order)",
			input: `git:
  paging:
    pager: delta --dark --paging=never
  autoFetch: true
`,
			expected: `git:
  pagers:
    - pager: delta --dark --paging=never
  autoFetch: true
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Moved git.paging object to git.pagers array"},
		},
		{
			name: "paging is moved to pagers array even if empty",
			input: `git:
  paging: {}
`,
			expected: `git:
  pagers: [{}]
`,
			expectedDidChange: true,
			expectedChanges:   []string{"Moved git.paging object to git.pagers array"},
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
