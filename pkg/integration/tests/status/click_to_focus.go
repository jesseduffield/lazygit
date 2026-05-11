package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickToFocus = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click in the status side panel to activate it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus()
		t.Views().Main().Lines(
			Contains("No changed files"),
		)

		t.Views().Status().Click(0, 0)
		t.Views().Status().IsFocused()
		t.Views().Main().ContainsLines(
			Contains(`   _`),
			Contains(`  | |                     (_) |`),
			Contains(`  | | __ _ _____   _  __ _ _| |_`),
			Contains("  | |/ _` |_  / | | |/ _` | | __|"),
			Contains(`  | | (_| |/ /| |_| | (_| | | |_`),
			Contains(`  |_|\__,_/___|\__, |\__, |_|\__|`),
			Contains(`                __/ | __/ |`),
			Contains(`               |___/ |___/`),
			Contains(``),
			Contains(`Copyright `),
		)
	},
})
