package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UnstageFromADiffWhoseHeadMoved = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging a line of the staged diff after HEAD moved since it was shown unstages nothing and shows an error",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("two")

		shell.UpdateFileAndAdd("file1", "one\nthree\n")
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
				Contains("-two"),
			)

		// HEAD moves back to the commit with the other line outside of lazygit, without
		// touching the index, and nothing has refreshed the diff since. The staged diff
		// now deletes the other line, at the same place.
		t.Shell().RunCommand([]string{"git", "reset", "--soft", "HEAD~"})

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
