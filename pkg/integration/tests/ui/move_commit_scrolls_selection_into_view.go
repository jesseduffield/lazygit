package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveCommitScrollsSelectionIntoView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a commit down scrolls it into view if it isn't visible",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(40)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("commit-40")).
			// Scroll the selected commit out of view with the mouse wheel
			ScrollWheelDown().
			ScrollWheelDown().
			OriginY(4).
			Press(keys.Commits.MoveDownCommit).
			SelectedLine(Contains("commit-40")).
			SelectedLineIdx(1).
			SelectedLineIsVisible()
	},
})
