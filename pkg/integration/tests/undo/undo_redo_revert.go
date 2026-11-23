package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoRedoRevert = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo and redo a revert, which must go back to the state before the revert and then restore the revert commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("f", "a")
		shell.Commit("A base")

		shell.UpdateFileAndAdd("f", "b")
		shell.Commit("B change")

		shell.Revert("HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		confirmUndo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Undo")).
				Content(MatchesRegexp(`Are you sure you want to soft reset to '.*'\?`)).
				Confirm()
		}

		confirmRedo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Redo")).
				Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Confirm()
		}

		t.Views().Commits().
			Focus().
			Lines(
				Contains(`Revert "B change"`).IsSelected(),
				Contains("B change"),
				Contains("A base"),
			)

		// The reverted commit has to survive the undo; only the revert is removed.
		t.Views().Commits().
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			Lines(
				Contains("B change").IsSelected(),
				Contains("A base"),
			)

		t.Views().Commits().
			Press(keys.Universal.Redo).
			Tap(confirmRedo).
			Lines(
				Contains(`Revert "B change"`).IsSelected(),
				Contains("B change"),
				Contains("A base"),
			)
	},
})
