package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickEmptyFollowedByConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry-picking multiple commits skips an empty commit and only cleans up after resolving subsequent conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("shared.txt", "base\n").Commit("add shared base").
			CreateFileAndAdd("conflict.txt", "base version\n").Commit("add conflict base").
			NewBranch("target").
			UpdateFileAndAdd("shared.txt", "target change\n").Commit("update shared on target").
			UpdateFileAndAdd("conflict.txt", "target version\n").Commit("update conflict on target").
			Checkout("master").
			NewBranch("source").
			UpdateFileAndAdd("shared.txt", "target change\n").Commit("match target shared").
			UpdateFileAndAdd("conflict.txt", "source version\n").Commit("add conflict on source").
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
				Contains("add conflict on source").IsSelected(),
				Contains("match target shared"),
				Contains("add conflict base"),
				Contains("add shared base"),
			).
			Press(keys.Universal.RangeSelectDown).
			Lines(
				Contains("add conflict on source").IsSelected(),
				Contains("match target shared").IsSelected(),
				Contains("add conflict base"),
				Contains("add shared base"),
			).
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("update conflict on target").IsSelected(),
				Contains("update shared on target"),
				Contains("add conflict base"),
				Contains("add shared base"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
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
				t.Shell().RunCommand([]string{"git", "rev-parse", "CHERRY_PICK_HEAD"})
			})

		t.Common().AcknowledgeConflicts()

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("conflict.txt")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			LineCount(EqualsInt(2)).
			Lines(
				Contains("target version"),
				Contains("source version"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Common().ContinueOnConflictsResolved("cherry-pick")

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("add conflict on source").IsSelected(),
				Contains("match target shared"),
				Contains("update conflict on target"),
			).
			SelectedLine(Contains("add conflict on source"))

		t.Views().Information().Content(DoesNotContain("commit copied"))

		t.Shell().RunCommandExpectError([]string{"git", "rev-parse", "CHERRY_PICK_HEAD"})
	},
})
