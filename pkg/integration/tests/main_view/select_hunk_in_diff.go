package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectHunkInDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle hunk selection in the focused main view, and step from hunk to hunk",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nFOUR\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Hunk mode widens the selection to the whole change block around the cursor —
		// which is lazygit's notion of a hunk, so the two changed lines and their
		// replacements are one block, and the isolated change further down is another.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("-four"),
				Contains("+THREE"),
				Contains("+FOUR"),
			).
			// In hunk mode the arrow keys step from block to block rather than by line.
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			).
			Press(keys.Universal.PrevItem).
			SelectedLines(
				Contains("-three"),
				Contains("-four"),
				Contains("+THREE"),
				Contains("+FOUR"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
			)
	},
})
