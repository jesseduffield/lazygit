package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchAndAutoForwardBranchesCheckedOutInOtherWorktreeWithStaleSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Auto-forward a main branch that is checked out in another worktree whose submodule is checked out at a different commit than the one recorded",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoForwardBranches = "onlyMainBranches"
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoSubmodule("sub", "sub")
		shell.Commit("add submodule")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.HardReset("HEAD^")

		shell.NewBranch("other")
		shell.AddWorktreeCheckout("master", "../linked-worktree")
		shell.RunCommand([]string{"git", "-C", "../linked-worktree", "submodule", "update", "--init"})
		shell.RunCommand([]string{"git", "-C", "../linked-worktree/sub", "commit", "--allow-empty", "-m", "newer"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("other").IsSelected(),
				Contains("master (worktree linked-worktree) ↓1"),
			)

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		t.Views().Branches().
			Lines(
				Contains("other").IsSelected(),
				Contains("master (worktree linked-worktree) ✓"),
			)
	},
})
