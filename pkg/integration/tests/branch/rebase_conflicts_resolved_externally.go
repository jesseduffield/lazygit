package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RebaseConflictsResolvedExternally = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resolve conflicts and continue rebase externally, verify popup auto-dismisses.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'first-change-branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Files().
			IsFocused().
			Tap(func() {
				t.Shell().UpdateFile("file", shared.FirstChangeFileContent)
				t.Shell().GitAdd("file")
			}).
			Press(keys.Universal.Refresh)

		t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue the rebase?"))

		t.Shell().RunCommandWithEnv([]string{"git", "rebase", "--continue"}, []string{
			"GIT_EDITOR=true",
			"GIT_SEQUENCE_EDITOR=true",
			"EDITOR=true",
			"VISUAL=true",
			"GIT_TERMINAL_PROMPT=0",
		})

		t.GlobalPress(keys.Universal.Refresh)
		t.Wait(500) // Allow time for the refresh to complete and the popup to auto-close

		t.Views().Files().IsFocused()
		t.Views().Information().Content(DoesNotContain("Rebasing"))
	},
})
