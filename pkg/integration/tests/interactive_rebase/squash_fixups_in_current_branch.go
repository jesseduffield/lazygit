package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashFixupsInCurrentBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squashes all fixups in the current branch.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1", "file1").
			Commit("master commit").
			NewBranch("branch").
			// Test the pathological case that the first commit of a branch is a
			// fixup for the last master commit below it. We _don't_ want this to
			// be squashed.
			UpdateFileAndAdd("file1", "changed file1").
			Commit("fixup! master commit").
			CreateNCommits(2).
			CreateFileAndAdd("fixup-file", "fixup content").
			Commit("fixup! commit 01")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectNextItem().
			SelectNextItem().
			Lines(
				Contains("fixup! commit 01"),
				Contains("commit 02"),
				Contains("commit 01").IsSelected(),
				Contains("fixup! master commit"),
				Contains("master commit"),
			).
			Press(keys.Commits.SquashAboveCommits).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Apply fixup commits")).
					Select(Contains("In current branch")).
					Confirm()
			}).
			Lines(
				Contains("commit 02"),
				Contains("commit 01").IsSelected(),
				Contains("fixup! master commit"),
				Contains("master commit"),
			)

		t.Views().Main().
			Content(Contains("fixup content"))
	},
})
