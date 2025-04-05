package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendWhenThereAreConflictsAndAmend = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends the last commit from the files panel while a rebase is stopped due to conflicts, and amends the commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setupForAmendTests(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		doTheRebaseForAmendTests(t, keys)

		t.Views().Files().
			Press(keys.Commits.AmendToCommit)

		t.ExpectPopup().Menu().
			Title(Equals("Amend commit")).
			Select(Equals("Yes, amend previous commit")).
			Confirm()

		t.Views().Files().IsEmpty()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick").Contains("commit three"),
				Contains("pick").Contains("<-- CONFLICT --- file1 changed in branch"),
				Contains("--- Commits ---"),
				Contains("commit two"),
				Contains("file1 changed in master"),
				Contains("base commit"),
			)

		checkCommitContainsChange(t, "commit two", "+branch")
	},
})
