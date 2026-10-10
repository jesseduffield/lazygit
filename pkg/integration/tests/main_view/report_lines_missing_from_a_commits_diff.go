package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ReportLinesMissingFromACommitsDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Adding lines to a custom patch that the diff puts where the commit's diff has none shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// A renderer that announces the metadata protocol and passes the diff through,
		// except that its hunk header puts the hunk ten lines further down.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007'; sed 's#@@ -1,3 +1,4 @@#@@ -11,3 +11,4 @@#'`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\nADD1\ntwo\nthree\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
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
			Content(Equals("Can't find the selected lines in the diff of 'file1'")).
			Confirm()

		// No patch is left behind.
		t.Views().Information().Content(DoesNotContain("Building patch"))
	},
})
