package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildMultiFileFromWholeCommitMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a custom patch spanning two files from a commit's whole-commit diff in the focused main view, by toggling each file's hunk, then apply it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.CreateFileAndAdd("file2", "alpha\nbeta\ngamma\ndelta\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\n")
		shell.UpdateFileAndAdd("file2", "alpha\nBETA\ngamma\ndelta\n")
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

		t.Views().Main().
			IsFocused().
			// Hunk mode is the default here, so the first change block of the multi-file
			// diff (file1's) is selected on focus.
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Move to the next change block, which is in file2, and toggle it in too.
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-beta"),
				Contains("+BETA"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("-beta"),
				Contains("+BETA"),
			)

		// The cumulative patch now spans both files.
		t.Views().Secondary().
			ContainsLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			ContainsLines(
				Contains("-beta"),
				Contains("+BETA"),
			)

		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		// Both files' changes reached the working tree.
		t.Views().Files().
			Focus().
			ContainsLines(
				Contains("file1"),
			).
			ContainsLines(
				Contains("file2"),
			)
	},
})
