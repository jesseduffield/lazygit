package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardLinesFromACommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard the selected lines of a commit's diff from the commit itself, in the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\nFOUR\nfive\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+TWO")).
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard lines from commit")).
			Content(Contains("Are you sure you want to discard the selected lines from this commit?")).
			Confirm()

		// The commit keeps its other change and has given up the one discarded, and the
		// selection carries on from the change that has taken its place.
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("+TWO")).
			SelectedLines(
				Contains("-four"),
			)

		// The rewrite is the commit's own business: nothing is left lying in the working
		// tree.
		t.Views().Files().IsEmpty()
	},
})
