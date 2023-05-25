package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Reset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset to a tag",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.CreateLightweightTag("tag", "HEAD^") // creating tag on commit "one"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Lines(
			Contains("two"),
			Contains("one"),
		)

		t.Views().Tags().
			Focus().
			Lines(
				Contains("tag").IsSelected(),
			).
			Press(keys.Commits.ViewResetOptions)

		t.ExpectPopup().Menu().
			Title(Contains("Reset to tag")).
			Select(Contains("Hard reset")).
			Confirm()

		t.Views().Commits().Lines(
			Contains("one"),
		)
	},
})
