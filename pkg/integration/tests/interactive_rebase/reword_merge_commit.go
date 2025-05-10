package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordMergeCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords a merge commit which is not the current head commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("first-branch").
			CreateFileAndAdd("file1.txt", "content").
			Commit("one").
			Checkout("master").
			Merge("first-branch").
			NewBranch("second-branch").
			EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ two").IsSelected(),
				Contains("CI ⏣─╮ Merge branch 'first-branch'"),
				Contains("CI │ ◯ one"),
				Contains("CI ◯─╯ base"),
			).
			SelectNextItem().
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Reword commit")).
					InitialText(Equals("Merge branch 'first-branch'")).
					Clear().
					Type("renamed merge").
					Confirm()
			}).
			Lines(
				Contains("CI ◯ two"),
				Contains("CI ⏣─╮ renamed merge").IsSelected(),
				Contains("CI │ ◯ one"),
				Contains("CI ◯ ╯ base"),
			)
	},
})
