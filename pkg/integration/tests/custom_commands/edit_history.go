package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditHistory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit an entry from the custom commands history",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig:  func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type("echo x").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type("ec").
			SuggestionLines(
				Equals("echo x"),
			).
			EditSuggestion(Equals("echo x")).
			InitialText(Equals("echo x"))
	},
})
