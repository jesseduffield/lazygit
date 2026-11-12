package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyToClipboardRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the names of a range of selected branches to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.EmptyCommit("initial commit")
		shell.NewBranch("branch-b")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-b").IsSelected(),
				Contains("branch-a"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.CopyToClipboard)

		t.ExpectToast(Equals("2 items copied to clipboard"))

		t.FileSystem().FileContent("clipboard", Equals("branch-b\nbranch-a"))
	},
})
