package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenCommandLogInEditor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open the full command log in the user's editor",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.Edit = "cp {{filename}} editor-output.txt"
		config.GetUserConfig().Gui.ShowRandomTip = false
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
			Select(Contains("Open command log in editor")).
			Confirm()

		t.ExpectToast(Equals("Command log opened in editor"))

		t.FileSystem().FileContent("editor-output.txt",
			Contains("Push").
				Contains("git push").
				Contains("master -> master").
				DoesNotContain("You can hide/focus"))
	},
})
