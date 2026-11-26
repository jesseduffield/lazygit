package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestOnlyOverACommitsDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Opening a diff line in the pull request is offered over a commit's own diff, and refused over the other diffs the main view shows",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")

		shell.UpdateFile("file1", "one\nTWO\nTHREE\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The working tree's diff is no commit of a branch, so no pull request has a
		// view of it and the command isn't offered there at all.
		t.Views().Files().
			Focus().
			SelectedLine(Contains("file1")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Tap(func() {
				// The command is bound right below the one asserted here, so a menu
				// showing that one would be showing this one too if it had it.
				t.Views().Menu().
					Content(Contains("Go to next file")).
					Content(DoesNotContain("Open pull request at selected line"))
			}).
			Cancel()

		// Over a commit's diff it is offered, and says so where the branch has no pull
		// request to open.
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("second commit")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Commits.OpenPullRequestInBrowser)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("No pull request found for this branch")).
			Confirm()

		// The pane previewing the custom patch shows the patch's lines at the numbers
		// the patch gives them, which are not the ones the pull request shows.
		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Commits.OpenPullRequestInBrowser)

		t.ExpectToast(Contains("Not available for the custom patch"))

		// In diffing mode the main view shows a diff against another ref rather than
		// the commit's own, and the pull request has no view of that either.
		t.Views().Commits().
			Focus().
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Diffing")).
			Select(MatchesRegexp(`Diff \w+`)).
			Confirm()

		t.Views().Commits().
			SelectNextItem().
			SelectedLine(Contains("first commit")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Commits.OpenPullRequestInBrowser)

		t.ExpectToast(Contains("Not available in diffing mode"))
	},
})
