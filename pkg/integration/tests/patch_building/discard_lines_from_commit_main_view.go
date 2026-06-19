package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardLinesFromCommitMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard a line from a local commit straight from the focused main view, without entering the commit files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("commit to remove from")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to remove from").IsSelected(),
				Contains("first commit"),
			).
			// Focus the commit's diff straight from the commits panel, rather than
			// entering the commit files panel first.
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+1st line"),
			).
			// Move down to the second line and discard just that one from the commit.
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+2nd line"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard lines from commit")).
			Content(Equals("Are you sure you want to discard the selected lines from this commit?")).
			Confirm()

		// After the rebase the commit keeps only the 1st and 3rd lines, and the selection
		// advances to the next surviving change (the 3rd line, now at the discarded line's
		// ordinal) rather than staying stale — like staging or discarding in the files panel.
		t.Views().Main().
			ContainsLines(
				Equals("+1st line"),
				Equals("+3rd line"),
			).
			Content(DoesNotContain("2nd line")).
			SelectedLines(
				Contains("+3rd line"),
			)
	},
})
