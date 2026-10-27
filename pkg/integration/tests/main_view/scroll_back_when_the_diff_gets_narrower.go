package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollBackWhenTheDiffGetsNarrower = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff scrolled sideways scrolls back when it is rendered again narrower than where it was scrolled to",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "ONE\ntwo\nthree\nfour\nfive "+strings.Repeat("word ", 60)+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-five"),
			).
			Press(keys.Universal.RangeSelectDown).
			SelectedLines(
				Contains("-five"),
				Contains("+five word"),
			).
			Press(keys.Universal.ScrollRight).
			Press(keys.Universal.ScrollRight).
			OriginX(52).
			// What is left of the diff fits in the pane, so the pane goes back to the left
			// edge.
			PressPrimaryAction().
			SelectedLines(
				Contains("+ONE"),
			).
			OriginX(0)
	},
})
