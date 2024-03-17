package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveUpdateRefTodo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move an update-ref item in the rebase todo list",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(3).
			NewBranch("branch2").
			CreateNCommitsStartingAt(3, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("CI commit 06"),
				Contains("pick").Contains("CI commit 05"),
				Contains("pick").Contains("CI commit 04"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit 03"),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI ◯ <-- YOU ARE HERE --- commit 01"),
			).
			NavigateToLine(Contains("update-ref")).
			Press(keys.Commits.MoveUpCommit).
			Press(keys.Commits.MoveUpCommit).
			Lines(
				Contains("pick").Contains("CI commit 06"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit 05"),
				Contains("pick").Contains("CI commit 04"),
				Contains("pick").Contains("CI commit 03"),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI ◯ <-- YOU ARE HERE --- commit 01"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("CI ◯ commit 06"),
				Contains("CI ◯ * commit 05"),
				Contains("CI ◯ commit 04"),
				Contains("CI ◯ commit 03"),
				Contains("CI ◯ commit 02"),
				Contains("CI ◯ commit 01"),
			)
	},
})
