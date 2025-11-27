package remote

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddForkRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the 'Add fork remote' command to add a fork remote and check out a branch from it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit")
		shell.CloneIntoRemote("origin")
		shell.NewBranch("feature")
		shell.Clone("fork")
		shell.Checkout("master")
		shell.RemoveBranch("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin").IsSelected(),
			).
			Press(keys.Branches.AddForkRemote)

		t.ExpectPopup().Prompt().
			Title(Equals("Fork owner (username/org). Use username:branch to check out a branch")).
			Type("fork:feature").
			Confirm()

		t.Views().Remotes().
			Lines(
				Contains("origin"),
				Contains("fork").IsSelected(),
			)

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("feature ✓"),
				Contains("master"),
			)
	},
})
