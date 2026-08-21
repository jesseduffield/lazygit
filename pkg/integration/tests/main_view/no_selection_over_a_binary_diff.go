package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoSelectionOverABinaryDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff with nothing selectable in it shows no selection, however it came to be showing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("text", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("binary", "\x00one\x00two\x00")
		shell.Commit("one")

		shell.UpdateFile("text", "one\nTWO\nthree\n")
		shell.UpdateFile("binary", "\x00one\x00TWO\x00")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// git says only that the file differs, so there is nothing to select — and a
		// refresh, which renders the same diff again, doesn't make one appear.
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("binary")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			Tap(func() {
				t.GlobalPress(keys.Universal.Refresh)
			}).
			SelectionIsHidden()

		// The same when acting on a diff of several files leaves nothing selectable in
		// it: staging the text file's only change leaves the binary one behind.
		t.Views().Files().
			Focus().
			NavigateToLine(Contains("▼ /")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Main.ToggleSelectHunk).
			PressPrimaryAction().
			SelectionIsHidden()
	},
})
