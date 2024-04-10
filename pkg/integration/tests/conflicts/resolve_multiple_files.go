package conflicts

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
	"github.com/lobes/lazytask/pkg/integration/tests/shared"
)

var ResolveMultipleFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensures that upon resolving conflicts for one file, the next file is selected",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFiles(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file1").IsSelected(),
				Contains("UU").Contains("file2"),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			PressPrimaryAction()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU").Contains("file2").IsSelected(),
			).
			PressEnter()

		// coincidentally these files have the same conflict
		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			PressPrimaryAction()

		t.Common().ContinueOnConflictsResolved()
	},
})
