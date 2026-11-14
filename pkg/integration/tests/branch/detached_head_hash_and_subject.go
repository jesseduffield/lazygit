package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DetachedHeadHashAndSubject = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show the commit hash and subject of a detached HEAD in the branches panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			Checkout("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextScreenMode).
			Lines(
				MatchesRegexp(`^  \* [0-9a-f]{8} \(HEAD detached at [0-9a-f]+\) {3}commit-02$`).IsSelected(),
				MatchesRegexp(`[0-9a-f]{8} master +commit-03$`),
			)
	},
})
