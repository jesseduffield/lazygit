package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragAcrossBranchBoundaries = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drop a dragged commit between the update-ref rows of a stack of branches",
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
			NewBranch("branch1b").
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit-05").IsSelected(),
				Contains("CI commit-04"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			// Once the drag starts, the update-ref rows appear
			ClickAndHold(1, 0).
			MouseMove(1, 1).
			Lines(
				Contains("CI commit-05").IsSelected(),
				Contains("CI commit-04"),
				Contains("drop here"),
				MatchesRegexp(`─ branch1\s*$`),
				Contains("─ branch1b"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			MouseMove(1, 3).
			Lines(
				Contains("CI commit-05").IsSelected(),
				Contains("CI commit-04"),
				MatchesRegexp(`─ branch1\s*$`),
				Contains("drop here"),
				Contains("─ branch1b"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			MouseRelease().
			Lines(
				Contains("CI commit-04"),
				Contains("CI * commit-05").IsSelected(),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)

		t.Git().
			BranchNamesAt("HEAD~1", []string{"branch1"}).
			BranchNamesAt("HEAD~2", []string{"branch1b"})

		t.Views().Commits().
			// Moving the pointer sideways doesn't show the update-ref rows yet
			ClickAndHold(1, 1).
			MouseMove(5, 1).
			Lines(
				Contains("CI commit-04"),
				Contains("CI * commit-05").IsSelected(),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			// Once the pointer is on another row, the update-ref row of branch1
			// appears above the pressed commit, which moves down by a row
			MouseMove(1, 0).
			Lines(
				Contains("drop here"),
				Contains("CI commit-04"),
				MatchesRegexp(`─ branch1\s*$`),
				Contains("CI * commit-05").IsSelected(),
				Contains("─ branch1b"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			MouseMove(1, 2).
			Lines(
				Contains("CI commit-04"),
				Contains("drop here"),
				MatchesRegexp(`─ branch1\s*$`),
				Contains("CI * commit-05").IsSelected(),
				Contains("─ branch1b"),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			).
			MouseRelease().
			Lines(
				Contains("CI commit-04"),
				Contains("CI commit-05").IsSelected(),
				Contains("CI * commit-03"),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)

		t.Git().
			BranchNamesAt("HEAD~2", []string{"branch1", "branch1b"})
	},
})
