package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ViewBranchDiffFromCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the branch diff of the checked-out branch from its commits, and of another branch from its sub-commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("base-file", "base\n").
			Commit("main 1").
			RenameCurrentBranch("main").
			NewBranch("other").
			CreateFileAndAdd("other-file", "other\n").
			Commit("other 1").
			Checkout("main").
			NewBranch("feature").
			CreateFileAndAdd("feature-file-1", "feature 1\n").
			Commit("feature 1").
			CreateFileAndAdd("feature-file-2", "feature 2\n").
			Commit("feature 2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature since main)")).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  A feature-file-1"),
				Equals("  A feature-file-2"),
			)

		// Refreshing the commits doesn't make the files follow the selected commit
		t.Shell().
			CreateFileAndAdd("feature-file-3", "feature 3\n").
			Commit("feature 3")
		t.RefreshInBackground()

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature since main)")).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  A feature-file-1"),
				Equals("  A feature-file-2"),
			).
			PressEscape()

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("other")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (other since main)")).
			Lines(
				Equals("A other-file").IsSelected(),
			)
	},
})
