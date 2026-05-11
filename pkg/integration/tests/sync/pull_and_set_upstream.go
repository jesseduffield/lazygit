package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullAndSetUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull a commit from the remote, setting the upstream branch in the process",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")

		// remove the 'two' commit so that we have something to pull from the remote
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("one"),
			)

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Files().IsFocused().Press(keys.Universal.Pull)

		t.ExpectPopup().Prompt().
			Title(Equals("Enter upstream as '<remote> <branchname>'")).
			SuggestionLines(Equals("origin master")).
			ConfirmFirstSuggestion()

		t.Views().Commits().
			Lines(
				Contains("two"),
				Contains("one"),
			)

		t.Views().Status().Content(Equals("✓ repo → master"))
	},
})
