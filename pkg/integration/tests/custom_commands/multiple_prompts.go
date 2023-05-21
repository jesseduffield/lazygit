package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MultiplePrompts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command with multiple prompts",
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press("a")

		t.ExpectPopup().Prompt().Title(Equals("Enter a file name")).Type("myfile").Confirm()

		t.ExpectPopup().Menu().Title(Equals("Choose file content")).Select(Contains("bar")).Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Are you sure?")).
			Content(Equals("Are you REALLY sure you want to make this file? Up to you buddy.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("myfile").IsSelected(),
			)

		t.Views().Main().Content(Contains("BAR"))
	},
})
