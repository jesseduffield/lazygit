package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var JumpToAnotherFileScrollsToTheLeftEdge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jumping to another file of a diff scrolled sideways shows it from the left edge, while jumping to another hunk doesn't",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		long := strings.Repeat("word ", 60)
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.CreateFileAndAdd("file2", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\ntwo\nthree "+long+"\nfour\nfive\nsix\nseven\neight\nnine "+long+"\nten\n")
		shell.UpdateFileAndAdd("file2", "one\ntwo\nthree "+long+"\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Universal.ScrollRight).
			OriginX(26).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-nine"),
			).
			OriginX(26).
			Press(keys.Main.NextFile).
			SelectedLines(
				Contains("diff --git a/file2 b/file2"),
			).
			OriginX(0).
			Press(keys.Universal.ScrollRight).
			OriginX(26).
			Press(keys.Main.PrevFile).
			SelectedLines(
				Contains("diff --git a/file1 b/file1"),
			).
			OriginX(0)
	},
})
