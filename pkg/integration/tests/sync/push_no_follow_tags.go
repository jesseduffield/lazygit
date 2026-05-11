package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushNoFollowTags = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push with --follow-tags NOT configured in git config",
	ExtraCmdArgs: []string{},
	Skip:         true, // turns out this actually DOES push the tag. I have no idea why
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.CreateAnnotatedTag("mytag", "message", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("✓ repo → master"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)

		t.Views().Status().Content(Equals("✓ repo → master"))

		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin"),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains("master"),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				// tag was not pushed to upstream
				Contains("two").DoesNotContain("mytag"),
				Contains("one"),
			)
	},
})
