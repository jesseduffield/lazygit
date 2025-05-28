package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetToDuplicateNamedUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset the current branch to an upstream branch when there is a competing tag name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			NewBranch("foo").
			EmptyCommit("commit 1").
			PushBranchAndSetUpstream("origin", "foo").
			EmptyCommit("commit 2").
			CreateLightweightTag("origin/foo", "HEAD")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("commit 2"),
			Contains("commit 1"),
		)
		t.Views().Tags().Focus().Lines(Contains("origin/foo"))

		t.Views().Remotes().Focus().
			Lines(Contains("origin")).
			PressEnter()
		t.Views().RemoteBranches().IsFocused().
			Lines(Contains("foo")).
			Press(keys.Commits.ViewResetOptions)
		t.ExpectPopup().Menu().
			Title(Contains("Reset to origin/foo")).
			Select(Contains("Hard reset")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("commit 1"),
		)

		t.Views().Tags().Focus().
			Lines(Contains("origin/foo")).
			Press(keys.Commits.ViewResetOptions)
		t.ExpectPopup().Menu().
			Title(Contains("Reset to origin/foo")).
			Select(Contains("Hard reset")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("commit 2"),
			Contains("commit 1"),
		)
	},
})
