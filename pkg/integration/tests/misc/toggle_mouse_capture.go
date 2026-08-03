package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleMouseCapture = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle mouse capture off and back on with the global keybinding",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.ToggleMouseCapture)
		t.ExpectToast(Equals("Mouse capture disabled"))

		t.GlobalPress(keys.Universal.ToggleMouseCapture)
		t.ExpectToast(Equals("Mouse capture enabled"))
	},
})
