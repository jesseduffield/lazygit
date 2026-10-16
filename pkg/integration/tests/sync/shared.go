package sync

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func createTwoBranchesReadyToForcePush(shell *Shell) {
	shell.EmptyCommit("one")
	shell.EmptyCommit("two")

	shell.NewBranch("other_branch")

	shell.CloneIntoRemote("origin")

	shell.SetBranchUpstream("master", "origin/master")
	shell.SetBranchUpstream("other_branch", "origin/other_branch")

	// remove the 'two' commit so that we have something to pull from the remote
	shell.HardReset("HEAD^")

	shell.Checkout("master")
	// doing the same for master
	shell.HardReset("HEAD^")
}

// Creates the branches branch1, branch2 and branch3, each on top of the
// previous one and pushed to origin, and then rewrites their commits the way
// rebasing the stack would, so that all of them have diverged from their
// remote branches. Also creates local-only at the tip of branch1, a branch
// without an upstream. Leaves branch3 checked out.
func createRebasedStackOfBranches(shell *Shell) {
	shell.EmptyCommit("base")
	shell.NewBranch("branch1")
	shell.EmptyCommit("one")
	shell.NewBranch("branch2")
	shell.EmptyCommit("two")
	shell.NewBranch("branch3")
	shell.EmptyCommit("three")

	shell.CloneIntoRemote("origin")
	shell.SetBranchUpstream("master", "origin/master")
	shell.SetBranchUpstream("branch1", "origin/branch1")
	shell.SetBranchUpstream("branch2", "origin/branch2")
	shell.SetBranchUpstream("branch3", "origin/branch3")

	shell.Checkout("branch1")
	shell.HardReset("master")
	shell.EmptyCommit("one-rebased")
	shell.NewBranch("local-only")
	shell.Checkout("branch2")
	shell.HardReset("branch1")
	shell.EmptyCommit("two-rebased")
	shell.Checkout("branch3")
	shell.HardReset("branch2")
	shell.EmptyCommit("three-rebased")
}

// Creates the branch "feature" with two commits on top of master, pushes it,
// and then rewrites those two commits and force-pushes them the way somebody
// else rebasing the branch would. The local branch stays where it was, so it
// has diverged from its remote branch without having any commits of its own.
// Leaves master checked out.
func createBranchRewrittenOnTheRemote(shell *Shell) {
	shell.EmptyCommit("one")
	shell.NewBranch("feature")
	shell.EmptyCommit("two")
	shell.EmptyCommit("three")

	shell.CloneIntoRemote("origin")
	shell.SetBranchUpstream("feature", "origin/feature")

	shell.CreateLightweightTag("before-rewrite", "feature")
	shell.HardReset("master")
	shell.EmptyCommit("two-rewritten")
	shell.EmptyCommit("three-rewritten")
	shell.RunCommand([]string{"git", "push", "--force", "origin", "feature"})
	shell.HardReset("before-rewrite")
	shell.RunCommand([]string{"git", "tag", "-d", "before-rewrite"})

	shell.Checkout("master")
}

// Creates the branches branch1, branch2 and branch3, each on top of the
// previous one and pushed to origin, and then rewrites their commits and
// force-pushes them the way somebody else rebasing the stack would. The local
// branches stay where they were, so all of them have diverged from their
// remote branches without having any commits of their own. Leaves master
// checked out.
func createStackRewrittenOnTheRemote(shell *Shell) {
	shell.EmptyCommit("base")
	shell.NewBranch("branch1")
	shell.EmptyCommit("one")
	shell.NewBranch("branch2")
	shell.EmptyCommit("two")
	shell.NewBranch("branch3")
	shell.EmptyCommit("three")

	shell.CloneIntoRemote("origin")
	shell.SetBranchUpstream("branch1", "origin/branch1")
	shell.SetBranchUpstream("branch2", "origin/branch2")
	shell.SetBranchUpstream("branch3", "origin/branch3")

	shell.CreateLightweightTag("before-rewrite", "branch3")
	shell.Checkout("branch1")
	shell.HardReset("master")
	shell.EmptyCommit("one-rewritten")
	shell.Checkout("branch2")
	shell.HardReset("branch1")
	shell.EmptyCommit("two-rewritten")
	shell.Checkout("branch3")
	shell.HardReset("branch2")
	shell.EmptyCommit("three-rewritten")
	shell.RunCommand([]string{"git", "push", "--force", "origin", "branch1", "branch2", "branch3"})

	// Put the local branches back where they were
	shell.HardReset("before-rewrite")
	shell.RunCommand([]string{"git", "branch", "-f", "branch2", "before-rewrite~"})
	shell.RunCommand([]string{"git", "branch", "-f", "branch1", "before-rewrite~2"})
	shell.RunCommand([]string{"git", "tag", "-d", "before-rewrite"})

	shell.Checkout("master")
}

func assertSuccessfullyPushed(t *TestDriver) {
	t.Views().Status().Content(Equals("✓ repo → master"))

	t.Views().Remotes().
		Focus().
		Lines(
			Contains("origin"),
		).
		PressEnter()

	t.Views().RemoteBranches().
		IsFocused().
		Lines(
			Contains("master"),
		).
		PressEnter()

	t.Views().SubCommits().
		IsFocused().
		Lines(
			Contains("two"),
			Contains("one"),
		)
}
