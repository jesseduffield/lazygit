package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoRedoMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo and redo a merge, which must go back to the state before the merge and then restore the merge commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("f", "a")
		shell.Commit("A base")

		shell.NewBranch("feature")
		shell.CreateFileAndAdd("g", "b")
		shell.Commit("B feature")

		shell.Checkout("master")
		shell.CreateFileAndAdd("h", "c")
		shell.Commit("C main")

		shell.Merge("feature")
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
				Contains("Merge branch 'feature'").IsSelected(),
				Contains("B feature"),
				Contains("C main"),
				Contains("A base"),
			)

		t.Views().Commits().
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			Lines(
				Contains("C main").IsSelected(),
				Contains("A base"),
			)

		t.Views().Commits().
			Press(keys.Universal.Redo).
			Tap(confirmRedo).
			Lines(
				Contains("Merge branch 'feature'").IsSelected(),
				Contains("B feature"),
				Contains("C main"),
				Contains("A base"),
			)
	},
})
