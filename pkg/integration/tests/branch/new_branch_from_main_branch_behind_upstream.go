package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchFromMainBranchBehindUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch from a main branch that is behind its upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
		shell.CloneIntoRemote("origin")
		shell.PushBranchAndSetUpstream("origin", "master")
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master ↓1").DoesNotContain("↑").IsSelected(),
			).
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Contains("New branch name (branch is off of 'origin/master')")).
			Type("new-branch").
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("new-branch").IsSelected(),
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			)
	},
})
