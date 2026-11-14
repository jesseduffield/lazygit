package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GotoBottomInBranchStack = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Go to the bottom of the commits list in steps, stopping at the beginning of each branch of a stack",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2).
			NewBranch("branch1").
			CreateNCommitsStartingAt(2, 3).
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 5)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ○ commit-06").IsSelected(),
				Contains("CI ○ commit-05"),
				Contains("CI ○ * commit-04"),
				Contains("CI ○ commit-03"),
				Contains("CI ○ commit-02"),
				Contains("CI ○ commit-01"),
			).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-05")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-03")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-01"))

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch2")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("commit-06")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-05")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-03")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-01"))

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit-03")).
			Press(keys.Universal.Edit).
			NavigateToLine(Contains("commit-06")).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("CI commit-06").IsSelected(),
				Contains("pick").Contains("CI commit-05"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit-04"),
				Contains("─── Commits"),
				Contains("CI ○ commit-03"),
				Contains("CI ○ commit-02"),
				Contains("CI ○ commit-01"),
			).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-05")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-04")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-03")).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit-01"))
	},
})
