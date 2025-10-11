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
				Contains("single change").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("1 commit copied"))
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
					Content(Contains("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Lines(
				Contains("single change"),
				Contains("two changes in one commit").IsSelected(),
				Contains("base"),
			).
			SelectPreviousItem()

		// Cherry-picked commit is empty
		t.Views().Main().Content(DoesNotContain("diff --git"))
	},
})
