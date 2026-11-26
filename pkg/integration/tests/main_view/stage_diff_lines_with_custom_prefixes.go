package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesWithCustomPrefixes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a line from the focused main view when diff.srcPrefix and diff.dstPrefix are set",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.45.0"),
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADD1\nADD2\ntwo\nthree\n")

		// git then puts these strings in front of the paths of a diff instead of
		// a/ and b/. Without a slash, nothing tells where they end.
		shell.SetConfig("diff.srcPrefix", "old:")
		shell.SetConfig("diff.dstPrefix", "new:")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			Contains("MM file1"),
		)
	},
})
