package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushWithUnknownHostPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a commit to a pre-configured upstream, where the SSH host must be verified",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")

		// simulate pushing to an unknown host by using a pre-push hook that prompts for host verification.
		shell.CopyHelpFile("pre-push-unknown-host", ".git/hooks/pre-push")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.ExpectPopup().Prompt().
			Title(Equals("SSH host verification (type 'yes', 'no', or fingerprint)")).
			Type("no").
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Host key verification failed")).
			Confirm()

		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.ExpectPopup().Prompt().
			Title(Equals("SSH host verification (type 'yes', 'no', or fingerprint)")).
			Type("yes").
			Confirm()

		t.Views().Status().Content(Equals("✓ repo → master"))

		assertSuccessfullyPushed(t)
	},
})
