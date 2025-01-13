package undo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UndoCommitSoft = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo/redo a commit using soft reset",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Undo.CommitReset = "soft"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		confirmUndo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Undo")).
				Content(MatchesRegexp(`Are you sure you want to soft reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Confirm()
		}

		confirmRedo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Redo")).
				Content(MatchesRegexp(`Are you sure you want to soft reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Confirm()
		}

		confirmAutostash := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Autostash?")).
				Content(MatchesRegexp(`You must stash and pop your changes to bring them across\. Do this automatically\? \(enter\/esc\)`)).
				Confirm()
		}

		t.Views().Commits().Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			Lines(
				Contains("one").IsSelected(),
			)

		t.Views().Files().
			Lines(
				Contains("A  file"),
			)

		t.Views().Commits().Focus().
			Press(keys.Universal.Redo).
			Tap(confirmRedo).
			Tap(confirmAutostash).
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			)

		t.Views().Files().
			IsEmpty()
	},
})
