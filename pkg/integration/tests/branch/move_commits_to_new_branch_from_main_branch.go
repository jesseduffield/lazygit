package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveCommitsToNewBranchFromMainBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch from the commits that you accidentally made on master",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CloneIntoRemote("origin")
		shell.PushBranchAndSetUpstream("origin", "master")
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
				Contains("master ↑2").IsSelected(),
			).
			Press(keys.Branches.MoveCommitsToNewBranch)

		t.ExpectPopup().Confirmation().
			Title(Equals("Move commits to new branch")).
			Content(Contains("This will take all unpushed commits and move them to a new branch (off of master).")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("New branch name (branch is off of 'master')")).
			Type("new branch").
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("new-branch").DoesNotContain("↑").IsSelected(),
				Contains("master ✓"),
			)

		t.Views().Commits().
			Lines(
				Contains("new commit 2").IsSelected(),
				Contains("new commit 1"),
				Contains("initial commit"),
			)
		t.Views().Files().
			Lines(
				Contains("M file1"),
			)
	},
})
