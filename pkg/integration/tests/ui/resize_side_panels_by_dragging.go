package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResizeSidePanelsByDragging = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resize the side panels by dragging the boundary between them and the main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			HasWidth(40).
			// The right edge of the files panel is part of the boundary
			Hover(38, 1)

		t.Views().Options().
			Content(Equals("Drag to resize | Double-click to reset"))

		// The left edge of the main view is part of the boundary, too
		t.Views().Main().
			ClickAndHold(-1, 3).
			MouseMove(-21, 3).
			MouseRelease()

		t.Views().Files().
			IsFocused().
			HasWidth(20)
		t.Views().Main().
			HasWidth(100)

		t.Views().Main().
			Hover(10, 3)

		t.Views().Options().
			Content(DoesNotContain("Drag to resize"))

		// A click on the boundary doesn't focus the main view
		t.Views().Main().
			Click(-1, 3)

		t.Views().Files().
			IsFocused().
			HasWidth(20)

		t.Views().Main().
			DoubleClick(-1, 3)

		t.Views().Files().
			HasWidth(40)
	},
})
