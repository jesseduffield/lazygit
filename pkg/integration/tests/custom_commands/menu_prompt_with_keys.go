package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MenuPromptWithKeys = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command with menu options that have keybindings",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: `echo {{.Form.Choice | quote}} > result.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Choice",
						Type:  "menu",
						Title: "Choose an option",
						Options: []config.CustomCommandMenuOption{
							{
								Name:        "first",
								Description: "First option",
								Value:       "FIRST",
								Key:         "1",
							},
							{
								Name:        "second",
								Description: "Second option",
								Value:       "SECOND",
								Key:         "H",
							},
							{
								Name:        "third",
								Description: "Third option",
								Value:       "THIRD",
								Key:         "3",
							},
						},
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press("a")

		t.ExpectPopup().Menu().
			Title(Equals("Choose an option"))

		// 'H' is normally a navigation key (ScrollLeft), so this tests that menu item
		// keybindings have proper precedence over non-essential navigation keys
		t.Views().Menu().Press("H")

		t.FileSystem().FileContent("result.txt", Equals("SECOND\n"))
	},
})
