package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordLastCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords the last (HEAD) commit",
	ExtraCmdArgs: "",
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
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Reword commit")).
					InitialText(Equals("commit 02")).
					Clear().
					Type("renamed 02").
					Confirm()
			}).
			Lines(
				Contains("renamed 02"),
				Contains("commit 01"),
			)
	},
})
