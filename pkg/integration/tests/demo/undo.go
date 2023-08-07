package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// TODO: fix confirmation view wrapping issue: https://github.com/jesseduffield/lazygit/issues/2872

var Undo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.NerdFontsVersion = "3"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Undo commands")
		t.Wait(1000)

		confirmCommitDrop := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Delete commit")).
				Content(Equals("Are you sure you want to delete this commit?")).
				Wait(500).
				Confirm()
		}

		confirmUndo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Undo")).
				Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Wait(500).
				Confirm()
		}

		confirmRedo := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Redo")).
				Content(MatchesRegexp(`Are you sure you want to hard reset to '.*'\? An auto-stash will be performed if necessary\.`)).
				Wait(500).
				Confirm()
		}

		t.Views().Commits().Focus().
			SetCaptionPrefix("Drop two commits").
			Wait(1000).
			Press(keys.Universal.Remove).
			Tap(confirmCommitDrop).
			Press(keys.Universal.Remove).
			Tap(confirmCommitDrop).
			SetCaptionPrefix("Undo the drops").
			Wait(1000).
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			Press(keys.Universal.Undo).
			Tap(confirmUndo).
			SetCaptionPrefix("Redo the drops").
			Wait(1000).
			Press(keys.Universal.Redo).
			Tap(confirmRedo).
			Press(keys.Universal.Redo).
			Tap(confirmRedo)
	},
})
