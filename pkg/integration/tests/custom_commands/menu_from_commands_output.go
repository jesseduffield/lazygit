package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MenuFromCommandsOutput = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using prompt response in menuFromCommand entries",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("foo").
			NewBranch("feature/foo").
			EmptyCommit("bar").
			NewBranch("feature/bar").
			EmptyCommit("baz")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: "git checkout {{ index .PromptResponses 1 }}",
				Prompts: []config.CustomCommandPrompt{
					{
						Type:         "input",
						Title:        "Which git command do you want to run?",
						InitialValue: "branch",
					},
					{
						Type:        "menuFromCommand",
						Title:       "Branch:",
						Command:     `git {{ index .PromptResponses 0 }} --format='%(refname:short)'`,
						Filter:      "(?P<branch>.*)",
						ValueFormat: `{{ .branch }}`,
						LabelFormat: `{{ .branch | green }}`,
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Git().CurrentBranchName("feature/bar")

		t.Views().Branches().
			Focus().
			Press("a")

		t.ExpectPopup().Prompt().
			Title(Equals("Which git command do you want to run?")).
			InitialText(Equals("branch")).
			Confirm()

		t.ExpectPopup().Menu().Title(Equals("Branch:")).Select(Equals("master")).Confirm()

		t.Git().CurrentBranchName("master")
	},
})
