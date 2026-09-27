package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UnfocusedListShowsSelectionWhenFilled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A list that gets its first item while the focus is elsewhere starts showing a selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			IsEmpty().
			SelectionIsHidden().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().IsFocused()

		t.Views().Files().
			Tap(func() {
				t.Shell().CreateFile("file2", "two\n")
				t.RefreshInBackground()
			}).
			Lines(Contains("file2")).
			SelectionIsInactive()
	},
})
