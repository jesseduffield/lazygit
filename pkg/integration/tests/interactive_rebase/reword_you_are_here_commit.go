package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordYouAreHereCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords the current HEAD commit in an interactive rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- commit 02").IsSelected(),
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
				Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- renamed 02").IsSelected(),
				Contains("commit 01"),
			)
	},
})
