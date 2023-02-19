package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func createTwoBranchesReadyToForcePush(shell *Shell) {
	shell.EmptyCommit("one")
	shell.EmptyCommit("two")

	shell.NewBranch("other_branch")

	shell.CloneIntoRemote("origin")

	shell.SetBranchUpstream("master", "origin/master")
	shell.SetBranchUpstream("other_branch", "origin/other_branch")

	// remove the 'two' commit so that we have something to pull from the remote
	shell.HardReset("HEAD^")

	shell.Checkout("master")
	// doing the same for master
	shell.HardReset("HEAD^")
}

var ForcePushMultipleUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Force push to only the upstream branch of the current branch because the user has push.default upstream",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("push.default", "upstream")

		createTwoBranchesReadyToForcePush(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("↓1 repo → master"))

		t.Views().Branches().
			Lines(
				Contains("master ↓1"),
				Contains("other_branch ↓1"),
			)

		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Equals("Your branch has diverged from the remote branch. Press 'esc' to cancel, or 'enter' to force push.")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("✓ repo → master"))

		t.Views().Branches().
			Lines(
				Contains("master ✓"),
				Contains("other_branch ↓1"),
			)
	},
})
