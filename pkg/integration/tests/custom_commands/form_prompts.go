package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FormPrompts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command referring prompt responses by name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press("a")

		t.ExpectPopup().Prompt().Title(Equals("Enter a file name")).Type("my file").Confirm()

		t.ExpectPopup().Menu().Title(Equals("Choose file content")).Select(Contains("bar")).Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Are you sure?")).
			Content(Equals("Are you REALLY sure you want to make this file? Up to you buddy.")).
			Confirm()

		t.Views().Files().
			Lines(
				Contains("my file").IsSelected(),
			)

		t.Views().Main().Content(Contains(`"BAR"`))
	},
})
