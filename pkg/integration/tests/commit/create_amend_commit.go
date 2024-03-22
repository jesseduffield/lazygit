package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateAmendCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create an amend commit for an existing commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			CreateFileAndAdd("fixup-file", "fixup content")
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
			Press(keys.Commits.CreateFixupCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create fixup commit")).
					Select(Contains("amend! commit with changes")).
					Confirm()
				t.ExpectPopup().CommitMessagePanel().
					Content(Equals("commit 02")).
					Type(" amended").Confirm()
			}).
			Lines(
				Contains("amend! commit 02"),
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			)

		if t.Git().Version().IsAtLeast(2, 32, 0) { // Support for auto-squashing "amend!" commits was added in git 2.32.0
			t.Views().Commits().
				Press(keys.Commits.SquashAboveCommits).
				Tap(func() {
					t.ExpectPopup().Menu().
						Title(Equals("Apply fixup commits")).
						Select(Contains("Above the selected commit")).
						Confirm()
				}).
				Lines(
					Contains("commit 03"),
					Contains("commit 02 amended").IsSelected(),
					Contains("commit 01"),
				)
		}
	},
})
