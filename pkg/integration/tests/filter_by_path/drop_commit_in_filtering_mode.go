package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropCommitInFilteringMode = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, then drop a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		filterByFilterFile(t, keys)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains(`both files`).IsSelected(),
				Contains(`only filterFile`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Drop commit or delete branch")).
					Select(Contains("Drop commit")).
					Confirm()
			}).
			Lines(
				Contains(`only filterFile`).IsSelected(),
			).
			Press(keys.Universal.Return).
			Lines(
				Contains(`none of the two`),
				Contains(`only otherFile`),
				Contains(`only filterFile`).IsSelected(),
			)
	},
})
