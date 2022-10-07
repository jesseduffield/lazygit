package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FormPrompts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command reffering prompt responses by name",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: `echo {{.Form.FileContent | quote}} > {{.Form.FileName | quote}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "FileName",
						Type:  "input",
						Title: "Enter a file name",
					},
					{
						Key:   "FileContent",
						Type:  "menu",
						Title: "Choose file content",
						Options: []config.CustomCommandMenuOption{
							{
								Name:        "foo",
								Description: "Foo",
								Value:       "FOO",
							},
							{
								Name:        "bar",
								Description: "Bar",
								Value:       `"BAR"`,
							},
							{
								Name:        "baz",
								Description: "Baz",
								Value:       "BAZ",
							},
						},
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
	Run: func(
		shell *Shell,
		input *Input,
		assert *Assert,
		keys config.KeybindingConfig,
	) {
		assert.WorkingTreeFileCount(0)

		input.PressKeys("a")

		assert.InPrompt()
		assert.MatchCurrentViewTitle(Equals("Enter a file name"))
		input.Type("my file")
		input.Confirm()

		assert.InMenu()
		assert.MatchCurrentViewTitle(Equals("Choose file content"))
		assert.MatchSelectedLine(Contains("foo"))
		input.NextItem()
		assert.MatchSelectedLine(Contains("bar"))
		input.Confirm()

		assert.InConfirm()
		assert.MatchCurrentViewTitle(Equals("Are you sure?"))
		input.Confirm()

		assert.WorkingTreeFileCount(1)
		assert.MatchSelectedLine(Contains("my file"))
		assert.MatchMainViewContent(Contains(`"BAR"`))
	},
})
