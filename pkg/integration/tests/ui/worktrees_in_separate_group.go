package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var WorktreesInSeparateGroup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify worktrees appears as its own side panel group when worktreesInSeparateGroup is enabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.WorktreesInSeparateGroup = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Default focus is files. Go prev to reach worktrees (its own group now).
		t.Views().Files().IsFocused().
			Press(keys.Universal.PrevBlock)

		t.Views().Worktrees().IsFocused().
			Press(keys.Universal.PrevBlock)

		// Prev from worktrees should go to status
		t.Views().Status().IsFocused()

		// Navigate forward: status -> worktrees -> files -> branches
		t.Views().Status().Press(keys.Universal.NextBlock)
		t.Views().Worktrees().IsFocused().
			Press(keys.Universal.NextBlock)
		t.Views().Files().IsFocused().
			Press(keys.Universal.NextBlock)
		t.Views().Branches().IsFocused()

		// Test jump keys: key 2 should jump to worktrees
		t.GlobalPress(keys.Universal.JumpToBlock[1])
		t.Views().Worktrees().IsFocused()

		// Key 3 should jump to files
		t.GlobalPress(keys.Universal.JumpToBlock[2])
		t.Views().Files().IsFocused()
	},
})
