package ui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BackgroundRefreshKeepsScrollPosition = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A background refresh doesn't scroll the selection back into view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		for i := range 20 {
			shell.CreateFile(fmt.Sprintf("file%02d", i), "")
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			SelectNextItem().
			SelectedLine(Contains("file00")).
			// Scroll the selection out of view with the mouse wheel
			ScrollWheelDown().
			ScrollWheelDown().
			OriginY(4).
			Tap(func() {
				t.Shell().CreateFile("aaa", "")
				t.RefreshInBackground()
			}).
			// The new file sorts before the selected one, so the selection has
			// moved down a line; the view must stay where the user left it though
			SelectedLineIdx(2).
			OriginY(4)
	},
})
