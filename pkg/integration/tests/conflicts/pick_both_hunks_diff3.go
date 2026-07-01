package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var PickBothHunksDiff3 = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pick both hunks of a conflict rendered in the diff3 style; the common ancestor must not be included",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("merge.conflictStyle", "diff3")
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			// the diff3 style renders the common ancestor between the two changes
			Content(Contains("<<<<<<< HEAD\nFirst Change")).
			Content(Contains("||||||| ")).
			Content(Contains("Original")).
			Press(keys.Main.PickBothHunks)

		t.Common().ContinueOnConflictsResolved("merge")

		t.Views().Files().IsEmpty()

		t.FileSystem().FileContent("file",
			/* EXPECTED:
			Equals("\nThis\nIs\nThe\nFirst Change\nSecond Change\nFile\n")
			ACTUAL: */
			Equals("\nThis\nIs\nThe\nFirst Change\nOriginal\nSecond Change\nFile\n"))
	},
})
