package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OriginalFileContent = `
This Is The
Original
File
..
`

var FirstChangeFileContent = `
This Is The
First Change
File
..
This Change Is Only In The First File
`

var SecondChangeFileContent = `
This Is The
Second Change
File
..
`

var MergeSecondFileFinalContent = `
This Is The
Second Change
File
..
This Change Is Only In The First File
`

var MergeFileTheirs = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Conflicting file can be resolved to 'their' version via merge-file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("original-branch").
			EmptyCommit("one").
			CreateFileAndAdd("file", OriginalFileContent).
			Commit("original").
			NewBranch("first-change-branch").
			UpdateFileAndAdd("file", FirstChangeFileContent).
			Commit("first change").
			Checkout("original-branch").
			NewBranch("second-change-branch").
			UpdateFileAndAdd("file", SecondChangeFileContent).
			Commit("second change").
			Checkout("first-change-branch").
			RunCommandExpectError([]string{"git", "merge", "--no-edit", "second-change-branch"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file").IsSelected(),
			)

		t.GlobalPress(keys.Files.OpenMergeOptions)

		t.ExpectPopup().Menu().
			Title(Equals("Resolve merge conflicts")).
			Select(Contains("Use Incoming")). // merge-file --theirs
			Confirm()

		t.Common().ContinueOnConflictsResolved("merge")

		t.Views().Files().IsEmpty()

		t.FileSystem().FileContent("file", Equals(MergeSecondFileFinalContent))
	},
})
