package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitLogFirstParent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show only the first parent of the commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) { },
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("MASTER")
		shell.NewBranch("feat/a")
		shell.EmptyCommit("A")
		shell.EmptyCommit("AA")
		shell.Checkout("master")
		shell.Merge("feat/a")
		shell.NewBranchFrom("feat/b", "master")
		shell.EmptyCommit("BB")
		shell.EmptyCommit("BBB")

		shell.Checkout("master")
		shell.Merge("feat/b")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Commits.OpenLogMenu)

		t.Views().Menu().
			Wait(1000).
			SelectNextItem().
			Wait(1000).
			PressPrimaryAction()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("feat/b"),
				Contains("feat/a"),
				Contains("MASTER"),
			)


		// t.Views().Files().
		// 	IsFocused().
		// 	Lines(
		// 		Equals("▼ /").IsSelected(),
		// 		Equals("  ?? myfile"),
		// 		Equals("  ?? myfile2"),
		// 	).
		// 	SelectNextItem().
		// 	PressPrimaryAction(). // stage file
		// 	Lines(
		// 		Equals("▼ /"),
		// 		Equals("  A  myfile").IsSelected(),
		// 		Equals("  ?? myfile2"),
		// 	).
		// 	SelectNextItem().
		// 	PressPrimaryAction(). // stage other file
		// 	Lines(
		// 		Equals("▼ /"),
		// 		Equals("  A  myfile"),
		// 		Equals("  A  myfile2").IsSelected(),
		// 	).
		// 	Press(keys.Files.CommitChanges)
		//
		// commitMessage := "my commit message"
		//
		// t.ExpectPopup().CommitMessagePanel().Type(commitMessage).Confirm()
		//
		// t.Views().Files().
		// 	IsEmpty()
		//
		// t.Views().Commits().
		// 	Focus().
		// 	Lines(
		// 		Contains(commitMessage).IsSelected(),
		// 	).
		// 	PressEnter()
		//
		// t.Views().CommitFiles().
		// 	IsFocused().
		// 	Lines(
		// 		Equals("▼ /"),
		// 		Equals("  A myfile"),
		// 		Equals("  A myfile2"),
		// 	)
	},
})
