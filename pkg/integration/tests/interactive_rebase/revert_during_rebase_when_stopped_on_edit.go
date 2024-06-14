package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertDuringRebaseWhenStoppedOnEdit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Revert a series of commits while stopped in a rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// TODO: use our revert UI once we support range-select for reverts
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "commits",
				Command: "git -c core.editor=: revert HEAD^ HEAD^^",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("master commit")
		shell.NewBranch("branch")
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 04").IsSelected(),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit"),
			).
			NavigateToLine(Contains("commit 03")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("commit 04"),
				Contains("<-- YOU ARE HERE --- commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit"),
			).
			Press("X").
			Lines(
				/* EXPECTED:
				Contains("pick").Contains("commit 04"),
				Contains(`<-- YOU ARE HERE --- Revert "commit 01"`).IsSelected(),
				Contains(`Revert "commit 02"`),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit"),
				ACTUAL: */
				Contains("pick").Contains("commit 04"),
				Contains("edit").Contains("<-- CONFLICT --- commit 03").IsSelected(),
				Contains(`Revert "commit 01"`),
				Contains(`Revert "commit 02"`),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
				Contains("master commit"),
			)
	},
})
