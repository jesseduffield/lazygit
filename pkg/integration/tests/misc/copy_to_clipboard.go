package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a branch name to the clipboard using custom clipboard command template",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.OS.CopyToClipboardCmd = "echo {{text}} > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
			).
			Press(keys.Universal.CopyToClipboard)

		t.Views().Files().
			Focus()

		t.GlobalPress(keys.Files.RefreshFiles)

		// Expect to see the clipboard file with contents
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("clipboard").IsSelected(),
			)

		t.Views().Main().Content(Contains("branch-a"))
	},
})
