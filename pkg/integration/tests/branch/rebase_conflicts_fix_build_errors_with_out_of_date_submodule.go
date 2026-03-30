package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseConflictsFixBuildErrorsWithOutOfDateSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts. While continue prompt is showing, fix build errors; get another prompt when continuing. Check that we don't stage submodules here.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		// Create an out-of-date submodule to verify that we don't try to stage it
		shell.
			EmptyCommit("Initial commit").
			CloneIntoSubmodule("submodule", "submodule").
			Commit("Add submodule").
			AddFileInWorktreeOrSubmodule("submodule", "file", "content").
			CommitInWorktreeOrSubmodule("submodule", "add file in submodule")

		shared.MergeConflictsSetup(shell)

		// Create an untracked file to verify that we don't try to stage it either
		shell.UpdateFile("untracked-file", "some untracked file")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().TopLines(
			Contains("first change"),
			Contains("original"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
				Contains("master"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'first-change-branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Rebasing"))

		popup := t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue the rebase?"))

		// While the popup is showing, fix some build errors
		t.Shell().UpdateFile("file", "make it compile again")

		// Continue
		popup.Confirm()

		t.Views().Files().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  MM file"),
				Equals("   M submodule (submodule)"),
				Equals("  ?? untracked-file"),
			)

		t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("Files have been modified since conflicts were resolved. Auto-stage them and continue?")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("Rebasing"))

		t.Views().Files().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("   M submodule (submodule)"),
				Equals("  ?? untracked-file"),
			)

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("first change").IsSelected(),
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
				Contains("original"),
			)

		t.Views().Main().
			Content(
				DoesNotContain("submodule").DoesNotContain("untracked-file"),
			)
	},
})
