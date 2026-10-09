package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollDiffHorizontally = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With wrapping turned off for diffs, the focused diff scrolls sideways, and stays scrolled as the selection moves and the diff is rendered again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		long := strings.Repeat("word ", 60)
		shell.UpdateFile("file1", "one "+long+"\ntwo "+long+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-one"),
			).
			OriginX(0).
			// A third of the pane's width at a time.
			Press(keys.Universal.ScrollRight).
			OriginX(26).
			Press(keys.Universal.ScrollLeft).
			OriginX(0).
			Press(keys.Universal.ScrollLeft).
			OriginX(0).
			Press(keys.Universal.ScrollRight).
			Press(keys.Universal.ScrollRight).
			OriginX(52).
			SelectNextItem().
			SelectedLines(
				Contains("-two"),
			).
			OriginX(52).
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			SelectedLines(
				Contains("-two"),
			).
			OriginX(52).
			PressPrimaryAction().
			SelectedLines(
				Contains("+one"),
			).
			OriginX(52).
			// Leaving the view gives up the horizontal scroll position.
			PressEscape()

		t.Views().Files().
			IsFocused()

		t.Views().Main().
			OriginX(0)
	},
})
