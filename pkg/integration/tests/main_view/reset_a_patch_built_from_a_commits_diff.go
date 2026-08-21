package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetAPatchBuiltFromACommitsDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset a custom patch built from a commit's diff without ever entering the commit's files, and stay in the diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Build the patch straight from the commit's diff, so that nothing has ever told
		// the commit files panel which commit it would be showing.
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Reset patch"))

		// Giving up the patch leaves the diff it was being built from, and the focus in it.
		t.Views().Information().Content(DoesNotContain("Building patch"))
		t.Views().Main().
			IsFocused().
			Content(Contains("-two"))
	},
})
