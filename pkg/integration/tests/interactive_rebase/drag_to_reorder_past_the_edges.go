package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragToReorderPastTheEdges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Dragging the first commit above the view, or the last one below the list, doesn't move it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit-03").IsSelected(),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			ClickAndHold(1, 0).
			MouseMove(1, -1).
			Lines(
				Contains("commit-03").IsSelected(),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			MouseRelease().
			Lines(
				Contains("commit-03").IsSelected(),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			NavigateToLine(Contains("commit-01")).
			ClickAndHold(1, 2).
			MouseMove(1, 3).
			Lines(
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01").IsSelected(),
			).
			MouseRelease().
			Lines(
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01").IsSelected(),
			)
	},
})
