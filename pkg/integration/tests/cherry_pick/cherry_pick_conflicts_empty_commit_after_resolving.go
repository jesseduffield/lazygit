package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CherryPickConflictsEmptyCommitAfterResolving = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick commits with conflicts, resolve them so that the commit becomes empty",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("first change").IsSelected(),
			).
			Press(keys.Commits.PasteCommits)

		t.ExpectPopup().Alert().
			Title(Equals("Cherry-pick")).
			Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file")).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Discard changes")).
			Select(Contains("Discard all changes")).
			Confirm()

		t.Common().ContinueOnConflictsResolved("cherry-pick")

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Focus().
			TopLines(
				// We have a bug with how the selection is updated in this case; normally you would
				// expect the "first change" commit to be selected because it was selected before
				// pasting, and we try to maintain that selection. This is broken for two reasons:
				// 1. We increment the selected line index after pasting by the number of pasted
				// commits; this is wrong because we skipped the commit that became empty. So
				// according to this bug, the "original" commit should be selected.
				// 2. We only update the selected line index after pasting if the currently selected
				// commit is not a rebase TODO commit, on the assumption that if it is, we are in a
				// rebase and the cherry-picked commits end up below the selection. In this case,
				// however, we still think we are cherry-picking because the final refresh after the
				// CheckMergeOrRebase in CherryPickHelper.Paste is async and hasn't completed yet;
				// so the "second-change-branch unrelated change" still has a "pick" action.
				//
				// We don't bother fixing it for now because it's a pretty niche case, and the
				// nature of the problem is only cosmetic.
				Contains("second-change-branch unrelated change").IsSelected(),
				Contains("first change"),
				Contains("original"),
			)
	},
})
