package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollDiffHorizontallyWithTheWheel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With wrapping turned off for diffs, the horizontal mouse wheel scrolls the diff sideways even when it isn't focused, and stops when the end of the widest line reaches the right edge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("one")

		shell.UpdateFile("file1", strings.Repeat("x", 83)+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused()

		t.Views().Main().
			OriginX(0).
			// Twice gui.scrollHeight columns at a time.
			ScrollWheelRight().
			OriginX(4).
			// The end of the line is at the right edge now.
			ScrollWheelRight().
			OriginX(6).
			ScrollWheelRight().
			OriginX(6).
			ScrollWheelLeft().
			OriginX(2).
			ScrollWheelLeft().
			OriginX(0)
	},
})
