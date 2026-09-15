package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleWhitespaceKeepsUnfocusedSelectionDimmed = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggling whitespace from the main view leaves the panel beneath it showing a dimmed selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")
		shell.UpdateFile("file1", "  one\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.ToggleWhitespaceInDiffView)

		t.Views().Files().
			SelectionIsInactive()
	},
})
