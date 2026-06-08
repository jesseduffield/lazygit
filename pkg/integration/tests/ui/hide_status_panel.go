package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideStatusPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify status panel is hidden when hideStatusPanel is enabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.HideStatusPanel = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Default focus is files
		t.Views().Files().IsFocused()

		// Next block should go to branches (skipping status)
		t.Views().Files().Press(keys.Universal.NextBlock)
		t.Views().Branches().IsFocused()

		// Continue forward: commits, stash
		t.Views().Branches().Press(keys.Universal.NextBlock)
		t.Views().Commits().IsFocused()

		t.Views().Commits().Press(keys.Universal.NextBlock)
		t.Views().Stash().IsFocused()

		// Wrapping forward from stash should go to files (skipping status)
		t.Views().Stash().Press(keys.Universal.NextBlock)
		t.Views().Files().IsFocused()

		// Cycling backwards from files should wrap to stash (skipping status)
		t.Views().Files().Press(keys.Universal.PrevBlock)
		t.Views().Stash().IsFocused()
	},
})
