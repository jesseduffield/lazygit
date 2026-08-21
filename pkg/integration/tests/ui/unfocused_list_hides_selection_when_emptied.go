package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UnfocusedListHidesSelectionWhenEmptied = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A list that loses its last item while the focus is elsewhere stops showing a selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")
		shell.UpdateFile("file1", "two\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(Contains("file1")).
			SelectionIsActive().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().IsFocused()

		t.Views().Files().
			SelectionIsInactive().
			Tap(func() {
				t.Shell().RunCommand([]string{"git", "checkout", "--", "file1"})
				t.RefreshInBackground()
			}).
			IsEmpty().
			SelectionIsHidden()
	},
})
