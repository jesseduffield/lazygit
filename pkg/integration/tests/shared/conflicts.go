package shared

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OriginalFileContent = `
This
Is
The
Original
File
`

var FirstChangeFileContent = `
This
Is
The
First Change
File
`

var SecondChangeFileContent = `
This
Is
The
Second Change
File
`

// prepares us for a rebase/merge that has conflicts
var MergeConflictsSetup = func(shell *Shell) {
	shell.
		NewBranch("original-branch").
		EmptyCommit("one").
		EmptyCommit("two").
		EmptyCommit("three").
		CreateFileAndAdd("file", OriginalFileContent).
		Commit("original").
		NewBranch("first-change-branch").
		UpdateFileAndAdd("file", FirstChangeFileContent).
		Commit("first change").
		Checkout("original-branch").
		NewBranch("second-change-branch").
		UpdateFileAndAdd("file", SecondChangeFileContent).
		Commit("second change").
		EmptyCommit("second-change-branch unrelated change").
		Checkout("first-change-branch")
}
