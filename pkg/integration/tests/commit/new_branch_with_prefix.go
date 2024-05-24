package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchWithPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Creating a new branch from a commit with a default name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.Git.BranchPrefix = "myprefix/"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit 1").
			EmptyCommit("commit 2").
			EmptyCommit("commit 3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 3").IsSelected(),
				Contains("commit 2"),
				Contains("commit 1"),
			).
			SelectNextItem().
			Press(keys.Universal.New).
			Tap(func() {
				branchName := "my-branch-name"
				t.ExpectPopup().Prompt().Title(Contains("New branch name")).Type(branchName).Confirm()
				t.Git().CurrentBranchName("myprefix/" + branchName)
			}).
			Lines(
				Contains("commit 2"),
				Contains("commit 1"),
			)
	},
})
