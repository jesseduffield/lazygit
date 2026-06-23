package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectionAfterMovingPatchOutMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a custom patch out of a commit from the focused main view re-establishes the (stale, multi-line) selection on a surviving change rather than leaving it painted over the shrunk diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("commit to move a patch out of")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to move a patch out of").IsSelected(),
				Contains("first commit"),
			).
			// Focus the commit's diff straight from the commits panel, rather than
			// entering the commit files panel first.
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+one"),
			).
			// Toggle just the first line into a custom patch, then leave a multi-line
			// range selected — the patch move below doesn't go through the focused-main-
			// view action handlers, so without the preserve net this stale range would be
			// left painted over the shrunk diff.
			PressPrimaryAction().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+four")).
			SelectedLines(
				Contains("+one"),
				Contains("+two"),
				Contains("+three"),
				Contains("+four"),
			)

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		// The moved line ('one') is gone from the commit, and the stale multi-line range
		// collapses to a single surviving change at the selection's top ordinal (now the
		// 'two' line) — just like discarding from the commit does.
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("+one")).
			ContainsLines(
				Equals("+two"),
				Equals("+three"),
				Equals("+four"),
				Equals("+five"),
			).
			SelectedLines(
				Contains("+two"),
			)
	},
})
