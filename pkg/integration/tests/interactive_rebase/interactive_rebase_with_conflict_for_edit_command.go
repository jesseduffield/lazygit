package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InteractiveRebaseWithConflictForEditCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase a branch interactively, and edit a commit that will conflict",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file.txt", "master content")
		shell.Commit("master commit")
		shell.NewBranchFrom("branch", "master^")
		shell.CreateNCommits(3)
		shell.CreateFileAndAdd("file.txt", "branch content")
		shell.Commit("this will conflict")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("this will conflict").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("initial commit"),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("master")).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'branch'")).
			Select(Contains("Interactive rebase")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("this will conflict")).
			Press(keys.Universal.Edit)

		t.Common().ContinueRebase()
		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Cancel()

		t.Views().Commits().
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("edit").Contains("<-- CONFLICT --- this will conflict").IsSelected(),
				Contains("--- Commits ---"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit"),
				Contains("initial commit"),
			)
	},
})
