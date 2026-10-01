package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EscapeDismissesSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Escape gives up a range selection, then hunk mode, before leaving the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nTWO\nTHREE\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// A sticky range: escape collapses it to the cursor line rather than leaving.
		t.Views().Main().
			IsFocused().
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-two"),
				Contains("-three"),
			).
			Press(keys.Universal.Return).
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			// Hunk mode the user asked for: escape goes back to line-by-line.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-two"),
				Contains("-three"),
				Contains("+TWO"),
				Contains("+THREE"),
			).
			Press(keys.Universal.Return).
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			// With nothing left to give up, escape leaves.
			Press(keys.Universal.Return)

		t.Views().Files().
			IsFocused()
	},
})
