package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedCommitRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the {{ .SelectedCommitRange }} template variable",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "global",
				Command: `git log --format="%s" {{.SelectedCommitRange.From}}^..{{.SelectedCommitRange.To}} > file.txt`,
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().
			Lines(
				Contains("commit 03").IsSelected(),
				Contains("commit 02"),
				Contains("commit 01"),
			)

		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 03\n"))

		t.Views().Commits().Focus().
			Press(keys.Universal.RangeSelectDown)

		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 03\ncommit 02\n"))
	},
})
