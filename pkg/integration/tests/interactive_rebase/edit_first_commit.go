package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edits the first commit, just to show that it's possible",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit 02"),
				MatchesRegexp("YOU ARE HERE.*commit 01").IsSelected(),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
