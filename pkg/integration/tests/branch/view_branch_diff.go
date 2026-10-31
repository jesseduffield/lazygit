package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ViewBranchDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the files a branch changed since it diverged from its base branch, leaving out what it merged in from there",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("base-file", "base\n").
			Commit("main 1").
			RenameCurrentBranch("main").
			NewBranch("feature").
			CreateFileAndAdd("dir/feature-file-1", "feature 1\n").
			Commit("feature 1").
			Checkout("main").
			CreateFileAndAdd("merged-file", "merged\n").
			Commit("main 2").
			Checkout("feature").
			Merge("main").
			CreateFileAndAdd("dir/feature-file-2", "feature 2\n").
			Commit("feature 2").
			Checkout("main").
			CreateFileAndAdd("unmerged-file", "unmerged\n").
			Commit("main 3").
			Checkout("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature").IsSelected(),
				Contains("main"),
			).
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature since main)")).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  A feature-file-1"),
				Equals("  A feature-file-2"),
			)

		// The diff stays as it was when it was asked for, even as the branch moves on and
		// the files are listed again
		t.Shell().
			CreateFileAndAdd("dir/feature-file-3", "feature 3\n").
			Commit("feature 3")

		t.Views().CommitFiles().
			IsFocused().
			Press(keys.Universal.IncreaseRenameSimilarityThreshold)

		t.ExpectToast(Contains("Changed rename similarity threshold"))

		t.Views().CommitFiles().
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  A feature-file-1"),
				Equals("  A feature-file-2"),
			).
			PressEscape()

		t.Views().Branches().
			IsFocused().
			NavigateToLine(Contains("main")).
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (main since main)")).
			Lines(
				Contains("(none)"),
			).
			PressEscape()

		// Asking again picks up the commit made since, and the upstream options offer
		// the branch diff too
		t.Views().Branches().
			IsFocused().
			NavigateToLine(Contains("feature")).
			Press(keys.Branches.SetUpstream)

		t.ExpectPopup().Menu().Title(Contains("Upstream")).
			Select(Contains("View branch diff since main")).Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature since main)")).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  A feature-file-1"),
				Equals("  A feature-file-2"),
				Equals("  A feature-file-3"),
			)
	},
})
