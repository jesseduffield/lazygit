package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateTag = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new tag on a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press(keys.Commits.CreateTag)

		t.ExpectPopup().Menu().
			Title(Equals("Create tag")).
			Select(Contains("Lightweight")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Tag name:")).
			Type("new-tag").
			Confirm()

		t.Views().Commits().
			Lines(
				MatchesRegexp(`new-tag.*two`).IsSelected(),
				MatchesRegexp(`one`),
			)

		t.Views().Tags().
			Focus().
			Lines(
				MatchesRegexp(`new-tag.*two`).IsSelected(),
			)

		t.Git().
			TagNamesAt("HEAD", []string{"new-tag"})
	},
})
