package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MenuFromCommandsOutput = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using prompt response in menuFromCommand entries",
	ExtraCmdArgs: "",
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
	Run: func(
		shell *Shell,
		input *Input,
		assert *Assert,
		keys config.KeybindingConfig,
	) {
		assert.CurrentBranchName("feature/bar")

		assert.WorkingTreeFileCount(0)
		input.SwitchToBranchesWindow()

		input.Press("a")

		assert.InPrompt()
		assert.CurrentViewTitle(Equals("Which git command do you want to run?"))
		assert.SelectedLine(Equals("branch"))
		input.Confirm()

		input.Menu(Equals("Branch:"), Equals("master"))

		assert.CurrentBranchName("master")
	},
})
