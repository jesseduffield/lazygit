package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropTodoCommitWithUpdateRef = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit during interactive rebase when there is an update-ref in the git-rebase-todo file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(3, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(3, 5)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 07").IsSelected(),
				Contains("commit 06"),
				Contains("commit 05"),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Edit).
			Focus().
			Lines(
				Contains("pick").Contains("commit 07"),
				Contains("pick").Contains("commit 06"),
				Contains("pick").Contains("commit 05"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("commit 04"),
				Contains("pick").Contains("commit 03"),
				Contains("<-- YOU ARE HERE --- commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 06")).
			Press(keys.Universal.Remove)

		t.Common().ContinueRebase()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 07"),
				Contains("commit 05"),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
