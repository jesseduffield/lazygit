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
			// The following behavior requires correction:
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("failed to parse line")).
					Confirm()
			}).
			Lines(
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			)
	},
})
