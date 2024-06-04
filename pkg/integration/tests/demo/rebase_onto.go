package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseOnto = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase with '--onto' flag. We start with a feature branch on the develop branch that we want to rebase onto the master branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(60)
		shell.NewBranch("develop")

		shell.SetAuthor("Joe Blow", "joeblow@gmail.com")

		shell.RandomChangeCommit("Develop commit 1")
		shell.RandomChangeCommit("Develop commit 2")
		shell.RandomChangeCommit("Develop commit 3")

		shell.SetAuthor("Jesse Duffield", "jesseduffield@gmail.com")

		shell.NewBranch("feature/demo")

		shell.RandomChangeCommit("Feature commit 1")
		shell.RandomChangeCommit("Feature commit 2")
		shell.RandomChangeCommit("Feature commit 3")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("feature/demo", "origin/feature/demo")
		shell.SetBranchUpstream("develop", "origin/develop")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Rebase from marked base commit")
		t.Wait(1000)

		// first we focus the commits view, then expand to show the branches against each commit
		// Then we go back to normal value, mark the last develop branch commit as the marked commit
		// Then go to the branches view and press 'r' on the master branch to rebase onto it
		// then we force push our changes.

		t.Views().Commits().
			Focus().
			Press(keys.Universal.PrevScreenMode).
			Wait(500).
			NavigateToLine(Contains("Develop commit 3")).
			Wait(500).
			Press(keys.Commits.MarkCommitAsBaseForRebase).
			Wait(1000).
			Press(keys.Universal.NextScreenMode).
			Wait(500)

		t.Views().Branches().
			Focus().
			Wait(500).
			NavigateToLine(Contains("master")).
			Wait(500).
			Press(keys.Branches.RebaseBranch).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Contains("Rebase 'feature/demo' from marked base")).
					Select(Contains("Simple rebase")).
					Confirm()
			}).
			Wait(1000).
			Press(keys.Universal.Push).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Contains("Force push")).
					Content(AnyString()).
					Wait(500).
					Confirm()
			})
	},
})
