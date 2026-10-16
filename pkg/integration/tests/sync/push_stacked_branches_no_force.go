package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesNoForce = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a stack of branches that are only ahead of their remote branches, without being asked to force push",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("base")
		shell.NewBranch("branch1")
		shell.NewBranch("branch2")
		shell.NewBranch("branch3")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.SetBranchUpstream("branch1", "origin/branch1")
		shell.SetBranchUpstream("branch2", "origin/branch2")
		shell.SetBranchUpstream("branch3", "origin/branch3")

		shell.Checkout("branch1")
		shell.EmptyCommit("one")
		shell.Checkout("branch2")
		shell.HardReset("branch1")
		shell.EmptyCommit("two")
		shell.Checkout("branch3")
		shell.HardReset("branch2")
		shell.EmptyCommit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("branch3 ↑3"),
				Contains("branch1 ↑1"),
				Contains("branch2 ↑2"),
				Contains("master ✓"),
			)

		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Menu().
			Title(Equals("Push")).
			ContainsLines(
				Contains("  branch2 ↑2"),
				Contains("  branch1 ↑1"),
			).
			Select(Contains("Push all these branches in addition to the current one")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ✓"),
				Contains("branch1 ✓"),
				Contains("branch2 ✓"),
				Contains("master ✓"),
			)
	},
})
