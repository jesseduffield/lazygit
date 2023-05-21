package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordYouAreHereCommitWithEditor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords the current HEAD commit in an interactive rebase with editor",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			SetConfig("core.editor", "sh -c 'echo renamed 02 >.git/COMMIT_EDITMSG'")
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
			Press(keys.Commits.RenameCommitWithEditor).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Reword in editor")).
					Content(Contains("Are you sure you want to reword this commit in your editor?")).
					Confirm()
			}).
			Lines(
				Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- renamed 02").IsSelected(),
				Contains("commit 01"),
			)
	},
})
