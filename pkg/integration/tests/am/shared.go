package am

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// setUpConflictingAm creates two commits on a branch, turns them into patches,
// and applies these to master with "git am". The first patch conflicts. If
// secondPatchConflicts is true, the second one conflicts with the resolution of
// the first one too. With threeWay, "git am" gets the --3way option, so that
// the conflicts show up in the files; otherwise it leaves the files alone.
func setUpConflictingAm(shell *Shell, secondPatchConflicts bool, threeWay bool) {
	shell.
		CreateFileAndAdd("file", "original\n").
		Commit("base").
		NewBranch("source").
		UpdateFileAndAdd("file", "source\n").
		Commit("change file")

	if secondPatchConflicts {
		shell.UpdateFileAndAdd("file", "source again\n").Commit("change file again")
	} else {
		shell.CreateFileAndAdd("file2", "two\n").Commit("add file2")
	}

	shell.
		RunShellCommand("git format-patch --stdout master..source > ../patches.mbox").
		Checkout("master").
		UpdateFileAndAdd("file", "master\n").
		Commit("conflicting change")

	amArgs := []string{"git", "am"}
	if threeWay {
		amArgs = append(amArgs, "--3way")
	}
	shell.RunCommandExpectError(append(amArgs, "../patches.mbox"))
}
