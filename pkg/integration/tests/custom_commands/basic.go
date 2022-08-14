package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Basic = components.NewIntegrationTest(components.NewIntegrationTestArgs{
	Description:  "Using a custom command to create a new file",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo:    func(shell *components.Shell) {},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: "touch myfile",
			},
		}
	},
	Run: func(
		shell *components.Shell,
		input *components.Input,
		assert *components.Assert,
		keys config.KeybindingConfig,
	) {
		assert.WorkingTreeFileCount(0)

		input.PressKeys("a")
		assert.WorkingTreeFileCount(1)
		assert.SelectedLineContains("myfile")
	},
})
