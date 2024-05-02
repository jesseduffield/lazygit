package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropWithCustomCommentChar = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit with the 'core.commentChar' option set to a custom character",
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
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Equals("Are you sure you want to drop the selected commit(s)?")).
					Confirm()
			}).
			Lines(
				Contains("commit 01").IsSelected(),
			)
	},
})
