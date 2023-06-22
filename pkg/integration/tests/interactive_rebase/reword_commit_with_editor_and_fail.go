package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordCommitWithEditorAndFail = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords a commit with editor, and fails because an empty commit message is given",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			SetConfig("core.editor", "sh -c 'echo </dev/null >.git/COMMIT_EDITMSG'")
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
			Press(keys.Commits.RenameCommitWithEditor).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Reword in editor")).
					Content(Contains("Are you sure you want to reword this commit in your editor?")).
					Confirm()
			}).
			Lines(
				Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- commit 02").IsSelected(),
				Contains("commit 01"),
			)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("exit status 1"))
	},
})
