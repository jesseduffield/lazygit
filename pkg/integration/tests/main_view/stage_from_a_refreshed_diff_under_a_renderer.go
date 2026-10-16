package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var StageFromARefreshedDiffUnderARenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a line of a renderer's diff works once the diff has been rendered again after a change, and after the index changed for another file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{shared.RecordsDiffRenderer}
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

		// The file changes outside of lazygit, and the diff is rendered again. What it was
		// made from goes with the diff now on screen.
		t.Shell().UpdateFile("file1", "one\nADD1\ntwo\nthree\nEND\n")
		t.GlobalPress(keys.Universal.Refresh)
		t.Views().Main().
			IsFocused().
			ContainsLines(
				Contains("+END"),
			).
			SelectedLines(
				Contains("+ADD1"),
			)

		// The index is rewritten for the other file, which leaves this file's diff alone.
		t.Shell().GitAdd("file2")

		t.Views().Main().
			PressPrimaryAction()

		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Contains("MM file1"),
				Contains("M  file2"),
			)
	},
})
