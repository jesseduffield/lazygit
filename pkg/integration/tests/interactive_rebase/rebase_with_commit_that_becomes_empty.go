package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseWithCommitThatBecomesEmpty = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Performs a rebase involving a commit that becomes empty during the rebase, and gets dropped.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		// It is important that we create two separate commits for the two
		// changes to the file, but only one commit for the same changes on our
		// branch; otherwise, the commit would be discarded at the start of the
		// rebase already.
		shell.CreateFileAndAdd("file", "change 1\n")
		shell.Commit("master change 1")
		shell.UpdateFileAndAdd("file", "change 1\nchange 2\n")
		shell.Commit("master change 2")
		shell.NewBranchFrom("branch", "HEAD^^")
		shell.CreateFileAndAdd("file", "change 1\nchange 2\n")
		shell.Commit("branch change")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("master")).
			Press(keys.Branches.RebaseBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Rebase 'branch'")).
			Select(Contains("Simple rebase")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("master change 2"),
				Contains("master change 1"),
				Contains("initial commit"),
			)
	},
})
