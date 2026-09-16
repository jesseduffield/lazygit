package ui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragBeyondViewport = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Dragging a range selection beyond the bottom of the panel doesn't scroll the view",
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
			OriginY(0).
			// The pointer ends up below the panel, so the range extends to a line
			// that isn't visible. Scrolling there is the drag autoscroller's job,
			// which scrolls line by line for as long as the pointer stays there;
			// the drag itself must leave the scroll position alone.
			ClickAndHold(1, 1).
			MouseMove(1, 8).
			MouseRelease().
			SelectedLineIdx(8).
			OriginY(0)
	},
})
