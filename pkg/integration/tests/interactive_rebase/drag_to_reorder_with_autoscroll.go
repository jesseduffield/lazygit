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
			ClickAndHold(1, 0).
			MouseMoveToBottom(1).
			OriginYAtLeast(3).
			MouseRelease().
			SelectedLines(
				Contains("commit-40"),
			).
			SelectedLineIdxAtLeast(3).
			GotoTop().
			TopLines(
				Contains("commit-39").IsSelected(),
			)
	},
})
