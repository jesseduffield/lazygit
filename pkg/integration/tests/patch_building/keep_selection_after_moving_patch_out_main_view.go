package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectionAfterMovingPatchOutMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a custom patch out of a commit leaves the focused main view's selection on a change that is still there, rather than painted over the diff the rewrite left behind",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
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

		// Take the first of the commit's three changed lines into a custom patch.
		t.Views().CommitFiles().
			IsFocused().
			PressEnter()

		t.Views().PatchBuilding().
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
			PressPrimaryAction().
			Press(keys.Universal.Return)

		// Leave a range selected over the diff, spanning the lines the patch holds. The
		// patch move doesn't go through the main view at all, so without a net nothing
		// would move this selection off lines that the rewrite takes away.
		t.Views().CommitFiles().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-one"),
			).
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
		// change that has taken its place — the same place in the diff's changes, which
		// is where the user was.
		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("+ONE")).
			SelectedLines(
				Contains("-three"),
			)
	},
})
