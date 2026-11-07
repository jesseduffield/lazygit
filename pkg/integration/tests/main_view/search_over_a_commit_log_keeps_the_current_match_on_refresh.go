package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchOverACommitLogKeepsTheCurrentMatchOnRefresh = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Searching a branch's commit log in the main view, which shows no selection, keeps the current match when a refresh renders the log again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			FilterOrSearch("commit-0").
			Press(keys.Universal.NextMatch).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'commit-0' (2 of 3)"))

				t.GlobalPress(keys.Universal.Refresh)
				t.Views().Search().Content(
					/* EXPECTED:
					Contains("matches for 'commit-0' (2 of 3)"))
					ACTUAL: */
					Contains("matches for 'commit-0' (1 of 3)"))
			})
	},
})
