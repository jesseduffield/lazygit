package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var SquashMergeUncommittedWithConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash merge a branch without committing, resolve the conflicts, and continue the merge later",
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
				Contains("first-change-branch").IsSelected(),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and leave uncommitted")).
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Select(Contains("View conflicts")).
			Confirm()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved content")
			}).
			Press(keys.Universal.Refresh)

		// The user asked not to commit, so we don't offer to continue the merge
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file"),
			)

		t.Views().Status().Content(Equals("(merging) repo → first-change-branch"))

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Merge options")).
			Select(Contains("continue")).
			Confirm()

		t.Views().Status().Content(Equals("repo → first-change-branch"))

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			TopLines(
				Contains("Squash merge second-change-branch into first-change-branch"),
				Contains("first change"),
			)
	},
})
