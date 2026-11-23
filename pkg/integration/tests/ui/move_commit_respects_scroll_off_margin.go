package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveCommitRespectsScrollOffMargin = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a commit down scrolls when the selection enters the bottom scroll-off margin",
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
			SelectedLineIdx(0).
			OriginY(0).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			SelectedLineIdx(2).
			OriginY(0).
			Press(keys.Commits.MoveDownCommit).
			SelectedLineIdx(3).
			OriginY(1)
	},
})
