package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchAndAutoForwardBranchesWorktreeAddedAfterStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Auto-forward skips a main branch that was externally checked out in a linked worktree after lazygit started",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoForwardBranches = "onlyMainBranches"
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
		shell.NewBranch("feature")
		shell.NewBranch("wt-branch")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.SetBranchUpstream("feature", "origin/feature")
		shell.Checkout("master")
		shell.HardReset("HEAD^")
		shell.Checkout("feature")
		shell.AddWorktreeCheckout("wt-branch", "../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("feature").IsSelected(),
				Contains("master ↓1").DoesNotContain("↑"),
				Contains("wt-branch (worktree linked-worktree)"),
			)

		// Switch the linked worktree to master externally.
		t.Shell().RunCommand([]string{"git", "-C", "../linked-worktree", "checkout", "master"})

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		t.Views().Branches().
			Lines(
				Contains("feature").IsSelected(),
				Contains("master (worktree linked-worktree) ↓1"),
				Contains("wt-branch").DoesNotContain("worktree"),
			)

		t.Views().Worktrees().
			Focus().
			NavigateToLine(Contains("linked-worktree")).
			PressPrimaryAction()

		t.Views().Files().
			Focus().
			IsEmpty()
	},
})
