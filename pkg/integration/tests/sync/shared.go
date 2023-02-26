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

func assertSuccessfullyPushed(t *TestDriver) {
	t.Views().Status().Content(Contains("✓ repo → master"))

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
