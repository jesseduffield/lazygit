package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var LockErrorPartwayThroughCherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick two commits while another git process takes index.lock partway through",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "original\n").
			Commit("base").
			NewBranch("source").
			CreateFileAndAdd("file1", "one\n").
			Commit("one").
			CreateFileAndAdd("file2", "two\n").
			Commit("two").
			Checkout("master")

		// The cherry-pick commits "one"; the lock then stops it when it tries
		// to pick "two".
		shared.HoldIndexLockAfterNextCommit(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
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
				Contains("two").IsSelected(),
				Contains("one"),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base").IsSelected(),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
					Confirm()
			}).
			/* EXPECTED:
			Lines(
				Contains("two"),
				Contains("one"),
				Contains("base"),
			)
			ACTUAL: */
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("cherry-pick is already in progress")).
					Confirm()
			}).
			Lines(
				Contains("─── Pending cherry-picks"),
				Contains("pick").Contains("two"),
				Contains("─── Commits"),
				Contains("one"),
				Contains("base"),
			)
	},
})
