package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Textbox = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command type multiline description",
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
				Command: `echo "{{.Form.Description}}" > output.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Description",
						Type:  "textbox",
						Title: "description",
					},
					{
						Type:  "confirm",
						Title: "Are you sure?",
						Body:  "Are you REALLY sure you want to make this file? Up to you buddy.",
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press("a")

		t.ExpectPopup().Textbox().Title(Equals("description")).Type("hello").NewLine().Type("world~!").Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Are you sure?")).
			Content(Equals("Are you REALLY sure you want to make this file? Up to you buddy.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("output.txt").IsSelected(),
			)

		t.Views().Main().Content(Contains("+hello\n+world~!"))
	},
})
