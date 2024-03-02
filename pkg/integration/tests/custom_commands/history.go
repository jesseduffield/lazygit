package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var History = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test that the custom commands history is saved correctly",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig:  func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type("echo 1").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(Contains("1")).
			Type("echo 2").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(
				// "echo 2" was typed last, so it should come first
				Contains("2"),
				Contains("1"),
			).
			Type("echo 3").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(
				Contains("3"),
				Contains("2"),
				Contains("1"),
			).
			Type("echo 1").
			Confirm()

		// Executing a command again should move it to the front:
		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(
				Contains("1"),
				Contains("3"),
				Contains("2"),
			)
	},
})
