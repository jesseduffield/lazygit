package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoKeepsFocusedMainViewSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undoing a commit rewrite while focused in the main view re-establishes the (stale, multi-line) selection on a surviving change rather than leaving it painted over the changed diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("commit to rewrite")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to rewrite").IsSelected(),
				Contains("first commit"),
			).
			// Focus the commit's diff straight from the commits panel.
			Press(keys.Universal.FocusMainView)

		// Discard the first line from the commit; the selection advances to the next
		// surviving change ('two').
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+one"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard lines from commit")).
			Content(Equals("Are you sure you want to discard the selected lines from this commit?")).
			Confirm()

		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("+one")).
			SelectedLines(
				Contains("+two"),
			).
			// Leave a multi-line range selected before undoing — undo rewrites the commit
			// outside the focused-main-view action handlers, so without the preserve net
			// this range would be left stale over the restored diff.
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+four")).
			SelectedLines(
				Contains("+two"),
				Contains("+three"),
				Contains("+four"),
			)

		// Undo is a global keybinding, so it fires while the main view holds focus.
		t.GlobalPress(keys.Universal.Undo)

		t.ExpectPopup().Confirmation().
			Title(Equals("Undo")).
			Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\?`)).
			Confirm()

		// The discarded line is back, and the stale multi-line range collapses to a
		// single surviving change at the selection's top ordinal (the first change line,
		// now 'one' again) rather than spanning arbitrary lines of the restored diff.
		t.Views().Main().
			IsFocused().
			ContainsLines(
				Equals("+one"),
				Equals("+two"),
				Equals("+three"),
				Equals("+four"),
				Equals("+five"),
			).
			SelectedLines(
				Contains("+one"),
			)
	},
})
