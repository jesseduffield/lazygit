package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SuggestionsCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command that uses a suggestions command in a prompt step",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-one")
		shell.EmptyCommit("blah")
		shell.NewBranch("branch-two")
		shell.EmptyCommit("blah")
		shell.NewBranch("branch-three")
		shell.EmptyCommit("blah")
		shell.NewBranch("branch-four")
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `git checkout {{.Form.Branch}}`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Branch",
						Type:  "input",
						Title: "Enter a branch name",
						Suggestions: config.CustomCommandSuggestions{
							Command: "git branch --format='%(refname:short)'",
						},
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-four").IsSelected(),
				Contains("branch-three"),
				Contains("branch-two"),
				Contains("branch-one"),
			).
			Press("a")

		t.ExpectPopup().Prompt().
			Title(Equals("Enter a branch name")).
			Type("three").
			SuggestionLines(Contains("branch-three")).
			ConfirmFirstSuggestion()

		t.Views().Branches().
			Lines(
				Contains("branch-three").IsSelected(),
				Contains("branch-four"),
				Contains("branch-two"),
				Contains("branch-one"),
			)
	},
})
