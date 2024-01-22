package filter_by_author

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var TypeAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by author using the typed in author",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter author to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter author:")).
			Type("Yang").
			SuggestionLines(Equals("Yang Wen-li <yang.wen-li@email.com>")).
			ConfirmFirstSuggestion()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 2"),
				Contains("commit 1"),
				Contains("commit 0"),
			)

		t.Views().Information().Content(Contains("Filtering by 'Yang Wen-li <yang.wen-li@email.com>'"))

		t.Views().Status().
			Focus().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter author to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter author:")).
			Type("Siegfried").
			SuggestionLines(Equals("Siegfried Kircheis <siegfried.kircheis@email.com>")).
			ConfirmFirstSuggestion()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 0"),
			)

		t.Views().Information().Content(Contains("Filtering by 'Siegfried Kircheis <siegfried.kircheis@email.com>'"))
	},
})
