package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchScrollsToAMatchOffTheEdge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With wrapping turned off for diffs, searching the focused diff scrolls it sideways to show a match past its edge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one "+strings.Repeat("word ", 60)+"needle\ntwo needle\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			FilterOrSearch("needle").
			SelectedLines(
				Contains("+one word"),
			).
			// The match is at column 305 of the diff; it ends up a third of the way
			// across the pane, which is 78 columns wide.
			OriginX(279).
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("+two needle"),
			).
			OriginX(0)
	},
})
