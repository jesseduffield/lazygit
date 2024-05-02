package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SortRemoteBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sort remote branches alphabetically or by date",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("first")
		shell.EmptyCommitWithDate("commit", "2023-04-07 10:00:00")
		shell.NewBranch("second")
		shell.EmptyCommitWithDate("commit", "2023-04-07 12:00:00")
		shell.NewBranch("third")
		shell.EmptyCommitWithDate("commit", "2023-04-07 11:00:00")
		shell.CloneIntoRemote("origin")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin").IsSelected(),
			).
			PressEnter()

		// sorted alphabetically by default
		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("first").IsSelected(),
				Contains("second"),
				Contains("third"),
			).
			SelectNextItem() // to test that the selection jumps back to the first when sorting

		t.Views().RemoteBranches().
			Press(keys.Branches.SortOrder)

		t.ExpectPopup().Menu().Title(Equals("Sort order")).
			Select(Contains("-committerdate")).
			Confirm()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("second").IsSelected(),
				Contains("third"),
				Contains("first"),
			)
	},
})
