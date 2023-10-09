package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterRemotes = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filtering remotes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit-one")
		shell.CloneIntoRemote("remote1")
		shell.CloneIntoRemote("remote2")
		shell.CloneIntoRemote("remote3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Remotes().
			Focus().
			Lines(
				Contains("remote1").IsSelected(),
				Contains("remote2"),
				Contains("remote3"),
			).
			FilterOrSearch("2").
			Lines(
				Contains("remote2").IsSelected(),
			).
			// cancel the filter
			PressEscape().
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Lines(
				Contains("remote1"),
				Contains("remote2").IsSelected(),
				Contains("remote3"),
			)
	},
})
