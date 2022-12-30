package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenameBranchAndPull = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rename a branch to no longer match its upstream, then pull from the upstream",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Git.AutoFetch = false
	},
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

		t.Views().Branches().
			Focus().
			Lines(
				Contains("master"),
			).
			Press(keys.Branches.RenameBranch).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("rename branch")).
					Content(Equals("This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Enter new branch name")).
					InitialText(Equals("master")).
					Type("-local").
					Confirm()
			}).
			Press(keys.Universal.PullFiles)

		t.Views().Commits().
			Lines(
				Contains("two"),
				Contains("one"),
			)
	},
})
