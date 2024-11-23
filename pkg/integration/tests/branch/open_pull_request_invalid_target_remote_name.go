package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestInvalidTargetRemoteName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open up a pull request, specifying a non-existing target remote",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Create an initial commit ('git branch set-upstream-to' bails out otherwise)
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")

		// Create a new branch
		shell.NewBranch("branch-1")

		// Create a couple of remotes
		shell.CloneIntoRemote("upstream")
		shell.CloneIntoRemote("origin")

		// To allow a pull request to be created from a branch, it must have an upstream set.
		shell.SetBranchUpstream("branch-1", "origin/branch-1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Open a PR for the current branch (i.e. 'branch-1')
		t.Views().
			Branches().
			Focus().
			Press(keys.Branches.ViewPullRequestOptions)

		t.ExpectPopup().
			Menu().
			Title(Equals("View create pull request options")).
			Select(Contains("Select branch")).
			Confirm()

		// Verify that we're prompted to enter the remote and enter the name of a non-existing one.
		t.ExpectPopup().
			Prompt().
			Title(Equals("Select target remote")).
			Type("non-existing-remote").
			Confirm()

		// Verify that this leads to an error being shown (instead of progressing to branch selection).
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("A remote named 'non-existing-remote' does not exist")).
			Confirm()
	},
})
