package worktree

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var CustomCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that custom commands work with worktrees by deleting a worktree via a custom command",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "d",
				Context: "worktrees",
				Command: "git worktree remove {{ .SelectedWorktree.Path | quote }}",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)"),
				Contains("linked-worktree"),
			).
			NavigateToLine(Contains("linked-worktree")).
			Press("d").
			Lines(
				Contains("repo (main)"),
			)
	},
})
