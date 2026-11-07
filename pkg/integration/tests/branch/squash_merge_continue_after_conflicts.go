package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var SquashMergeContinueAfterConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash merge and commit a branch, and continue after resolving the conflicts",
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
			Select(Contains("Squash merge and commit")).
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Select(Contains("View conflicts")).
			Confirm()

		t.Views().Status().Content(Equals("(merging) repo → first-change-branch"))

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved content")
			}).
			Press(keys.Universal.Refresh)

		t.Common().ContinueOnConflictsResolved("merge")

		t.Views().Status().Content(Equals("repo → first-change-branch"))

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Lines(
				Contains("Squash merge second-change-branch into first-change-branch").IsSelected(),
				Contains("first change"),
				Contains("original"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			)

		t.FileSystem().PathNotPresent(".git/SQUASH_MSG")
	},
})
