package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ForcePushTriangular = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push to a remote, requiring a force push because the branch is behind the remote push branch but not the upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.22.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("push.default", "current")

		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.NewBranch("feature")
		shell.SetBranchUpstream("feature", "origin/master")
		shell.EmptyCommit("two")
		shell.PushBranch("origin", "feature")

		// remove the 'two' commit so that we are behind the push branch
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("✓ repo → feature"))

		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		// This results in an attempt to push normally, which fails with an error:
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Updates were rejected. Please fetch and examine the remote changes before pushing again."))

		/* EXPECTED:
		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Equals("Your branch has diverged from the remote branch. Press <esc> to cancel, or <enter> to force push.")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("✓ repo → feature"))

		t.Views().Remotes().Focus().
			Lines(Contains("origin")).
			PressEnter()

		t.Views().RemoteBranches().IsFocused().
			Lines(
				Contains("feature"),
				Contains("master"),
			).
			PressEnter()

		t.Views().SubCommits().IsFocused().
			Lines(Contains("one"))
		*/
	},
})
