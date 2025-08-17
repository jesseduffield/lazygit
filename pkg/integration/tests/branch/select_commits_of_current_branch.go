package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectCommitsOfCurrentBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select all commits of the current branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("master 01")
		shell.EmptyCommit("master 02")
		shell.NewBranch("branch1")
		shell.CreateNCommits(2)
		shell.NewBranchFrom("branch2", "master")
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master 02"),
				Contains("master 01"),
			).
			Press(keys.Commits.SelectCommitsOfCurrentBranch).
			Lines(
				Contains("commit 03").IsSelected(),
				Contains("commit 02").IsSelected(),
				Contains("commit 01").IsSelected(),
				Contains("master 02"),
				Contains("master 01"),
			).
			PressEscape().
			Lines(
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master 02"),
				Contains("master 01"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch2").IsSelected(),
				Contains("branch1"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
				Contains("master 02"),
				Contains("master 01"),
			).
			Press(keys.Commits.SelectCommitsOfCurrentBranch).
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01").IsSelected(),
				Contains("master 02"),
				Contains("master 01"),
			)
	},
})
