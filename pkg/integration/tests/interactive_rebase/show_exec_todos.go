package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowExecTodos = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show exec todos in the rebase todo list",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "commits",
				Command: "git -c core.editor=: rebase -i -x false HEAD^^",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press("X").
			Tap(func() {
				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("Rebasing (2/4)Executing: false")).Confirm()
			}).
			Lines(
				Contains("exec").Contains("false"),
				Contains("pick").Contains("CI commit 03"),
				Contains("CI ◯ <-- YOU ARE HERE --- commit 02"),
				Contains("CI ◯ commit 01"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("exit status 1")).Confirm()
			}).
			Lines(
				Contains("CI ◯ <-- YOU ARE HERE --- commit 03"),
				Contains("CI ◯ commit 02"),
				Contains("CI ◯ commit 01"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("CI ◯ commit 03"),
				Contains("CI ◯ commit 02"),
				Contains("CI ◯ commit 01"),
			)
	},
})
