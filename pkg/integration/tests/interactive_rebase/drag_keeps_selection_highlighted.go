package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragKeepsSelectionHighlighted = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep the original commit range highlighted while dragging sideways",
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
			Press(keys.Universal.RangeSelectDown).
			ClickAndHold(1, 1).
			SelectedLines(
				Contains("commit-05"),
				Contains("commit-04"),
				Contains("commit-03"),
			).
			MouseMove(10, 1).
			SelectedLines(
				Contains("commit-05"),
				Contains("commit-04"),
				Contains("commit-03"),
			).
			MouseRelease()
	},
})
