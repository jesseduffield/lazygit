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
				Key:         "x",
				Description: "My Custom Commands",
				CommandMenu: []config.CustomCommand{
					{
						Key:     "j",
						Context: "global",
						Command: "echo j",
						Output:  "popup",
					},
					{
						Key:     "H",
						Context: "global",
						Command: "echo H",
						Output:  "popup",
					},
					{
						Key:     "y",
						Context: "global",
						Command: "echo y",
						Output:  "popup",
					},
					{
						Key:     "<down>",
						Context: "global",
						Command: "echo down",
						Output:  "popup",
					},
				},
			},
		}
		cfg.GetUserConfig().Keybinding.Universal.ConfirmMenu = "y"
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
						Contains("j echo j"),
						Contains("H echo H"),
						Contains("  echo y"),
						Contains("  echo down"),
					)
				t.GlobalPress("j")
				/* EXPECTED:
				t.ExpectPopup().Alert().Title(Equals("echo j")).Content(Equals("j")).Confirm()
				ACTUAL: */
				// The menu stays open; 'j' didn't trigger the command; instead, it selected the
				// next item, which we can confirm by pressing enter:
				t.GlobalPress(keys.Universal.ConfirmMenu)
				t.ExpectPopup().Alert().Title(Equals("echo H")).Content(Equals("H")).Confirm()
			}).
			Press("x").
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands"))
				t.GlobalPress("H")
				/* EXPECTED:
				t.ExpectPopup().Alert().Title(Equals("echo H")).Content(Equals("H")).Confirm()
				ACTUAL: */
				// The menu stays open:
				t.ExpectPopup().Menu().
					Title(Equals("My Custom Commands"))
			})
	},
})
