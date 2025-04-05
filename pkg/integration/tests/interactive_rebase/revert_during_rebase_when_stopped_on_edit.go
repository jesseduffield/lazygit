package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertDuringRebaseWhenStoppedOnEdit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Revert a series of commits while stopped in a rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("master commit 1")
		shell.EmptyCommit("master commit 2")
		shell.NewBranch("branch")
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit 2"),
				Contains("master commit 1"),
			).
			NavigateToLine(Contains("commit 03")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick").Contains("commit 04"),
				Contains("--- Commits ---"),
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit 2"),
				Contains("master commit 1"),
			).
			SelectNextItem().
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.RevertCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Revert commit")).
					Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
					Confirm()
			}).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick").Contains("commit 04"),
				Contains("--- Commits ---"),
				Contains(`Revert "commit 01"`),
				Contains(`Revert "commit 02"`),
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("commit 01").IsSelected(),
				Contains("master commit 2"),
				Contains("master commit 1"),
			)
	},
})
