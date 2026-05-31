package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusMainViewFromSecondary = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focus the primary main view with its keybinding while the secondary main view is focused",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Two separate hunks so we can stage one and leave the other unstaged,
		// splitting the main panel into a primary (unstaged) and secondary (staged) view
		shell.
			CreateFileAndAdd("file1", "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\n").
			Commit("initial commit").
			UpdateFile("file1", "line1-changed\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10-changed\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Stage the first hunk so file1 is partially staged and the main panel splits
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

		// Focus the primary main view from the side panel
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Toggle across to the secondary main view
		t.Views().Main().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		// The focus-main-view key returns focus to the primary view from the secondary
		t.Views().Secondary().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused()
	},
})
