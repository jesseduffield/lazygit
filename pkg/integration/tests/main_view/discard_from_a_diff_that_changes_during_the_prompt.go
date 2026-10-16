package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardFromADiffThatChangesDuringThePrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a line of a diff that changes while the confirmation is open discards nothing and shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
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
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard change")).
			Content(Contains("Are you sure you want to discard this change"))

		// Another line takes the place of the selected one outside of lazygit while the
		// user is being asked.
		t.Shell().UpdateFile("file1", "one\nOTHER\ntwo\nthree\n")

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard change")).
			Content(Contains("Are you sure you want to discard this change")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("The diff of 'file1' has changed since it was shown. Check the updated diff and select the lines again.")).
			Confirm()

		t.FileSystem().FileContent("file1",
			Equals("one\nOTHER\ntwo\nthree\n"))
	},
})
