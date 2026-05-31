package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusSecondaryView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focus the secondary main view with the keybinding and switch between main views",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Two separate hunks so we can stage one and leave the other unstaged
		shell.
			CreateFileAndAdd("file1", "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\n").
			Commit("initial commit").
			UpdateFile("file1", "line1-changed\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10-changed\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// No secondary view yet — pressing the focus-secondary key is a no-op
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusSecondaryView).
			IsFocused()

		// Enter staging view and stage first hunk to create a partially staged file
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			PressPrimaryAction(). // stage first hunk
			PressEscape()

		// Back in Files panel with file1 partially staged — split main panel is visible
		t.Views().Files().
			IsFocused()

		// Focus secondary view from side panel
		t.Views().Files().
			Press(keys.Universal.FocusSecondaryView)

		t.Views().Secondary().
			IsFocused()

		// Switch from secondary to primary main view
		t.Views().Secondary().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused()

		// Switch from primary back to secondary main view
		t.Views().Main().
			Press(keys.Universal.FocusSecondaryView)

		t.Views().Secondary().
			IsFocused()

		// Escape back to side panel
		t.Views().Secondary().
			PressEscape()

		t.Views().Files().
			IsFocused()
	},
})
