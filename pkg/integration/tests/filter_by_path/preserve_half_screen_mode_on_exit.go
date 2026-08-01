package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PreserveHalfScreenModeOnExit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Exiting commit filtering preserves configured half screen mode",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ScreenMode = "half"
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			IsFocused()

		t.Views().Files().
			IsInvisible()

		filterByFilterFile(t, keys)
		postFilterTest(t)

		t.Views().Files().
			IsInvisible()

		t.Views().Commits().
			PressEscape()

		t.Views().Files().
			IsInvisible()

		t.Views().Commits().
			Lines(
				Contains(`none of the two`),
				Contains(`both files`).IsSelected(),
				Contains(`only otherFile`),
				Contains(`only filterFile`),
			)
	},
})
