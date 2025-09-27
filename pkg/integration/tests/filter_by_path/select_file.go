package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, by finding file in UI and filtering on it",
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
				Contains(`both files`),
				Contains(`only otherFile`),
				Contains(`only filterFile`),
			).
			NavigateToLine(Contains(`both files`)).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals(`▼ /`).IsSelected(),
				Equals(`  M filterFile`),
				Equals(`  M otherFile`),
			).
			SelectNextItem().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().Title(Equals("Filtering")).
			Select(Contains("Filter by 'filterFile'")).Confirm()

		postFilterTest(t)
	},
})
