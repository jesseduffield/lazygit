package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the tag to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateLightweightTag("super.l000ongtag", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("tag").IsSelected(),
			).
			Press(keys.Universal.CopyToClipboard)

		t.ExpectToast(Equals("'super.l000ongtag' copied to clipboard"))

		t.FileSystem().FileContent("clipboard", Equals("super.l000ongtag"))
	},
})
