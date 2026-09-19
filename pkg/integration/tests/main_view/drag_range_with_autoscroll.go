package main_view

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragRangeWithAutoscroll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep scrolling while dragging a range selection at the edge of the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		fileContent := "base\n"
		shell.CreateFileAndAdd("file1", fileContent)
		shell.Commit("one")
		for i := 1; i <= 40; i++ {
			fileContent += fmt.Sprintf("line %d\n", i)
		}
		shell.UpdateFile("file1", fileContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The diff is longer than the view, so holding the pointer at the bottom edge has
		// to keep scrolling and extending the selection rather than stopping there.
		t.Views().Main().
			IsFocused().
			ClickAndHold(1, 6).
			MouseMoveToBottom(1).
			OriginYAtLeast(3).
			SelectedLineIdxAtLeast(9).
			MouseRelease()
	},
})
