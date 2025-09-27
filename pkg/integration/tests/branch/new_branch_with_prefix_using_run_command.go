package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchWithPrefixUsingRunCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Creating a new branch with a branch prefix using a runCommand",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.BranchPrefix = "myprefix/{{ runCommand \"echo dynamic\" }}/"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit 1")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 1").IsSelected(),
			).
			SelectNextItem().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Contains("New branch name")).
					InitialText(Equals("myprefix/dynamic/")).
					Type("my-branch").
					Confirm()
				t.Git().CurrentBranchName("myprefix/dynamic/my-branch")
			})
	},
})
