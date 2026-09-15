package shared

import (
	"fmt"

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

// A conflict-marker-size that isn't git's default of 7. It's set for file types
// whose regular content tends to contain marker-looking lines, e.g.
// documentation about merging, or test scripts.
const CustomConflictMarkerSize = 32

// Makes git write conflict markers of CustomConflictMarkerSize characters into
// the file that the setups below create conflicts in. Call this before one of
// them.
var SetCustomConflictMarkerSize = func(shell *Shell) {
	shell.CreateFileAndAdd(".gitattributes",
		fmt.Sprintf("file conflict-marker-size=%d\n", CustomConflictMarkerSize)).
		Commit("set a custom conflict marker size")
}

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

	shell.RunCommandExpectError([]string{"git", "merge", "--no-edit", "second-change-branch"})
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

	shell.RunCommandExpectError([]string{"git", "merge", "--no-edit", "second-change-branch"})
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

	shell.RunCommandExpectError([]string{"git", "merge", "--no-edit", "second-change-branch"})
}

var CreateMergeConflictFileForMergeFileTests = func(shell *Shell, originalFileContent string, currentChangeFileContent string, incomingChangeFileContent string) {
	shell.
		NewBranch("original-branch").
		EmptyCommit("one").
		CreateFileAndAdd("file", originalFileContent).
		Commit("original").
		NewBranch("current-change-branch").
		UpdateFileAndAdd("file", currentChangeFileContent).
		Commit("first change").
		Checkout("original-branch").
		NewBranch("incoming-change-branch").
		UpdateFileAndAdd("file", incomingChangeFileContent).
		Commit("second change").
		Checkout("current-change-branch").
		RunCommandExpectError([]string{"git", "merge", "--no-edit", "incoming-change-branch"})
}
