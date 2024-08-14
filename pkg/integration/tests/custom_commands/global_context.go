package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GlobalContext = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ensure global context works",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("my change")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:        "X",
				Context:    "global",
				Command:    "touch myfile",
				ShowOutput: false,
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// commits
		t.Views().Commits().
			Focus().
			Press("X")

		t.Views().Files().
			Focus().
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")
		t.GlobalPress(keys.Files.RefreshFiles)

		// branches
		t.Views().Branches().
			Focus().
			Press("X")

		t.Views().Files().
			Focus().
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")
		t.GlobalPress(keys.Files.RefreshFiles)

		// files
		t.Views().Files().
			Focus().
			Press("X")

		t.Views().Files().
			Focus().
			Lines(Contains("myfile"))

		t.Shell().DeleteFile("myfile")
	},
})
