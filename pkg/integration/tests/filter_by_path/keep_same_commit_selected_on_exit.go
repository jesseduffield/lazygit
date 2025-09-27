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
		filterByFilterFile(t, keys)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains(`both files`).IsSelected(),
				Contains(`only filterFile`),
			).
			Tap(func() {
				t.Views().Main().
					ContainsLines(
						Equals("    both files"),
						Equals("---"),
						Equals(" filterFile | 2 +-"),
						Equals(" 1 file changed, 1 insertion(+), 1 deletion(-)"),
					)
			}).
			PressEscape().
			Lines(
				Contains(`none of the two`),
				Contains(`both files`).IsSelected(),
				Contains(`only otherFile`),
				Contains(`only filterFile`),
			)

		t.Views().Main().
			ContainsLines(
				Equals("    both files"),
				Equals("---"),
				Equals(" filterFile | 2 +-"),
				Equals(" otherFile  | 2 +-"),
				Equals(" 2 files changed, 2 insertions(+), 2 deletions(-)"),
			)
	},
})
