package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FixupFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Tries to fixup the first commit, which results in an error message",
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
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Commits.MarkCommitAsFixup).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Equals("There's no commit below to squash into")).
					Confirm()
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
