package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragSelectsDiffLineRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Dragging in the main view's diff selects the range from the line the drag started on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nFOUR\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Hunk mode is on, so the mouse-down alone selects the whole block; the drag
		// anchors the range where the mouse went down instead, one line at a time.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("-four"),
				Contains("+THREE"),
				Contains("+FOUR"),
			).
			ClickAndHold(0, 8).
			MouseMove(0, 9).
			SelectedLines(
				Contains("-four"),
				Contains("+THREE"),
			).
			MouseMove(0, 10).
			SelectedLines(
				Contains("-four"),
				Contains("+THREE"),
				Contains("+FOUR"),
			).
			MouseRelease().
			SelectedLines(
				Contains("-four"),
				Contains("+THREE"),
				Contains("+FOUR"),
			)
	},
})
