package custom_commands

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var BasicCmdFromConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command to create a new file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: "touch myfile",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press("a").
			Lines(
				Contains("myfile"),
			)
	},
})
