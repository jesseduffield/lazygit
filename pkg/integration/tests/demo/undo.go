package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Undo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Undo",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Undo commands")
		t.Wait(1000)

		confirmCommitDrop := func() {
			t.ExpectPopup().Confirmation().
				Title(Equals("Drop commit")).
				Content(Equals("Are you sure you want to drop the selected commit(s)?")).
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
			Tap(confirmUndo)
	},
})
