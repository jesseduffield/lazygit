package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchWhenSortedByDate = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fetch a branch while sort order is by date; verify that branch stays selected",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithDate("commit", "2023-04-07 10:00:00"). // first master commit, older than branch2
			EmptyCommitWithDate("commit", "2023-04-07 12:00:00"). // second master commit, newer than branch2
			NewBranch("branch1").                                 // branch1 will be checked out, so its date doesn't matter
			EmptyCommitWithDate("commit", "2023-04-07 11:00:00"). // branch2 commit, date is between the two master commits
			NewBranch("branch2").
			Checkout("master").
			CloneIntoRemote("origin").
			SetBranchUpstream("master", "origin/master"). // upstream points to second master commit
			HardReset("HEAD^").                           // rewind to first master commit
			Checkout("branch1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().Title(Equals("Sort order")).
			Select(Contains("-committerdate")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("* branch1").IsSelected(),
				Contains("branch2"),
				Contains("master ↓1"),
			).
			NavigateToLine(Contains("master")).
			Press(keys.Branches.FetchRemote).
			Lines(
				Contains("* branch1"),
				Contains("master").IsSelected(),
				Contains("branch2"),
			)
	},
})
