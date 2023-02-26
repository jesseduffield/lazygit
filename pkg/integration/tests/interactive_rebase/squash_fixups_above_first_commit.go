package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashFixupsAboveFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squashes all fixups above the first (initial) commit.",
	ExtraCmdArgs: "",
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
				t.ExpectPopup().Confirmation().
					Title(Equals("Create fixup commit")).
					Content(Contains("Are you sure you want to create a fixup! commit for commit")).
					Confirm()
			}).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Commits.SquashAboveCommits).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Squash all 'fixup!' commits above selected commit (autosquash)")).
					Content(Contains("Are you sure you want to squash all fixup! commits above")).
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
