package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateTag = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new tag on branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			NewBranch("new-branch").
			EmptyCommit("new commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				MatchesRegexp(`\*\s*new-branch`).IsSelected(),
				MatchesRegexp(`master`),
			).
			SelectNextItem().
			Press(keys.Branches.CreateTag)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Tag name")).
			Type("new-tag").
			Confirm()

		t.Views().Tags().Focus().
			Lines(
				MatchesRegexp(`new-tag`).IsSelected(),
			)

		t.Git().
			TagNamesAt("HEAD", []string{}).
			TagNamesAt("master", []string{"new-tag"})
	},
})
