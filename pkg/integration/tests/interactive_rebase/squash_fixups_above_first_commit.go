package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashFixupsAboveFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squashes all fixups above the first (initial) commit.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2).
			CreateFileAndAdd("fixup-file", "fixup content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Commits.CreateFixupCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create fixup commit")).
					Select(Contains("fixup! commit")).
					Confirm()
			}).
			NavigateToLine(Contains("commit 01").DoesNotContain("fixup!")).
			Press(keys.Commits.SquashAboveCommits).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Apply fixup commits")).
					Select(Contains("Above the selected commit")).
					Confirm()
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("fixup content"))
	},
})
