package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomCommandsSubmenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using custom commands from a custom commands menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:         "x",
				Description: "My Custom Commands",
				CommandMenu: []config.CustomCommand{
					{
						Key:     "1",
						Context: "global",
						Command: "touch myfile-global",
					},
					{
						Key:     "2",
						Context: "files",
						Command: "touch myfile-files",
					},
					{
						Key:     "3",
						Context: "commits",
						Command: "touch myfile-commits",
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			IsEmpty().
			Press("x").
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands")).
					Lines(
						Contains("1 touch myfile-global"),
						Contains("2 touch myfile-files"),
					).
					Select(Contains("touch myfile-files")).Confirm()
			}).
			Lines(
				Contains("myfile-files"),
			)

		t.Views().Commits().
			Focus().
			Press("x").
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands")).
					Lines(
						Contains("1 touch myfile-global"),
						Contains("3 touch myfile-commits"),
					)
				t.GlobalPress("3")
			})

		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Contains("myfile-commits"),
				Contains("myfile-files"),
			)
	},
})
