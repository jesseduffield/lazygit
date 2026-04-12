package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConditionalPromptFalseString = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Conditional prompt is skipped when condition is bare false or template false",
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
				Command: `echo "{{.Form.Choice}}" > result.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Choice",
						Type:  "menu",
						Title: "Pick one",
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
						},
					},
					{
						Key:       "Skipped1",
						Type:      "input",
						Title:     "This is always skipped (false)",
						Condition: `false`,
					},
					{
						Key:       "Skipped2",
						Type:      "input",
						Title:     "This is always skipped (template false)",
						Condition: `{{ eq "a" "b" }}`,
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press("a")

		t.ExpectPopup().Menu().Title(Equals("Pick one")).Select(Contains("foo")).Confirm()

		// Both conditional prompts skipped, file created directly
		t.Views().Files().
			Focus().
			Lines(
				Contains("result.txt").IsSelected(),
			)

		t.FileSystem().FileContent("result.txt", Equals("FOO\n"))
	},
})
