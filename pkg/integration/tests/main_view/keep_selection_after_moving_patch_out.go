package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectionAfterMovingPatchOut = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a custom patch out of a commit leaves the focused main view's selection on a change that is still there",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "ONE\ntwo\nTHREE\nfour\nFIVE\n")
		shell.Commit("commit to move a patch out of")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move a patch out of").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Take the first modification into a custom patch.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-one"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-one"),
				Contains("+ONE"),
			).
			PressPrimaryAction()

		// Keep a range selected across lines that the pending patch will remove from
		// the commit and a later change that will remain.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("-one")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+THREE")).
			SelectedLines(
				Contains("-one"),
				Contains("+ONE"),
				Contains(" two"),
				Contains("-three"),
				Contains("+THREE"),
			)

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		// The moved lines are gone from the commit, so the range collapses onto the
		// change that has taken their place.
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("+ONE")).
			SelectedLines(
				Contains("-three"),
			)
	},
})
