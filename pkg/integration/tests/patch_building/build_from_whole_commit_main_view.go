package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildFromWholeCommitMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a custom patch from the whole-commit diff in a commits panel's focused main view (without entering the commit files), then apply it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.CreateFileAndAdd("file2", "alpha\nbeta\ngamma\ndelta\n")
		shell.Commit("first commit")

		// One commit touching two files, so the whole-commit diff is multi-file.
		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\n")
		shell.UpdateFileAndAdd("file2", "alpha\nBETA\ngamma\ndelta\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.NextItem).
			PressEnter()

		// Focus the whole-commit diff straight from the sub-commits panel, rather than
		// entering the commit files panel first.
		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			// The selection anchors on the first change line of the multi-file diff,
			// which belongs to file1.
			SelectedLines(
				Contains("-three"),
			).
			// `a` extends to the whole change block, then space toggles just file1's block
			// into the custom patch; file2's change, also in this diff, stays out.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction().
			// The selection is re-established on the same block after the toggle's
			// re-render (which, when the secondary view first appears, re-wraps the
			// narrower diff).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			)

		t.Views().Information().Content(Contains("Building patch"))

		// The secondary view shows the cumulative patch live — only file1's toggled
		// block, not file2's change from the same commit.
		t.Views().Secondary().
			ContainsLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			Content(DoesNotContain("BETA"))

		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		// Only file1's toggled block reached the working tree; file2 is untouched.
		t.Views().Files().
			Focus().
			Lines(
				Contains("file1").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("THREE")).
			Content(DoesNotContain("BETA"))
	},
})
