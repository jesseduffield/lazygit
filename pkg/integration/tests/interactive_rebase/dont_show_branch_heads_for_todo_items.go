package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DontShowBranchHeadsForTodoItems = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check that branch heads are shown for normal commits during interactive rebase, but not for todo items",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetAppState().GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(4, 3).
			NewBranch("branch3").
			CreateNCommitsStartingAt(3, 7)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 09"),
				Contains("CI commit 08"),
				Contains("CI commit 07"),
				Contains("CI * commit 06"),
				Contains("CI commit 05"),
				Contains("CI commit 04"),
				Contains("CI commit 03"),
				Contains("CI * commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 04")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("CI commit 09"),
				Contains("pick").Contains("CI commit 08"),
				Contains("pick").Contains("CI commit 07"),
				Contains("update-ref").Contains("branch2"),
				Contains("pick").Contains("CI commit 06"), // no star on this entry, even though branch2 points to it
				Contains("pick").Contains("CI commit 05"),
				Contains("CI <-- YOU ARE HERE --- commit 04"),
				Contains("CI commit 03"),
				Contains("CI * commit 02"), // this star is fine though
				Contains("CI commit 01"),
			)
	},
})
