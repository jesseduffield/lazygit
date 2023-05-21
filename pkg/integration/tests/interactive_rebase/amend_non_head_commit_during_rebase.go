package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendNonHeadCommitDuringRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Tries to amend a commit that is not the head while already rebasing, resulting in an error message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- commit 02"),
				Contains("commit 01"),
			)

		for _, commit := range []string{"commit 01", "commit 03"} {
			t.Views().Commits().
				NavigateToLine(Contains(commit)).
				Press(keys.Commits.AmendToCommit)

			t.ExpectPopup().Alert().
				Title(Equals("Error")).
				Content(Contains("Can't perform this action during a rebase")).
				Confirm()
		}
	},
})
