package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Rename = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rename a branch, replacing spaces in the name with dashes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master"),
			).
			Press(keys.Branches.RenameBranch).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Contains("Enter new branch name")).
					InitialText(Equals("master")).
					Clear().
					Type("new branch name").
					Confirm()
			}).
			Lines(
				Contains("new-branch-name"),
			)
	},
})
