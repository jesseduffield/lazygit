package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestSelectRemoteAndTargetBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open up a pull request, specifying a remote and target branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.OpenLink = "echo {{link}} > /tmp/openlink"
	},
	SetupRepo: func(shell *Shell) {
		// Create an initial commit ('git branch set-upstream-to' bails out otherwise)
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")

		// Create a new branch and a remote that has that branch
		shell.NewBranch("branch-1")
		shell.CloneIntoRemote("upstream")

		// Create another branch and a second remote. The first remote doesn't have this branch.
		shell.NewBranch("branch-2")
		shell.CloneIntoRemote("origin")

		// To allow a pull request to be created from a branch, it must have an upstream set.
		shell.SetBranchUpstream("branch-2", "origin/branch-2")

		shell.RunCommand([]string{"git", "remote", "set-url", "origin", "https://github.com/my-personal-fork/lazygit"})
		shell.RunCommand([]string{"git", "remote", "set-url", "upstream", "https://github.com/jesseduffield/lazygit"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Open a PR for the current branch (i.e. 'branch-2')
		t.Views().
			Branches().
			Focus().
			Press(keys.Branches.ViewPullRequestOptions)

		t.ExpectPopup().
			Menu().
			Title(Equals("View create pull request options")).
			Select(Contains("Select branch")).
			Confirm()

		// Verify that we're prompted to enter the remote
		t.ExpectPopup().
			Prompt().
			Title(Equals("Select target remote")).
			SuggestionLines(
				Equals("origin"),
				Equals("upstream")).
			ConfirmSuggestion(Equals("upstream"))

		// Verify that we're prompted to enter the target branch and that only those branches
		// present in the selected remote are listed as suggestions (i.e. 'branch-2' is not there).
		t.ExpectPopup().
			Prompt().
			Title(Equals("branch-2 → upstream/")).
			SuggestionLines(
				Equals("branch-1"),
				Equals("master")).
			ConfirmSuggestion(Equals("master"))

		// Verify that the expected URL is used (by checking the openlink file)
		//
		// Please note that when targeting a different remote - like it's done here in this test -
		// the link is not yet correct. Thus, this test is expected to fail once this is fixed.
		t.FileSystem().FileContent(
			"/tmp/openlink",
			Equals("https://github.com/my-personal-fork/lazygit/compare/master...branch-2?expand=1\n"))
	},
})
