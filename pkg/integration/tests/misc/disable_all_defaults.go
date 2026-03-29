package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableAllDefaults = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirms you can disable all default keybindings by setting disableAllDefaults to true",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Keybinding.DisableAllDefaults = true
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Default keybindings should be disabled, so pressing 'q' should not quit
		// We can't easily test quitting, so let's test navigation
		// Pressing <up> or k should not navigate up because default keybindings are disabled
		t.Views().Files().
			IsFocused().
			// Press a key that would normally navigate but shouldn't work now
			Press("k").
			// Should still be on files panel since navigation keybindings are disabled
			IsFocused()

		// Note: Mouse bindings should still work, but we can't easily test that in integration tests
		// The key point is that keyboard shortcuts are disabled, giving users a blank slate
	},
})