package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoSelectionOverABinaryDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff with nothing selectable in it shows no selection, and a refresh doesn't bring one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("binary", "\x00one\x00two\x00")
		shell.Commit("one")

		shell.UpdateFile("binary", "\x00one\x00TWO\x00")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// git says only that the file differs, so there is nothing to select — and a
		// refresh, which renders the same diff again, doesn't make one appear.
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			Tap(func() {
				t.GlobalPress(keys.Universal.Refresh)
			}).
			SelectionIsHidden()
	},
})
