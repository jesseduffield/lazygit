package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MultiplePrompts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command with multiple prompts",
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
				Command: `echo "{{index .PromptResponses 1}}" > {{index .PromptResponses 0}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Type:  "input",
						Title: "Enter a file name",
					},
					{
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
								Value:       "BAR",
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
		input.Type("myfile")
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
		assert.MatchSelectedLine(Contains("myfile"))
		assert.MatchMainViewContent(Contains("BAR"))
	},
})
