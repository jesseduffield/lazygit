package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropTodoCommitWithUpdateRefShowBranchHeads = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit during interactive rebase when there is an update-ref in the git-rebase-todo file (with experimentalShowBranchHeads on)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.ExperimentalShowBranchHeads = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			NewBranch("mybranch").
			CreateNCommitsStartingAt(3, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("(*) commit 06").IsSelected(),
				Contains("commit 05"),
				Contains("commit 04"),
				Contains("(*) commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Focus().
			Lines(
				Contains("pick").Contains("(*) commit 06"),
				Contains("pick").Contains("commit 05"),
				Contains("pick").Contains("commit 04"),
				Contains("update-ref").Contains("master"),
				Contains("pick").Contains("(*) commit 03"),
				Contains("pick").Contains("commit 02"),
				Contains("<-- YOU ARE HERE --- commit 01"),
			).
			NavigateToLine(Contains("commit 05")).
			Press(keys.Universal.Remove)

		t.Common().ContinueRebase()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("(*) commit 06"),
				Contains("commit 04"),
				Contains("(*) commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
