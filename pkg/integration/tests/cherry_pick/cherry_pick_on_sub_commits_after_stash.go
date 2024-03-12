package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickOnSubCommitsAfterStash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste cherry-picked commits onto sub-commits view after stash",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("first-branch").
			NewBranch("second-branch").
			Checkout("first-branch").
			EmptyCommit("one").
			Checkout("second-branch").
			CreateFileAndAdd("file-unstaged", "content").
			EmptyCommit("two").
			UpdateFile("file-unstaged", "new content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("second-branch"),
				Contains("first-branch"),
				Contains("master"),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("two").IsSelected(),
				Contains("base"),
			).
			// copy commit 'two'
			Press(keys.Commits.CherryPickCopy).
			Tap(func() {
				t.Views().Information().Content(Contains("1 commit copied"))
			})

		t.Views().Branches().
			Focus().
			Lines(
				Contains("second-branch").IsSelected(),
				Contains("first-branch"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("one").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				// cherry-picked commit will be deleted after confirmation
				t.Views().Information().Content(Contains("1 commit copied"))
			}).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Stash changes")).
					Content(Contains("Are you sure you want to stash all changes?")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the copied commits onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commit copied"))
			}).
			PressEscape()

		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-branch").IsSelected(),
				Contains("second-branch"),
				Contains("master"),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
				Contains("base"),
			)

		t.Views().Stash().
			Lines(
				MatchesRegexp(`\ds .* Stash all changes`),
			)
	},
})
