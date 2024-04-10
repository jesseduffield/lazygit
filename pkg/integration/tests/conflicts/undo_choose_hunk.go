package conflicts

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
	"github.com/lobes/lazytask/pkg/integration/tests/shared"
)

var UndoChooseHunk = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Chooses a hunk when resolving a merge conflict and then undoes the choice",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFileMultiple(shell)
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
			Content(Contains("<<<<<<< HEAD\nFirst Change")).
			// explicitly asserting on the selection because sometimes the content renders
			// before the selection is ready for user input
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			PressPrimaryAction().
			// choosing the first hunk
			Content(DoesNotContain("<<<<<<< HEAD\nFirst Change")).
			Press(keys.Universal.Undo).
			Content(Contains("<<<<<<< HEAD\nFirst Change"))
	},
})
