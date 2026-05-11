package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickCommitThatBecomesEmpty = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry-pick a commit that becomes empty at the destination",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			CreateFileAndAdd("file1", "change 1\n").
			CreateFileAndAdd("file2", "change 2\n").
			Commit("two changes in one commit").
			NewBranchFrom("branch", "HEAD^").
			CreateFileAndAdd("file1", "change 1\n").
			Commit("single change").
			CreateFileAndAdd("file3", "change 3\n").
			Commit("unrelated change").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("branch"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("unrelated change").IsSelected(),
				Contains("single change"),
				Contains("base"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("2 commits copied"))
			})

		t.Views().Commits().
			Focus().
			Lines(
				Contains("two changes in one commit").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
					Confirm()
			})

		if t.Git().Version().IsAtLeast(2, 45, 0) {
			t.Views().Commits().
				Lines(
					Contains("unrelated change"),
					Contains("single change"),
					Contains("two changes in one commit").IsSelected(),
					Contains("base"),
				).
				SelectPreviousItem()

			// Cherry-picked commit is empty
			t.Views().Main().Content(DoesNotContain("diff --git"))
		} else {
			t.Views().Commits().
				// We have a bug with how the selection is updated in this case; normally you would
				// expect the "two changes in one commit" commit to be selected because it was
				// selected before pasting, and we try to maintain that selection. This is broken
				// for two reasons:
				// 1. We increment the selected line index after pasting by the number of pasted
				// commits; this is wrong because we skipped the commit that became empty. So
				// according to this bug, the "base" commit should be selected.
				// 2. We only update the selected line index after pasting if the currently selected
				// commit is not a rebase TODO commit, on the assumption that if it is, we are in a
				// rebase and the cherry-picked commits end up below the selection. In this case,
				// however, we still think we are cherry-picking because the final refresh after the
				// CheckMergeOrRebase in CherryPickHelper.Paste is async and hasn't completed yet;
				// so the "unrelated change" still has a "pick" action.
				//
				// Since this only happens for older git versions, we don't bother fixing it.
				Lines(
					Contains("unrelated change").IsSelected(),
					Contains("two changes in one commit"),
					Contains("base"),
				)
		}
	},
})
