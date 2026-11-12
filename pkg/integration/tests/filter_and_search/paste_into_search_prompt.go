package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteIntoSearchPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste multi-line text into the search prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.StartSearch)

		t.ExpectSearch().
			Paste("\ntwo\n")

		t.Views().Search().
			IsFocused().
			Content(Equals("two"))

		t.ExpectSearch().
			Confirm()

		t.Views().Search().
			IsVisible().
			Content(Contains("matches for 'two' (1 of 1)"))

		t.Views().Commits().
			Lines(
				Contains("three"),
				Contains("two").IsSelected(),
				Contains("one"),
			)
	},
})
