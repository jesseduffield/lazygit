package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectDiffLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view selects the first visible change line, and the arrow keys move the selection",
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
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// The selection starts on the first change line rather than at the top of the
		// diff, so that it lands on something to act on without the view jumping.
		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+THREE"),
			).
			Press(keys.Universal.PrevItem).
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Universal.GotoTop).
			SelectedLines(
				Contains("diff --git a/file1 b/file1"),
			).
			Press(keys.Universal.Return)

		t.Views().Files().
			IsFocused()
	},
})
