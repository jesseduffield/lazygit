package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashFixupsAbove = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squashes all fixups above a commit and checks that the selected line stays correct.",
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
					Select(Contains("fixup! commit")).
					Confirm()
			}).
			Lines(
				Contains("fixup! commit 02"),
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			).
			Press(keys.Commits.SquashAboveCommits).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Apply fixup commits")).
					Select(Contains("Above the selected commit")).
					Confirm()
			}).
			Lines(
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			)

		t.Views().Main().
			Content(Contains("fixup content"))
	},
})
