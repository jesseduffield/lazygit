package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchAndAutoForwardBranchesAllBranchesCheckedOutInOtherWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fetch from remote and auto-forward branches with config set to 'allBranches'; check that a branch checked out by another worktree is forwarded there, unless that worktree has changes or is mid-rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoForwardBranches = "allBranches"
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
		shell.NewBranch("feature")
		shell.NewBranch("diverged")
		shell.NewBranch("dirty")
		shell.NewBranch("rebasing")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.SetBranchUpstream("feature", "origin/feature")
		shell.SetBranchUpstream("diverged", "origin/diverged")
		shell.SetBranchUpstream("dirty", "origin/dirty")
		shell.SetBranchUpstream("rebasing", "origin/rebasing")
		shell.Checkout("master")
		shell.HardReset("HEAD^")
		shell.Checkout("feature")
		shell.HardReset("HEAD~2")
		shell.Checkout("diverged")
		shell.HardReset("HEAD~2")
		shell.EmptyCommit("local")
		shell.Checkout("dirty")
		shell.HardReset("HEAD~2")
		shell.Checkout("rebasing")
		shell.HardReset("HEAD^")
		shell.NewBranch("checked-out")

		shell.AddWorktreeCheckout("master", "../linked-worktree")

		shell.AddWorktreeCheckout("dirty", "../dirty-worktree")
		shell.UpdateFile("../dirty-worktree/file01.txt", "changed")

		shell.AddWorktreeCheckout("rebasing", "../rebasing-worktree")
		// the failing exec stops the rebase after picking commit-02, with HEAD
		// detached from the branch
		shell.RunCommandExpectError([]string{"git", "-C", "../rebasing-worktree", "rebase", "--exec", "false", "HEAD^"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("checked-out").IsSelected(),
				Contains("dirty (worktree dirty-worktree) ↓2").DoesNotContain("↑"),
				Contains("diverged ↓2↑1"),
				Contains("feature ↓2").DoesNotContain("↑"),
				Contains("master (worktree linked-worktree) ↓1").DoesNotContain("↑"),
				Contains("rebasing (worktree rebasing-worktree) ↓1").DoesNotContain("↑"),
			)

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		// AutoForwardBranches is "allBranches": feature gets forwarded, and so
		// does master in the worktree it is checked out in
		t.Views().Branches().
			Lines(
				Contains("checked-out").IsSelected(),
				Contains("dirty (worktree dirty-worktree) ↓2"),
				Contains("diverged ↓2↑1"),
				Contains("feature ✓"),
				Contains("master (worktree linked-worktree) ✓"),
				Contains("rebasing (worktree rebasing-worktree) ↓1"),
			)

		// The files of the linked worktree have moved along with master
		t.Views().Worktrees().
			Focus().
			NavigateToLine(Contains("linked-worktree")).
			PressPrimaryAction()

		t.Views().Files().
			Focus().
			IsEmpty()

		// The rebase in progress was left alone
		t.Views().Worktrees().
			Focus().
			NavigateToLine(Contains("rebasing-worktree")).
			PressPrimaryAction()

		t.Views().Commits().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("─── Commits"),
				Contains("commit-02"),
				Contains("commit-01"),
			)
	},
})
