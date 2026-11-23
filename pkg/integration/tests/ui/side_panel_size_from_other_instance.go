package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SidePanelSizeFromOtherInstance = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A side panel size that another lazygit instance saved takes effect when the window regains focus",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			HasWidth(40)

		// The other instance ran in a narrower terminal; the ratio applies here
		// too, because it's the only size saved for the layout
		t.UpdateAppStateFile(func(state *config.AppState) {
			state.SidePanelSizes = map[config.SidePanelLayout][]config.SidePanelSize{
				config.SidePanelLayoutNormal: {{ScreenSize: 80, Ratio: 0.25}},
			}
		})

		t.Views().Files().
			HasWidth(40)

		t.FocusIn()

		t.Views().Files().
			HasWidth(30)
	},
})
