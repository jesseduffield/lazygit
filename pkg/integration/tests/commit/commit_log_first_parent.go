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
			Wait(500).
			SelectNextItem().
			SelectedLine(Contains("--first-parent")).
			Wait(500).
			PressPrimaryAction()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("feat/b"),
				Contains("feat/a"),
				Contains("MASTER"),
			).
			Wait(500)

		t.Views().Branches().
			Focus().
			Wait(500).
			PressEnter().
			Wait(500)

		t.Views().SubCommits().
			Press(keys.Commits.OpenLogMenu)

		t.Views().Menu().
			SelectedLine(Contains("--first-parent")).
			PressPrimaryAction()

		t.Views().SubCommits().
			Lines(
				Contains("feat/b"),
				Contains("feat/a"),
				Contains("MASTER"),
			)
	},
})
