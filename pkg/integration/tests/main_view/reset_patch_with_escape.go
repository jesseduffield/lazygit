package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetPatchWithEscape = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset a custom patch after escaping from the focused commit diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+file1 content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		// Leave the focused diff and then the commit files panel. Escape at the top
		// level gives up the patch.
		t.Views().Main().PressEscape()
		t.Views().CommitFiles().IsFocused().PressEscape()
		t.Views().Commits().IsFocused().PressEscape()

		t.Views().Information().Content(DoesNotContain("Building patch"))
	},
})
