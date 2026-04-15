package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConditionalPromptFalseValue = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Entering literal false as form input does not incorrectly skip a conditional prompt",
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
				Command: `echo "{{.Form.Word}} {{.Form.Extra}}" > result.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Word",
						Type:  "input",
						Title: "Enter a word",
					},
					{
						Key:       "Extra",
						Type:      "input",
						Title:     "Enter extra",
						Condition: `{{ eq .Form.Word "false" }}`,
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press("a")

		t.ExpectPopup().Prompt().Title(Equals("Enter a word")).Type("false").Confirm()

		// Condition {{ eq .Form.Word "false" }} evaluates to true, so prompt should appear
		t.ExpectPopup().Prompt().Title(Equals("Enter extra")).Type("baz").Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("result.txt").IsSelected(),
			)

		t.FileSystem().FileContent("result.txt", Equals("false baz\n"))
	},
})
