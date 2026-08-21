package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickSelectsDiffLine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Clicking a line of the main view's diff focuses the view and selects that line",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused()

		// The click both focuses the view and points at a line, so that line is selected
		// rather than the first change.
		t.Views().Main().
			Click(0, 4).
			IsFocused().
			SelectionIsActive().
			SelectedLines(
				Contains("@@ -1,5 +1,5 @@"),
			).
			// A click in the already-focused view moves the selection to the clicked line.
			Click(0, 7).
			SelectedLines(
				Contains("-three"),
			)
	},
})
