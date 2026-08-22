package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RestoreUpstreamBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Restore an upstream branch that was deleted on the remote",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			EmptyCommit("base commit").
			NewBranch("feature").
			EmptyCommit("on feature").
			PushBranchAndSetUpstream("origin", "feature").
			Checkout("master").
			RunCommand([]string{"git", "-C", "../origin", "branch", "-D", "feature"}).
			RunCommand([]string{"git", "fetch", "origin", "--prune"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature").Contains("upstream gone"),
			)

		t.Views().Branches().
			NavigateToLine(Contains("feature")).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Restore upstream branch")).
					Confirm()
			})

		// the "upstream gone" message is gone and the remote branch is recreated
		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("feature").DoesNotContain("upstream gone"),
			)

		t.Git().AssertRemoteBranchExists("origin", "feature")
	},
})
