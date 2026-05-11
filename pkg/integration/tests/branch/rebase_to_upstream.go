package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebaseToUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rebase the current branch to the selected branch upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			EmptyCommit("ensure-master").
			EmptyCommit("to-be-added"). // <- this will only exist remotely
			PushBranchAndSetUpstream("origin", "master").
			RenameCurrentBranch("master-local").
			HardReset("HEAD~1").
			NewBranchFrom("base-branch", "master-local").
			EmptyCommit("base-branch-commit").
			NewBranch("target").
			EmptyCommit("target-commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("target-commit"),
			Contains("base-branch-commit"),
			Contains("ensure-master"),
		)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("target").IsSelected(),
				Contains("base-branch"),
				Contains("master-local"),
			).
			SelectNextItem().
			Lines(
				Contains("target"),
				Contains("base-branch").IsSelected(),
				Contains("master-local"),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Rebase checked-out branch onto upstream of selected branch")).
					Tooltip(Contains("Disabled: The selected branch has no upstream (or the upstream is not stored locally)")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("Disabled: The selected branch has no upstream (or the upstream is not stored locally)"))
					}).
					Cancel()
			}).
			SelectNextItem().
			Lines(
				Contains("target"),
				Contains("base-branch"),
				Contains("master-local").IsSelected(),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Rebase checked-out branch onto origin/master...")).
					Confirm()
				t.ExpectPopup().Menu().
					Title(Equals("Rebase 'target'")).
					Select(Contains("Simple rebase")).
					Confirm()
			})

		t.Views().Commits().Lines(
			Contains("target-commit"),
			Contains("base-branch-commit"),
			Contains("to-be-added"),
			Contains("ensure-master"),
		)
	},
})
