package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoRedoCherryPick = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo and redo a cherry-pick, which must go back to the state before the pick and then restore the picked commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("f", "a")
		shell.Commit("A base")

		shell.NewBranch("feature")
		shell.UpdateFileAndAdd("f", "a\nb")
		shell.Commit("B feature")

		shell.Checkout("master")
		shell.CreateFileAndAdd("g", "c")
		shell.Commit("C main")

		shell.RunCommand([]string{"git", "cherry-pick", "feature"})
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
				Contains("B feature").IsSelected(),
				Contains("C main"),
				Contains("A base"),
			)

		// Undo goes back to the state before the cherry-pick, which means the picked
		// commit is gone but C main is still there.
		t.Views().Commits().
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			Lines(
				Contains("C main").IsSelected(),
				Contains("A base"),
			)

		// Redo brings the cherry-picked commit back.
		t.Views().Commits().
			Press(keys.Universal.Redo).
			Tap(confirmRedo).
			Lines(
				Contains("B feature").IsSelected(),
				Contains("C main"),
				Contains("A base"),
			)
	},
})
