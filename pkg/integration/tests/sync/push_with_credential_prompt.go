package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushWithCredentialPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a commit to a pre-configured upstream, where credentials are required",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")

		// actually getting a password prompt is tricky: it requires SSH'ing into localhost under a newly created, restricted, user.
		// This is not easy to do in a cross-platform way, nor is it easy to do in a docker container.
		// If you can think of a way to do it, please let me know!
		shell.CopyHelpFile("pre-push", ".git/hooks/pre-push")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		// correct credentials are: username=username, password=password

		t.ExpectPopup().Prompt().
			Title(Equals("Username")).
			Type("username").
			Confirm()

		// enter incorrect password
		t.ExpectPopup().Prompt().
			Title(Equals("Password")).
			Type("incorrect password").
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("incorrect username/password")).
			Confirm()

		t.Views().Status().Content(Equals("↑1 repo → master"))

		// try again with correct password
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.ExpectPopup().Prompt().
			Title(Equals("Username")).
			Type("username").
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Password")).
			Type("password").
			Confirm()

		t.Views().Status().Content(Equals("✓ repo → master"))

		assertSuccessfullyPushed(t)
	},
})
