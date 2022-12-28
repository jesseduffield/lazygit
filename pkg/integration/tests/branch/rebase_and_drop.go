package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseAndDrop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase onto another branch, deal with the conflicts. Also mark a commit to be dropped before continuing.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
		// addin a couple additional commits so that we can drop one
		shell.EmptyCommit("to remove")
		shell.EmptyCommit("to keep")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			TopLines(
				Contains("to keep"),
				Contains("to remove"),
				Contains("first change"),
				Contains("original"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch").IsSelected(),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Confirmation().
			Title(Equals("Rebasing")).
			Content(Contains("Are you sure you want to rebase 'first-change-branch' on top of 'second-change-branch'?")).
			Confirm()

		t.Views().Information().Content(Contains("rebasing"))

		t.ExpectPopup().Confirmation().
			Title(Equals("Auto-merge failed")).
			Content(Contains("Conflicts!")).
			Confirm()

		t.Views().Files().IsFocused().
			SelectedLine(MatchesRegexp("UU.*file"))

		t.Views().Commits().
			Focus().
			TopLines(
				MatchesRegexp(`pick.*to keep`).IsSelected(),
				MatchesRegexp(`pick.*to remove`),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			TopLines(
				MatchesRegexp(`pick.*to keep`),
				MatchesRegexp(`drop.*to remove`).IsSelected(),
				MatchesRegexp("YOU ARE HERE.*second-change-branch unrelated change"),
				MatchesRegexp("second change"),
				MatchesRegexp("original"),
			)

		t.Views().Files().
			Focus().
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			PressPrimaryAction()

		t.ExpectPopup().Confirmation().
			Title(Equals("continue")).
			Content(Contains("all merge conflicts resolved. Continue?")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("rebasing"))

		t.Views().Commits().TopLines(
			Contains("to keep"),
			Contains("second-change-branch unrelated change").IsSelected(),
			Contains("second change"),
			Contains("original"),
		)
	},
})
