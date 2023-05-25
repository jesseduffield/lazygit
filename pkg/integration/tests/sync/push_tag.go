package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushTag = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a specific tag",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")

		shell.CloneIntoRemote("origin")

		shell.CreateAnnotatedTag("mytag", "message", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("mytag"),
			).
			Press(keys.Branches.PushTag)

		t.ExpectPopup().Prompt().
			Title(Equals("Remote to push tag 'mytag' to:")).
			InitialText(Equals("origin")).
			SuggestionLines(
				Contains("origin"),
			).
			Confirm()

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
				Contains("two").Contains("mytag"),
				Contains("one"),
			)
	},
})
