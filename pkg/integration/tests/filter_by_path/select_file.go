package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, by finding file in UI and filtering on it",
	ExtraCmdArgs: "",
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
				Contains(`only filterFile`).IsSelected(),
				Contains(`only otherFile`),
				Contains(`both files`),
			).
			PressEnter()

		// when you click into the commit itself, you see all files from that commit
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains(`filterFile`).IsSelected(),
			).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().Title(Equals("Filtering")).Select(Contains("filter by 'filterFile'")).Confirm()

		postFilterTest(t)
	},
})
