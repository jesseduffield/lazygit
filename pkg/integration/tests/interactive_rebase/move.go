package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Move = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Directly move a commit all the way down and all the way back up",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 04").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 04").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("commit 04").IsSelected(),
			).
			// assert nothing happens upon trying to move beyond the last commit
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("commit 04").IsSelected(),
			).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 04").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 03"),
				Contains("commit 04").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			// assert nothing happens upon trying to move beyond the first commit
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
