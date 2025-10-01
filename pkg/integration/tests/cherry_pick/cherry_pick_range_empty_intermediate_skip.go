package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickRangeEmptyIntermediateSkip = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry-picking a range with an intermediate empty commit and skipping it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("shared.txt", "original\n").Commit("add shared file on master").
			NewBranch("target").
			UpdateFileAndAdd("shared.txt", "target change\n").Commit("update shared file on target").
			Checkout("master").
			NewBranch("source").
			CreateFileAndAdd("unique1.txt", "content1\n").Commit("add unique1 on source").
			UpdateFileAndAdd("shared.txt", "target change\n").Commit("match target change").
			CreateFileAndAdd("unique2.txt", "content2\n").Commit("add unique2 on source").
			Checkout("target")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("target"),
				Contains("source"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("add unique2 on source").IsSelected(),
				Contains("match target change"),
				Contains("add unique1 on source"),
				Contains("add shared file on master"),
			).
			Press(keys.Universal.RangeSelectDown).
			Lines(
				Contains("add unique2 on source").IsSelected(),
				Contains("match target change").IsSelected(),
				Contains("add unique1 on source"),
				Contains("add shared file on master"),
			).
			Press(keys.Universal.RangeSelectDown).
			Lines(
				Contains("add unique2 on source").IsSelected(),
				Contains("match target change").IsSelected(),
				Contains("add unique1 on source").IsSelected(),
				Contains("add shared file on master"),
			).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("3 commits copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("update shared file on target").IsSelected(),
				Contains("add shared file on master"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 3 copied commit(s) onto this branch?")).
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
				t.Views().Information().Content(DoesNotContain("commits copied"))
			}).
			TopLines(
				Contains("add unique2 on source").IsSelected(),
				Contains("add unique1 on source"),
				Contains("update shared file on target"),
			)
	},
})
