package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CopyFromAStaleDiffUnderARenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copying lines of a renderer's diff whose file changed since it was shown copies nothing and shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{shared.RecordsDiffRenderer}
		// Emulate the clipboard by writing to a file.
		cfg.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
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
			Press(keys.Universal.CopyToClipboard)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("The diff of 'file1' has changed since it was shown. Check the updated diff and select the lines again.")).
			Confirm()
	},
})
