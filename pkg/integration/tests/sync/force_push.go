package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ForcePush = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push to a remote with new commits, requiring a force push",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")

		// remove the 'two' commit so that we have something to pull from the remote
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("↓1 repo → master"))

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

		t.Views().Remotes().Focus().
			Lines(Contains("origin")).
			PressEnter()

		t.Views().RemoteBranches().IsFocused().
			Lines(Contains("master")).
			PressEnter()

		t.Views().SubCommits().IsFocused().
			Lines(Contains("one"))
	},
})
