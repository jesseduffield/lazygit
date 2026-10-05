package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveOnWhenThePatchMarksRewrapTheDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Taking the first line into a custom patch moves the selection on to the next change, although the marks that come with the patch narrow the diff and wrap a line above the selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file", "one\nTWO\nthree\nFOUR\nfive\n")
		// git draws the diffstat line of a file with this many changes as wide as the
		// view, so the columns the marks take make that line wrap.
		shell.CreateFileAndAdd("many", strings.Repeat("line\n", 300))
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+TWO")).
			PressPrimaryAction().
			MarkedLines(
				Contains("+TWO"),
			).
			SelectedLines(
				Contains("-four"),
			)
	},
})
