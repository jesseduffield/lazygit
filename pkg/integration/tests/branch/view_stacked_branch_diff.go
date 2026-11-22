package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ViewStackedBranchDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the branch diff of a branch stacked on others, since the closest of them or any other",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("base-file", "base\n").
			Commit("main 1").
			RenameCurrentBranch("main").
			NewBranch("feature-a").
			CreateFileAndAdd("a-file", "a\n").
			Commit("a").
			NewBranch("feature-b").
			CreateFileAndAdd("b-file", "b\n").
			Commit("b").
			NewBranch("feature-c").
			CreateFileAndAdd("c-file", "c\n").
			Commit("c")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("feature-c")).
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature-c since feature-b)")).
			Lines(
				Equals("A c-file").IsSelected(),
			).
			PressEscape()

		t.Views().Branches().
			IsFocused().
			Press(keys.Branches.SetUpstream)

		t.ExpectPopup().Menu().Title(Contains("Upstream")).
			ContainsLines(
				Contains("View branch diff since feature-b"),
				Contains("View branch diff since feature-a"),
				Contains("View branch diff since main"),
			).
			Select(Contains("View branch diff since feature-a")).Confirm()

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (feature-c since feature-a)")).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  A b-file"),
				Equals("  A c-file"),
			)
	},
})
