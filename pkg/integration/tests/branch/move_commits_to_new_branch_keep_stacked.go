package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveCommitsToNewBranchKeepStacked = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch from the commits that you accidentally made on the wrong branch; choosing stacked on current branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.BranchPrefix = "myprefix/"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CloneIntoRemote("origin")
		shell.PushBranchAndSetUpstream("origin", "master")
		shell.NewBranch("feature")
		shell.EmptyCommit("feature branch commit")
		shell.PushBranchAndSetUpstream("origin", "feature")
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("new commit 1")
		shell.EmptyCommit("new commit 2")
		shell.UpdateFile("file1", "file1 changed")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Lines(
				Contains("M file1"),
			)
		t.Views().Branches().
			Focus().
			Lines(
				Contains("feature ↑2").IsSelected(),
				Contains("master ✓"),
			).
			Press(keys.Branches.MoveCommitsToNewBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Move commits to new branch")).
			Select(Contains("New branch stacked on current branch (feature)")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("New branch name (branch is off of 'feature')")).
			InitialText(Equals("myprefix/")).
			Type("new branch").
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("myprefix/new-branch").DoesNotContain("↑").IsSelected(),
				Contains("feature ✓"),
				Contains("master ✓"),
			)

		t.Views().Commits().
			Lines(
				Contains("new commit 2").IsSelected(),
				Contains("new commit 1"),
				Contains("* feature branch commit"),
				Contains("initial commit"),
			)
		t.Views().Files().
			Lines(
				Contains("M file1"),
			)
	},
})
