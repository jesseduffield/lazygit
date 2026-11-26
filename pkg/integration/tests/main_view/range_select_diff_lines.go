package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RangeSelectDiffLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a range of diff lines in the focused main view, both sticky and with shift",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nFOUR\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// A sticky range is extended by the plain arrow keys, and pressing the key again
		// collapses it back to the cursor line.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-three"),
				Contains("-four"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			SelectedLines(
				Contains("-four"),
			).
			// A non-sticky range only grows while shift is held, and a plain move
			// collapses it again.
			Press(keys.Universal.RangeSelectDown).
			SelectedLines(
				Contains("-four"),
				Contains("+THREE"),
			).
			Press(keys.Universal.RangeSelectUp).
			SelectedLines(
				Contains("-four"),
			).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+THREE"),
			)
	},
})
