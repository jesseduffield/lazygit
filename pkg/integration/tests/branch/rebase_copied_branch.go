package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseCopiedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Make a copy of a branch, rebase it, check that the original branch is unaffected",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.AppState.GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("master 1").
			EmptyCommit("master 2").
			NewBranchFrom("branch1", "master^").
			EmptyCommit("branch 1").
			EmptyCommit("branch 2").
			NewBranch("branch2")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("CI * branch 2"),
			Contains("CI branch 1"),
			Contains("CI master 1"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch2").IsSelected(),
				Contains("branch1"),
				Contains("master"),
			).
			NavigateToLine(Contains("master")).
			Press(keys.Branches.RebaseBranch).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Rebase 'branch2'")).
					Select(Contains("Simple rebase")).
					Confirm()
			})

		t.Views().Commits().Lines(
			Contains("CI branch 2"),
			Contains("CI branch 1"),
			Contains("CI master 2"),
			Contains("CI master 1"),
		)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch1")).
			PressPrimaryAction()

		t.Views().Commits().Lines(
			Contains("CI branch 2"),
			Contains("CI branch 1"),
			Contains("CI master 1"),
		)
	},
})
