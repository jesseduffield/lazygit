package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EnterAndDoubleClickFocusFileDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter and double-click focus a file's diff in the main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "base\n")
		shell.Commit("add file")
		shell.UpdateFile("file1", "changed\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-base")).
			PressEscape()

		t.Views().Files().
			IsFocused().
			Click(3, 0).
			Click(3, 0)
		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-base")).
			PressEscape()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("add file").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+base")).
			PressEscape()

		t.Views().CommitFiles().
			IsFocused().
			Click(3, 0).
			Click(3, 0)
		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+base"))
	},
})
