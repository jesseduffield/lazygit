package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NavigateByHunkAndFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jump from hunk to hunk and from file to file in the focused main view of a commit's diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.CreateFileAndAdd("file2", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
		shell.UpdateFileAndAdd("file2", "one\ntwo\nTHREE\n")
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
			// Hunk navigation moves between change blocks, which is what lazygit means by
			// a hunk: a run of changes bounded by context, of which one @@ hunk may hold
			// several.
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-nine"),
			).
			Press(keys.Main.PrevHunk).
			SelectedLines(
				Contains("-three"),
			).
			// File navigation lands on the top of the next file's diff, which for a
			// parseable diff is its header.
			Press(keys.Main.NextFile).
			SelectedLines(
				Contains("diff --git a/file2 b/file2"),
			).
			Press(keys.Main.NextFile).
			SelectedLines(
				Contains("diff --git a/file2 b/file2"),
			).
			Press(keys.Main.PrevFile).
			SelectedLines(
				Contains("diff --git a/file1 b/file1"),
			).
			// A range that only grows while shift is held is a plain selection again once
			// we jump elsewhere, rather than stretching to wherever we land.
			Press(keys.Main.NextHunk).
			Press(keys.Universal.RangeSelectDown).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-nine"),
			).
			// A sticky range does stretch to it.
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Main.PrevHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
				Contains(" four"),
				Contains(" five"),
				Contains(" six"),
				Contains(" seven"),
				Contains(" eight"),
				Contains("-nine"),
			)
	},
})
