package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollDiffHorizontallyStopsAtTheEnd = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Scrolling a diff sideways with H and L stops when the end of the widest line reaches the right edge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")

		shell.UpdateFile("file1", strings.Repeat("x", 83)+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			OriginX(0).
			// A third of the pane's width would be 26 columns, but the end of the line is
			// at the right edge already after 6.
			Press(keys.Universal.ScrollRight).
			OriginX(6).
			Press(keys.Universal.ScrollRight).
			OriginX(6).
			Press(keys.Universal.ScrollLeft).
			OriginX(0)
	},
})
