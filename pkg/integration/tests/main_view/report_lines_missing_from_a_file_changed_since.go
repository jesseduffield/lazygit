package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ReportLinesMissingFromAFileChangedSince = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging lines that are gone from the file since its diff was shown shows an error",
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
			)

		// The change is undone outside of lazygit, and nothing has refreshed the diff
		// since.
		t.Shell().UpdateFile("file1", "one\ntwo\nthree\n")

		t.Views().Main().
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("Can't find the selected lines in the diff of 'file1'")).
			Confirm()
	},
})
