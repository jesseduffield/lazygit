package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesWithMnemonicPrefixes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a line from the focused main view when diff.mnemonicPrefix is set",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADD1\nADD2\ntwo\nthree\n")

		// git then puts i/ and w/ in front of the paths of the working tree's diff
		// instead of a/ and b/.
		shell.SetConfig("diff.mnemonicPrefix", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// The user chose these prefixes, so the diff shows them.
		t.Views().Main().
			IsFocused().
			ContainsLines(
				Equals("diff --git i/file1 w/file1"),
			).
			SelectedLines(
				Contains("+ADD1"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			/* EXPECTED:
			Contains("MM file1"),
			ACTUAL: */
			Contains(" M file1"),
		)
	},
})
