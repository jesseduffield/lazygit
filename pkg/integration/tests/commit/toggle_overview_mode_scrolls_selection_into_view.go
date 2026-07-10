package commit

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleOverviewModeScrollsSelectionIntoView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enabling overview mode scrolls the moved selection into view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// A long tail of tagged commits keeps the overview list taller than
		// the viewport, so the view can't get away with not scrolling; a long
		// run of plain commits on top lets us scroll the selection far away
		// from where it ends up after collapsing.
		shell.EmptyCommit("base")
		for i := 1; i <= 80; i++ {
			shell.EmptyCommit(fmt.Sprintf("tagged-%02d", i))
			shell.CreateLightweightTag(fmt.Sprintf("tag-%02d", i), "HEAD")
		}
		for i := 1; i <= 60; i++ {
			shell.EmptyCommit(fmt.Sprintf("plain-%02d", i))
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// plain-01 is deep down the list, so navigating to it scrolls the
		// view; every commit above it is plain except the HEAD commit, so
		// collapsing moves the selection to the very top of the overview,
		// which must scroll back into view.
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("plain-01")).
			Press(keys.Commits.ToggleOverviewMode).
			SelectedLine(Contains("plain-60")).
			TopLines(
				Contains("plain-60"),
			)
	},
})
