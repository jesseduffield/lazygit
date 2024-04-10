package branch

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var DetachedHead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch on detached head",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			Checkout("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				MatchesRegexp(`\*.*HEAD`).IsSelected(),
				MatchesRegexp(`master`),
			).
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(MatchesRegexp(`^New branch name \(branch is off of '[0-9a-f]+'\)$`)).
			Type("new-branch").
			Confirm()

		t.Views().Branches().
			Lines(
				MatchesRegexp(`\* new-branch`).IsSelected(),
				MatchesRegexp(`master`),
			)

		t.Git().CurrentBranchName("new-branch")
	},
})
