package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var UnstageFromADiffWhoseIndexChangedUnderARenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging a line of a renderer's staged diff after the index changed since it was shown unstages nothing and shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{shared.RecordsDiffRenderer}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\nADD1\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			)

		// Another line takes the place of the selected one in the index outside of
		// lazygit, and nothing has refreshed the diff since.
		t.Shell().UpdateFileAndAdd("file1", "one\nOTHER\ntwo\nthree\n")

		t.Views().Secondary().
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("The diff of 'file1' has changed since it was shown. Check the updated diff and select the lines again.")).
			Confirm()

		t.Views().Files().
			Lines(
				Contains("M  file1"),
			)
	},
})
