package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ViewRemoteBranchDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the files a remote branch changed since it diverged from its base branch, without checking it out",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("base-file", "base\n").
			Commit("main 1").
			RenameCurrentBranch("main").
			NewBranch("feature").
			CreateFileAndAdd("dir/feature-file", "feature\n").
			Commit("feature 1").
			Checkout("main").
			CreateFileAndAdd("main-file", "main\n").
			Commit("main 2").
			CloneIntoRemote("origin").
			RunCommand([]string{"git", "branch", "-D", "feature"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin").IsSelected(),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("feature").IsSelected(),
				Contains("main"),
			).
			Press(keys.Branches.ViewBranchDiff)

		t.Views().CommitFiles().
			IsFocused().
			Title(Equals("Diff files (origin/feature since origin/main)")).
			Lines(
				Equals("▼ dir").IsSelected(),
				Equals("  A feature-file"),
			).
			PressEscape()

		t.Views().RemoteBranches().
			IsFocused()
	},
})
