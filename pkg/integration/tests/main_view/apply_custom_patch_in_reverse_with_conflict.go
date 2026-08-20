package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyCustomPatchInReverseWithConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a multi-file custom patch in reverse when one file conflicts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "file1 content\nmore file1 content\n")
		shell.UpdateFileAndAdd("file2", "file2 content\nmore file2 content\n")
		shell.Commit("second commit")
		shell.UpdateFileAndAdd("file1", "file1 content\nmore file1 content\neven more file1\n")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("second commit"),
				Contains("first commit"),
			).
			NavigateToLine(Contains("second commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M file1"),
				Equals("  M file2"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+more file1 content")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+more file2 content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().Content(Contains("+more file1 content").Contains("+more file2 content"))
		t.Common().SelectPatchOption(Contains("Apply patch in reverse"))

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Applied patch to 'file1' with conflicts.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Equals("UU file1").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			ContainsLines(
				Contains("file1 content"),
				Contains("<<<<<<< ours").IsSelected(),
				Contains("more file1 content").IsSelected(),
				Contains("even more file1").IsSelected(),
				Contains("=======").IsSelected(),
				Contains(">>>>>>> theirs"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /"),
				Equals("  M  file1").IsSelected(),
				Equals("  M  file2"),
			)
		t.Views().Secondary().ContainsLines(
			Contains(" file1 content"),
			Contains("-more file1 content"),
			Contains("-even more file1"),
		)
	},
})
