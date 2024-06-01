package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ForcePushRemoteBranchNotStoredLocally = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a branch whose remote branch is not stored locally, requiring a force push",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.Clone("some-remote")

		// remove the 'two' commit so that we have something to pull from the remote
		shell.HardReset("HEAD^")

		shell.SetConfig("branch.master.remote", "../some-remote")
		shell.SetConfig("branch.master.pushRemote", "../some-remote")
		shell.SetConfig("branch.master.merge", "refs/heads/master")

		shell.CreateFileAndAdd("file1", "file1 content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Contains("? repo → master"))

		// We're behind our upstream now, so we expect to be asked to force-push
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Equals("Your branch has diverged from the remote branch. Press <esc> to cancel, or <enter> to force push.")).
			Confirm()

		// Make a new local commit
		t.Views().Files().IsFocused().Press(keys.Files.CommitChanges)
		t.ExpectPopup().CommitMessagePanel().Type("new").Confirm()

		t.Views().Commits().
			Lines(
				Contains("new"),
				Contains("one"),
			)

		// Pushing this works without needing to force push
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		// Now add the clone as a remote just so that we can check if what we
		// pushed arrived there correctly
		t.Views().Remotes().Focus().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Equals("New remote name:")).Type("some-remote").Confirm()
		t.ExpectPopup().Prompt().
			Title(Equals("New remote url:")).Type("../some-remote").Confirm()
		t.Views().Remotes().Lines(
			Contains("some-remote").IsSelected(),
		).
			PressEnter()

		t.Views().RemoteBranches().IsFocused().Lines(
			Contains("master").IsSelected(),
		).
			PressEnter()

		t.Views().SubCommits().IsFocused().Lines(
			Contains("new"),
			Contains("one"),
		)
	},
})
