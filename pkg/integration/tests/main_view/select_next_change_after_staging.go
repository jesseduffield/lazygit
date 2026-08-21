package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextChangeAfterStaging = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After staging from the focused main view the selection lands on the change that took the place of the one staged",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADD1\nADD2\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Line by line, each press leaves the selection on the next change to stage.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("+ADD2"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("-nine"),
			).
			// A hunk goes the same way: the block after the one staged is what is left
			// in its place, and here that is nothing, so the last change stays selected.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			).
			PressPrimaryAction().
			Tap(func() {
				t.Views().Files().Lines(
					Contains("M  file1"),
				)
			})
	},
})
