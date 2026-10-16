package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageAfterChangesThatLeaveTheDiffAlone = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a line still works after changes outside of lazygit that leave the file's diff as it was",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("file2", "a\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADD1\ntwo\nthree\n")
		shell.UpdateFile("file2", "b\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains(" M file1"),
				Contains(" M file2"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			)

		// Outside of lazygit, and with nothing refreshing the diff: the file is written
		// with the content it has, the index is rewritten for another file, and HEAD
		// moves to a commit of that other file.
		t.Shell().
			UpdateFile("file1", "one\nADD1\ntwo\nthree\n").
			GitAdd("file2").
			Commit("two")

		t.Views().Main().
			PressPrimaryAction()

		t.Views().Files().
			Lines(
				Contains("M  file1"),
			)
		t.Views().Secondary().
			ContainsLines(
				Contains("+ADD1"),
			)
	},
})
