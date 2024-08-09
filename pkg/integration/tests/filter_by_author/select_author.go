package filter_by_author

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits using the currently highlighted commit's author when the commit view is active",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetAppState().GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLineIdx(0).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Filter by 'Paul Oberstein <paul.oberstein@email.com>'")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 7"),
				Contains("commit 6"),
				Contains("commit 5"),
				Contains("commit 4"),
				Contains("commit 3"),
				Contains("commit 2"),
				Contains("commit 1"),
				Contains("commit 0"),
			)

		t.Views().Information().Content(Contains("Filtering by 'Paul Oberstein <paul.oberstein@email.com>'"))

		t.Views().Commits().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Stop filtering")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("SK commit 0")).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Filter by 'Siegfried Kircheis <siegfried.kircheis@email.com>'")).
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 0"),
			)

		t.Views().Information().Content(Contains("Filtering by 'Siegfried Kircheis <siegfried.kircheis@email.com>'"))
	},
})
