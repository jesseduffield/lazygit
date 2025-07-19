package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordCommitInFilteringMode = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, then reword a commit",
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
			SelectNextItem().
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Clear().
					Type("new message").
					Confirm()
			}).
			Lines(
				Contains(`both files`),
				Contains(`new message`).IsSelected(),
			).
			Press(keys.Universal.Return).
			Lines(
				Contains(`none of the two`),
				Contains(`both files`),
				Contains(`only otherFile`),
				Contains(`new message`).IsSelected(),
			)
	},
})
