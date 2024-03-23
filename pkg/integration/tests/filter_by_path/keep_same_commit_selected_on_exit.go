package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSameCommitSelectedOnExit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When exiting filtering mode, keep the same commit selected if possible",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains(`none of the two`).IsSelected(),
				Contains(`only filterFile`),
				Contains(`only otherFile`),
				Contains(`both files`),
			).Press(keys.Universal.FilteringMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Filtering")).
					Select(Contains("Enter path to filter by")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Enter path:")).
					Type("filterF").
					SuggestionLines(Equals("filterFile")).
					ConfirmFirstSuggestion()
			}).
			Lines(
				Contains(`only filterFile`).IsSelected(),
				Contains(`both files`),
			).
			SelectNextItem().
			PressEscape().
			Lines(
				Contains(`none of the two`),
				Contains(`only filterFile`),
				Contains(`only otherFile`),
				Contains(`both files`).IsSelected(),
			)
	},
})
