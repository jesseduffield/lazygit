package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Push = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a commit to a pre-configured upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("↑1 repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		assertSuccessfullyPushed(t)
	},
})
