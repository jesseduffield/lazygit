package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RunCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command that uses runCommand template function in a prompt step",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `git checkout {{.Form.Branch}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:          "Branch",
						Type:         "input",
						Title:        "Enter a branch name",
						InitialValue: "myprefix/{{ runCommand \"echo dynamic\" }}/",
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press("a")

		t.ExpectPopup().Prompt().
			Title(Equals("Enter a branch name")).
			InitialText(Contains("myprefix/dynamic/")).
			Confirm()
	},
})
