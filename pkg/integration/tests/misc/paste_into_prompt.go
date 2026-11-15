package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteIntoPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste multi-line text into a prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Contains("New branch name")).
			Paste("\nmy\tnew\nbranch\n")

		t.Views().Prompt().
			IsFocused().
			Content(Equals("my new branch"))
	},
})
