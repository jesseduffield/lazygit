package reflog

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DoNotShowBranchMarkersInReflogSubcommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that no branch heads are shown in the subcommits view of a reflog entry",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.AppState.GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch1")
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.NewBranch("branch2")
		shell.EmptyCommit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Check that the local commits view does show a branch marker for branch1
		t.Views().Commits().
			Lines(
				Contains("CI three"),
				Contains("CI * two"),
				Contains("CI one"),
			)

		t.Views().Branches().
			Focus().
			// Check out branch1
			NavigateToLine(Contains("branch1")).
			PressPrimaryAction().
			// Look at the subcommits of branch2
			NavigateToLine(Contains("branch2")).
			PressEnter().
			// Check that we see a marker for branch1 here (but not for
			// branch2), even though branch1 is checked out
			Tap(func() {
				t.Views().SubCommits().
					IsFocused().
					Lines(
						Contains("CI three"),
						Contains("CI * two"),
						Contains("CI one"),
					).
					PressEscape()
			}).
			// Check out branch2 again
			NavigateToLine(Contains("branch2")).
			PressPrimaryAction()

		t.Views().ReflogCommits().
			Focus().
			TopLines(
				Contains("checkout: moving from branch1 to branch2").IsSelected(),
			).
			PressEnter().
			// Check that the subcommits view for a reflog entry doesn't show
			// any branch markers
			Tap(func() {
				t.Views().SubCommits().
					IsFocused().
					Lines(
						Contains("CI three"),
						Contains("CI two"),
						Contains("CI one"),
					)
			})
	},
})
