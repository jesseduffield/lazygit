package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordLastCommitOfStackedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords the last commit of a branch in the middle of a stack",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"master"}
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(2, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 05").IsSelected(),
				Contains("CI commit 04"),
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 03")).
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Reword commit")).
					InitialText(Equals("commit 03")).
					Clear().
					Type("renamed 03").
					Confirm()
			}).
			Lines(
				Contains("CI commit 05"),
				Contains("CI commit 04"),
				Contains("CI * renamed 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			)
	},
})
