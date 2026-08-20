package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartPatchFromAnotherCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirm replacing a custom patch when selecting lines from another commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("first commit")
		shell.CreateFileAndAdd("file2", "file2 content")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file2").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+file2 content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().Content(Contains("file2"))

		t.Views().Main().PressEscape()
		t.Views().CommitFiles().IsFocused().PressEscape()
		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("first commit")).
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

		t.ExpectPopup().Confirmation().
			Title(Contains("Discard patch")).
			Content(Contains("You can only build a patch from one commit/stash-entry at a time. Discard current patch?")).
			Confirm()

		t.Views().Secondary().Content(Contains("file1").DoesNotContain("file2"))
	},
})
