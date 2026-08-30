package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragToReorderInRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drag rebase todos without allowing real commits to move",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(5)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit-01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("commit-05"),
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("─── Commits"),
				Contains("commit-01").IsSelected(),
			).
			NavigateToLine(Contains("commit-05")).
			ClickAndHold(1, 1).
			MouseMove(1, 6).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("commit-05").IsSelected(),
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("drop here"),
				Contains("─── Commits"),
				Contains("commit-01"),
			).
			MouseRelease().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-05").IsSelected(),
				Contains("─── Commits"),
				Contains("commit-01"),
			).
			NavigateToLine(Contains("commit-01")).
			ClickAndHold(1, 6).
			MouseMove(1, 4).
			SelectedLines(
				Contains("commit-05"),
				Contains("─── Commits"),
				Contains("commit-01"),
			).
			MouseRelease().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-05").IsSelected(),
				Contains("─── Commits").IsSelected(),
				Contains("commit-01").IsSelected(),
			)
	},
})
