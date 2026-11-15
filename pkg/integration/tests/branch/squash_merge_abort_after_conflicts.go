package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var SquashMergeAbortAfterConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash merge and commit a branch, and abort when there are conflicts",
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
			Select(Contains("Abort the merge")).
			Confirm()

		t.Views().Status().Content(Equals("repo → first-change-branch"))

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			TopLines(
				Contains("first change"),
				Contains("original"),
			)

		t.FileSystem().PathNotPresent(".git/SQUASH_MSG")
	},
})
