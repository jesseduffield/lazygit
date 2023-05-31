package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenPullRequestNoUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open up a pull request with a missing upstream branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().
			Branches().
			Focus().
			Press(keys.Branches.CreatePullRequest)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Cannot open a pull request for a branch with no upstream")).
			Confirm()
	},
})
