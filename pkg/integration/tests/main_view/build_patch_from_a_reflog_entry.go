package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildPatchFromAReflogEntry = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Take lines of the diff of a reflog entry into a custom patch, which can then be applied",
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
		t.Views().ReflogCommits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().ContainsLines(
			Contains("-two"),
			Contains(" three"),
		)

		// A reflog entry is never a commit we may rewrite, so the patch can be applied
		// but not moved out of the commit it came from.
		t.Common().SelectPatchOption(Contains("Apply patch in reverse"))

		t.Views().Files().
			Focus().
			Lines(
				Contains("M").Contains("file1"),
			)
	},
})
