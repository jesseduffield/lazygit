package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SetAuthorRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set author on a range of commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("user.email", "Bill@example.com")
		shell.SetConfig("user.name", "Bill Smith")

		shell.EmptyCommit("fourth")
		shell.EmptyCommit("third")
		shell.EmptyCommit("second")
		shell.EmptyCommit("first")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("BS").Contains("first").IsSelected(),
				Contains("BS").Contains("second"),
				Contains("BS").Contains("third"),
				Contains("BS").Contains("fourth"),
			)

		t.Views().Commits().
			Focus().
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains(" Set author")). // adding space at start to distinguish from 'reset author'
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Set author")).
					Type("John Smith <John@example.com>").
					Confirm()
			}).
			PressEscape().
			Lines(
				Contains("BS").Contains("first"),
				Contains("JS").Contains("second"),
				Contains("JS").Contains("third").IsSelected(),
				Contains("BS").Contains("fourth"),
			)
	},
})
