package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveLinesFromMainViewSecondary = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a line from the custom patch by pressing space on it in the secondary (custom-patch) pane of the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Line mode, so we can include only some of a hunk's additions and exercise the
		// renumbering the aggregated patch applies to the included ones.
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		// Three consecutive additions in a single hunk.
		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "one\nADDED1\nADDED2\nADDED3\ntwo\nthree\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextItem).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		// Include ADDED2 and ADDED3 in the patch but not ADDED1. Excluding ADDED1 shifts
		// the new-file line numbers of ADDED2/ADDED3 in the aggregated patch, which is what
		// used to make removing them by line number from the secondary act on the wrong line.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADDED1"),
			).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+ADDED2"),
			).
			Press(keys.Universal.RangeSelectDown).
			SelectedLines(
				Contains("+ADDED2"),
				Contains("+ADDED3"),
			).
			PressPrimaryAction()

		// The cumulative patch holds ADDED2 and ADDED3 only.
		t.Views().Secondary().
			Content(Contains("+ADDED2")).
			Content(Contains("+ADDED3")).
			Content(DoesNotContain("ADDED1"))

		// Tab into the secondary (custom-patch) pane; the selection lands on its first
		// change line, ADDED2.
		t.Views().Main().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+ADDED2"),
			).
			// Discarding from the commit makes no sense in the custom-patch preview, so it's
			// disabled here (you remove from the patch with space instead).
			Press(keys.Universal.Remove)

		t.ExpectToast(Contains("Cannot discard from the custom patch view"))

		t.Views().Secondary().
			// Space removes the selected line from the patch — and removes ADDED2, not some
			// other line resolved from its shifted line number.
			PressPrimaryAction().
			// The selection lands on the next surviving change, ADDED3.
			SelectedLines(
				Contains("+ADDED3"),
			).
			Content(Contains("+ADDED3")).
			Content(DoesNotContain("ADDED2")).
			Content(DoesNotContain("ADDED1"))

		// Applying confirms only ADDED3 was ever in the patch.
		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		t.Views().Files().
			Focus().
			Lines(
				Contains("file1").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("ADDED3")).
			Content(DoesNotContain("ADDED1")).
			Content(DoesNotContain("ADDED2"))
	},
})
