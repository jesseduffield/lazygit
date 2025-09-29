package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickEmptyAutoStash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry-picking an empty commit after auto-stashing restores stash list and copy mode",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("source").
			CreateFileAndAdd("shared.txt", "content\n").Commit("add shared file on source").
			Checkout("master").
			CreateFileAndAdd("shared.txt", "content\n").Commit("add shared file on master").
			UpdateFile("shared.txt", "local change\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("shared.txt"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("master"),
				Contains("source"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("add shared file on source").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("1 commit copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("add shared file on master").IsSelected(),
				Contains("base"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Cherry-pick produced no changes")).
					ContainsLines(
						Contains("Skip this cherry-pick"),
						Contains("Create empty commit and continue"),
						Contains("Cancel"),
					).
					Select(Contains("Skip this cherry-pick")).
					Confirm()
			}).
			Tap(func() {
				t.Shell().RunCommandExpectError([]string{"git", "rev-parse", "CHERRY_PICK_HEAD"})
				t.Views().Stash().IsEmpty()
				t.Views().Files().Lines(
					Contains("shared.txt"),
				)
				t.Views().Information().Content(DoesNotContain("commit copied"))
			}).
			Lines(
				Contains("add shared file on master").IsSelected(),
				Contains("base"),
			)
	},
})
