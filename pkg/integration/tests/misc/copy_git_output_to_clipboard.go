package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyGitOutputToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy streamed git output from the command log to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.Views().Status().Content(Equals("✓ repo → master"))

		t.GlobalPress(keys.Universal.ExtrasMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Command log")).
			Select(Contains("Copy last git output to clipboard")).
			Confirm()

		t.ExpectToast(Equals("Git output copied to clipboard"))

		t.FileSystem().FileContent("clipboard",
			Contains("master -> master").
				Contains("git push").
				Contains("Push"))
	},
})
