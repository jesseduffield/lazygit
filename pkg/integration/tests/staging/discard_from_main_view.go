package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard a hunk from the working tree in the focused main view, without diving into the staging view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// Two separate change blocks, far enough apart to stay distinct hunks.
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			// `a` extends the selection to the whole change block around the cursor.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard change")).
			Content(Contains("Are you sure you want to discard this change")).
			Confirm()

		// The first block was discarded from the working tree, leaving line 3 at its
		// committed value; the other block stays.
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("THREE")).
			ContainsLines(
				Contains("-nine"),
				Contains("+NINE"),
			)
	},
})
