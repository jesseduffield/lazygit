package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchCollapsesTheSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Searching the focused main view leaves a single line selected at the match, whatever was selected before",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// The match is in another block entirely, so the hunk selection the cursor
			// has just left goes with it.
			FilterOrSearch("NINE").
			SelectedLines(
				Contains("+NINE"),
			).
			// And the same for a hunk selected while a search is on: the next match is
			// not part of it either.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			).
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("+NINE"),
			)
	},
})
