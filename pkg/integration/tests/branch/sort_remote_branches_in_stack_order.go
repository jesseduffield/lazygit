package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortRemoteBranchesInStackOrder = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sort remote branches that share a committer date so that a branch comes before the ones it is based on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.41.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.RemoteBranchSortOrder = "date"
	},
	SetupRepo: func(shell *Shell) {
		// Rebasing a stack of branches gives every commit it creates the same
		// committer date; give these commits one date for the same effect.
		date := "2024-01-01 10:00:00"

		shell.
			EmptyCommitWithDate("base", date).
			NewBranch("branch-c").
			EmptyCommitWithDate("c", date).
			NewBranch("branch-a").
			EmptyCommitWithDate("a", date).
			NewBranch("branch-b").
			EmptyCommitWithDate("b", date).
			NewBranchFrom("unrelated", "master").
			EmptyCommitWithDate("unrelated", date).
			Checkout("master").
			CloneIntoRemote("origin")
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
				Contains("branch-b").IsSelected(),
				Contains("branch-a"),
				Contains("branch-c"),
				Contains("unrelated"),
				Contains("master"),
			)
	},
})
