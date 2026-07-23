package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragToReorder = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drag a selected commit range multiple rows in one operation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(5)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.RangeSelectDown).
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			ClickAndHold(1, 1).
			MouseMove(1, 3).
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("drop here"),
				Contains("commit-01"),
			).
			PressEscape().
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			MouseMove(1, 4).
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			MouseRelease().
			ClickAndHold(1, 1).
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			MouseMove(1, 3).
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("drop here"),
				Contains("commit-01"),
			).
			SelectNextItem().
			SelectedLines(
				Contains("commit-03"),
			).
			MouseRelease().
			TopLines(
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-01"),
			).
			ClickAndHold(1, 2).
			MouseMove(1, 0).
			TopLines(
				Contains("drop here"),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-01"),
			).
			MouseRelease().
			TopLines(
				Contains("commit-05").IsSelected(),
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			ClickAndHold(1, 1).
			MouseRelease().
			SelectedLines(
				Contains("commit-04"),
			)
	},
})
