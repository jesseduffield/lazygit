package shell_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteSuggestionWhenNoneSelected = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Press the delete key again after deleting the last entry in the shell command history",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo:    func(shell *Shell) {},
	SetupConfig:  func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			Type("echo 1").
			Confirm()

		// Deleting the only suggestion leaves the list focused but empty, so
		// there is nothing for a second delete to act on.
		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			SuggestionLines(Contains("echo 1")).
			DeleteSuggestion(Contains("echo 1"))

		t.Views().Suggestions().
			IsFocused().
			IsEmpty().
			Press(keys.Universal.Remove)

		t.ExpectToast(Equals("Disabled: No item selected"))
	},
})
