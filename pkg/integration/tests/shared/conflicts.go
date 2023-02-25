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

var CreateMergeConflictFile = func(shell *Shell) {
	MergeConflictsSetup(shell)

	shell.RunShellCommandExpectError("git merge --no-edit second-change-branch")
}

var CreateMergeCommit = func(shell *Shell) {
	CreateMergeConflictFile(shell)
	shell.UpdateFileAndAdd("file", SecondChangeFileContent)
	shell.ContinueMerge()
}

// creates a merge conflict where there are two files with conflicts and a separate file without conflicts
var CreateMergeConflictFiles = func(shell *Shell) {
	shell.
		NewBranch("original-branch").
		EmptyCommit("one").
		EmptyCommit("two").
		EmptyCommit("three").
		CreateFileAndAdd("file1", OriginalFileContent).
		CreateFileAndAdd("file2", OriginalFileContent).
		Commit("original").
		NewBranch("first-change-branch").
		UpdateFileAndAdd("file1", FirstChangeFileContent).
		UpdateFileAndAdd("file2", FirstChangeFileContent).
		Commit("first change").
		Checkout("original-branch").
		NewBranch("second-change-branch").
		UpdateFileAndAdd("file1", SecondChangeFileContent).
		UpdateFileAndAdd("file2", SecondChangeFileContent).
		// this file is not changed in the second branch
		CreateFileAndAdd("file3", "content").
		Commit("second change").
		EmptyCommit("second-change-branch unrelated change").
		Checkout("first-change-branch")

	shell.RunShellCommandExpectError("git merge --no-edit second-change-branch")
}

// These 'multiple' variants are just like the short ones but with longer file contents and with multiple conflicts within the file.

var OriginalFileContentMultiple = `
This
Is
The
Original
File
..
It
Is
Longer
Than
The
Other
Options
`

var FirstChangeFileContentMultiple = `
This
Is
The
First Change
File
..
It
Is
Longer
Than
The
Other
Other First Change
`

var SecondChangeFileContentMultiple = `
This
Is
The
Second Change
File
..
It
Is
Longer
Than
The
Other
Other Second Change
`

var CreateMergeConflictFileMultiple = func(shell *Shell) {
	shell.
		NewBranch("original-branch").
		EmptyCommit("one").
		EmptyCommit("two").
		EmptyCommit("three").
		CreateFileAndAdd("file", OriginalFileContentMultiple).
		Commit("original").
		NewBranch("first-change-branch").
		UpdateFileAndAdd("file", FirstChangeFileContentMultiple).
		Commit("first change").
		Checkout("original-branch").
		NewBranch("second-change-branch").
		UpdateFileAndAdd("file", SecondChangeFileContentMultiple).
		Commit("second change").
		EmptyCommit("second-change-branch unrelated change").
		Checkout("first-change-branch")

	shell.RunShellCommandExpectError("git merge --no-edit second-change-branch")
}
