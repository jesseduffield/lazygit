package shell_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteFromHistory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete an entry from the custom commands history",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig:  func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		createCustomCommand := func(command string) {
			t.GlobalPress(keys.Universal.ExecuteShellCommand)
			t.ExpectPopup().Prompt().
				Title(Equals("Shell command:")).
				Type(command).
				Confirm()
		}

		createCustomCommand("echo 1")
		createCustomCommand("echo 2")
		createCustomCommand("echo 3")

		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			SuggestionLines(
				Contains("3"),
				Contains("2"),
				Contains("1"),
			).
			DeleteSuggestion(Contains("2")).
			SuggestionLines(
				Contains("3"),
				Contains("1"),
			)
	},
})
