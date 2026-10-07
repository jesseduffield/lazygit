package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesWithNoPrefixes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a line of a file in a directory named b from the focused main view when diff.noprefix is set",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("b/file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("b/file1", "one\nADD1\nADD2\ntwo\nthree\n")

		// git then leaves out the a/ and b/ in front of the paths of a diff. The b/
		// that the paths of this file start with is its directory.
		shell.SetConfig("diff.noprefix", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ b").IsSelected(),
				Contains(" M file1"),
			).
			SelectNextItem().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			Contains("▼ b"),
			/* EXPECTED:
			Contains("MM file1"),
			ACTUAL: */
			Contains(" M file1"),
		)
	},
})
