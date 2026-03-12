package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConditionalPrompts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command with conditional prompts that are skipped based on form values",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "files",
				Command: `echo "{{.Form.Choice}}{{if .Form.Detail}} {{.Form.Detail}}{{end}}" > result.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Key:   "Choice",
						Type:  "menu",
						Title: "Choose an option",
						Options: []config.CustomCommandMenuOption{
							{
								Name:        "first",
								Description: "First option",
								Value:       "FIRST",
								Key:         "1",
							},
							{
								Name:        "second",
								Description: "Second option",
								Value:       "SECOND",
								Key:         "H",
							},
						},
					},
					{
						Key:       "Detail",
						Type:      "input",
						Title:     "Enter detail for second option",
						Condition: `{{ eq .Form.Choice "SECOND" }}`,
					},
				},
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Test 1: Select "first" via key — conditional prompt should be skipped
		t.Views().Files().
			IsFocused().
			Press("a")

		t.ExpectPopup().Menu().
			Title(Equals("Choose an option"))

		t.Views().Menu().Press("1")

		// Detail prompt should be skipped, file should be created directly
		t.Views().Files().
			Focus().
			Lines(
				Contains("result.txt").IsSelected(),
			)

		t.FileSystem().FileContent("result.txt", Equals("FIRST\n"))

		// Test 2: Select "second" via key — conditional prompt should appear
		t.Shell().DeleteFile("result.txt")
		t.GlobalPress(keys.Files.RefreshFiles)

		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press("a")

		t.ExpectPopup().Menu().
			Title(Equals("Choose an option"))

		t.Views().Menu().Press("H")

		// Detail prompt should appear because Choice == "SECOND"
		t.ExpectPopup().Prompt().Title(Equals("Enter detail for second option")).Type("extra").Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Contains("result.txt").IsSelected(),
			)

		t.FileSystem().FileContent("result.txt", Equals("SECOND extra\n"))
	},
})
