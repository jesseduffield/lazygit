package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// NOTE: we're getting a weird offset in the popup prompt for some reason. Not sure what's behind that.

var MenuFromCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using menuFromCommand prompt type",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("foo").
			EmptyCommit("bar").
			EmptyCommit("baz").
			NewBranch("feature/foo")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `echo "{{index .PromptResponses 0}} {{index .PromptResponses 1}} {{ .SelectedLocalBranch.Name }}" > output.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Type:        "menuFromCommand",
						Title:       "Choose commit message",
						Command:     `git log --oneline --pretty=%B`,
						Filter:      `(?P<commit_message>.*)`,
						ValueFormat: `{{ .commit_message }}`,
						LabelFormat: `{{ .commit_message | yellow }}`,
					},
					{
						Type:         "input",
						Title:        "Description",
						InitialValue: `{{ if .SelectedLocalBranch.Name }}Branch: #{{ .SelectedLocalBranch.Name }}{{end}}`,
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty()

		t.Views().Branches().
			Focus().
			Press("a")

		t.ExpectPopup().Menu().Title(Equals("Choose commit message")).Select(Contains("bar")).Confirm()

		t.ExpectPopup().Prompt().Title(Equals("Description")).Type(" my branch").Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("output.txt").IsSelected(),
			)

		t.Views().Main().Content(Contains("bar Branch: #feature/foo my branch feature/foo"))
	},
})
