package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick a merge commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("first-branch").
			NewBranch("second-branch").
			Checkout("first-branch").
			Checkout("second-branch").
			CreateFileAndAdd("file1.txt", "content").
			Commit("one").
			CreateFileAndAdd("file2.txt", "content").
			Commit("two").
			Checkout("master").
			Merge("second-branch").
			Checkout("first-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-branch"),
				Contains("master"),
				Contains("second-branch"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("⏣─╮ Merge branch 'second-branch'").IsSelected(),
				Contains("│ ◯ two"),
				Contains("│ ◯ one"),
				Contains("◯ ╯ base"),
			).
			// copy the merge commit
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("1 commit copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base").IsSelected(),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commit copied"))
			}).
			Lines(
				Contains("Merge branch 'second-branch'").IsSelected(),
				Contains("base"),
			)

		t.Views().Main().ContainsLines(
			Contains("Merge branch 'second-branch'"),
			Contains("---"),
			Contains("file1.txt | 1 +"),
			Contains("file2.txt | 1 +"),
			Contains("2 files changed, 2 insertions(+)"),
		)
	},
})
