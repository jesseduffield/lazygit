package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to the first (initial) commit.",
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
			Press(keys.Commits.AmendToCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Amend commit")).
					Content(Contains("Are you sure you want to amend this commit with your staged files?")).
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
