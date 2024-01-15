package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashDownFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Tries to squash down the first commit, which results in an error message",
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
			Press(keys.Commits.SquashDown).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: There's no commit below to squash into"))
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
