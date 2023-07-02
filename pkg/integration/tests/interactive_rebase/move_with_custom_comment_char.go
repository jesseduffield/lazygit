package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveWithCustomCommentChar = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Directly moves a commit down and back up with the 'core.commentChar' option set to a custom character",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("core.commentChar", ";")
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit 01"),
				Contains("commit 02").IsSelected(),
			).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			)
	},
})
