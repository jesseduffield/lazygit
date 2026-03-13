package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Keybindings = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Display the keybindings with custom commands",
	SetupRepo:   func(shell *Shell) {},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: "touch myfile",
			},
			{
				Key:     "p",
				Context: "files",
				Command: "fake pull",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Press("?")

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			ContainsLines(
				Contains("a touch myfile"),
				DoesNotContain("a Stage all"),
			)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			ContainsLines(
				Contains("p fake pull"),
			)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			ContainsLines(
				Contains("p Pull"),
			)
	},
})
