package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepClickedCommitSelectedAfterFocusIn = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep a clicked commit selected when focus-in immediately precedes the click",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit-02").IsSelected(),
				Contains("commit-01"),
			).
			FocusInAndClick(1, 1).
			SelectedLine(Contains("commit-01"))
	},
})
