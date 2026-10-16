package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditHunkOfAStaleDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Editing the hunk of a diff that changed since it was shown opens no editor and shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Stand in for the editor: keep a copy of the patch it was given, outside the
		// repo so that the files panel keeps saying what the test is about, and hand it
		// back as it is, so that it is staged.
		cfg.GetUserConfig().OS.EditAtLineAndWait = "cp {{filename}} ../edited-patch"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADD1\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(" M file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			)

		// Another line takes the place of the selected one outside of lazygit, and
		// nothing has refreshed the diff since.
		t.Shell().UpdateFile("file1", "one\nOTHER\ntwo\nthree\n")

		t.Views().Main().
			Press(keys.Main.EditSelectHunk)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("The diff of 'file1' has changed since it was shown. Check the updated diff and select the lines again.")).
			Confirm()

		t.Views().Files().
			Lines(
				Contains(" M file1"),
			)
	},
})
