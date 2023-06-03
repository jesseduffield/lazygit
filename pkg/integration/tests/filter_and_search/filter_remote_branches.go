package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterRemoteBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering remote branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-apple")
		shell.EmptyCommit("commit-one")
		shell.NewBranch("branch-grape")
		shell.NewBranch("branch-orange")

		shell.CloneIntoRemote("origin")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains(`origin`).IsSelected(),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			Lines(
				Contains(`branch-apple`).IsSelected(),
				Contains(`branch-grape`),
				Contains(`branch-orange`),
			).
			FilterOrSearch("grape").
			Lines(
				Contains(`branch-grape`).IsSelected(),
			).
			// cancel the filter
			PressEscape().
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Lines(
				Contains(`branch-apple`),
				Contains(`branch-grape`).IsSelected(),
				Contains(`branch-orange`),
			).
			// return to remotes view
			PressEscape()

		t.Views().Remotes().
			IsFocused().
			Lines(
				Contains(`origin`).IsSelected(),
			)
	},
})
