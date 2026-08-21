package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDeletedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Acting on the whole of a file's block in a directory's diff acts on the file: staging a deletion, and unstaging an addition",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("fileA", "a\n")
		shell.CreateFileAndAdd("fileB", "b1\nb2\n")
		shell.Commit("one")

		shell.UpdateFile("fileA", "a\nfromA\n")
		shell.DeleteFile("fileB")
		shell.CreateFileAndAdd("fileC", "c1\nc2\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ /").IsSelected(),
				Contains(" M fileA"),
				Contains(" D fileB"),
				Contains("A  fileC"),
			).
			Press(keys.Universal.FocusMainView)

		// Select everything the deleted file contributes to the diff.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("-b1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("-b2")).
			SelectedLines(
				Contains("-b1"),
				Contains("-b2"),
			).
			PressPrimaryAction()

		// The file is staged as deleted. Applying its lines as a patch would have left
		// an empty file in the index instead, which is not what deleting a file means.
		t.Views().Files().Lines(
			Contains("▼ /"),
			Contains(" M fileA"),
			Contains("D  fileB"),
			Contains("A  fileC"),
		)

		// The same the other way round: taking the whole of an added file back out of
		// the index leaves it untracked, rather than tracked and empty.
		t.Views().Main().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			NavigateToLine(Contains("+c1")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+c2")).
			SelectedLines(
				Contains("+c1"),
				Contains("+c2"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			Contains("▼ /"),
			Contains(" M fileA"),
			Contains("D  fileB"),
			Contains("?? fileC"),
		)
	},
})
