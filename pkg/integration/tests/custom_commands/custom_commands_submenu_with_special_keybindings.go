package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomCommandsSubmenuWithSpecialKeybindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using custom commands from a custom commands menu with keybindings that conflict with builtin ones",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:         config.Keybinding{"x"},
				Description: "My Custom Commands",
				CommandMenu: []config.CustomCommand{
					{
						Key:     config.Keybinding{"j"},
						Context: "global",
						Command: "echo j",
						Output:  "popup",
					},
					{
						Key:     config.Keybinding{"H"},
						Context: "global",
						Command: "echo H",
						Output:  "popup",
					},
					{
						Key:     config.Keybinding{"y"},
						Context: "global",
						Command: "echo y",
						Output:  "popup",
					},
					{
						Key:     config.Keybinding{"<down>"},
						Context: "global",
						Command: "echo down",
						Output:  "popup",
					},
				},
			},
		}
		cfg.GetUserConfig().Keybinding.Universal.ConfirmMenu = config.Keybinding{"y"}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			IsEmpty().
			Press(config.Keybinding{"x"}).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands")).
					Lines(
						Contains("j echo j"),
						Contains("H echo H"),
						Contains("  echo y"),
						Contains("  echo down"),
					)
				t.GlobalPress(config.Keybinding{"j"})
				t.ExpectPopup().Alert().Title(Equals("echo j")).Content(Equals("j")).Confirm()
			}).
			Press(config.Keybinding{"x"}).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands"))
				t.GlobalPress(config.Keybinding{"H"})
				t.ExpectPopup().Alert().Title(Equals("echo H")).Content(Equals("H")).Confirm()
			})
	},
})
