package reflog

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick a reflog commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.HardReset("HEAD^^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().ReflogCommits().
			Focus().
			Lines(
				Contains("reset: moving to HEAD^^").IsSelected(),
				Contains("commit: three"),
				Contains("commit: two"),
				Contains("commit (initial): one"),
			).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("1 commit copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("one").IsSelected(),
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
					Select(Contains("Create empty commit and continue")).
					Confirm()
			}).
			Tap(func() {
				t.Shell().RunCommandExpectError([]string{"git", "rev-parse", "CHERRY_PICK_HEAD"})
			}).
			Lines(
				Contains("three"),
				Contains("one").IsSelected(),
			)
	},
})
