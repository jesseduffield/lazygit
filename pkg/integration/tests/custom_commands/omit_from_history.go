package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OmitFromHistory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Omitting a runtime custom command from history if it begins with space",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type("echo aubergine").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(Contains("aubergine")).
			SuggestionLines(DoesNotContain("tangerine")).
			Type(" echo tangerine").
			Confirm()

		t.GlobalPress(keys.Universal.ExecuteCustomCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			SuggestionLines(Contains("aubergine")).
			SuggestionLines(DoesNotContain("tangerine")).
			Cancel()
	},
})
