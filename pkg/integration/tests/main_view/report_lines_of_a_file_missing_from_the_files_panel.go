package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ReportLinesOfAFileMissingFromTheFilesPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging lines that the diff puts in a file the files panel doesn't have shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// A renderer that announces the metadata protocol and passes the diff through,
		// except that its +++ line names a file in another directory.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007'; sed 's#+++ b/#+++ b/elsewhere/#'`},
		}
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
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("Can't find the file 'elsewhere/file1' that the selected lines belong to")).
			Confirm()

		t.Views().Files().Lines(
			Contains(" M file1"),
		)
	},
})
