package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragToReorderWithAutoscroll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep scrolling commits while a dragged commit is held at the panel edge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(40)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			TopLines(
				Contains("commit-40").IsSelected(),
			).
			// Click and hold the first commit
			ClickAndHold(1, 0).
			// Move the mouse to the bottom of the panel to trigger autoscroll
			MouseMoveToBottom(1).
			// Verify that the view scrolls
			OriginYAtLeast(3).
			// Move the mouse back into the viewport
			MouseMove(1, 1).
			// This keeps the scroll as it was
			OriginYAtLeast(3).
			MouseRelease().
			SelectedLines(
				Contains("commit-40"),
			).
			SelectedLineIdxAtLeast(3).
			// Scroll back to verify that the original commit is no longer at the top
			GotoTop().
			TopLines(
				Contains("commit-39").IsSelected(),
			)
	},
})
