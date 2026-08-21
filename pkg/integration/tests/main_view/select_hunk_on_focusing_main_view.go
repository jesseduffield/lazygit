package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectHunkOnFocusingMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When hunk mode is the default, focusing the main view selects the first whole change block",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = true
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

		// No key press needed: the whole block is selected just by focusing.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Hunk mode being the configured default, it isn't something escape gives up:
			// escape leaves the view.
			Press(keys.Universal.Return)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			// A click on a change line keeps hunk mode and selects that line's block.
			Click(0, 14).
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			).
			// A click on a context line points at it precisely, so it stays a single line.
			Click(0, 12).
			SelectedLines(
				Contains(" seven"),
			)
	},
})
