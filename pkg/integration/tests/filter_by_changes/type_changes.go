package filter_by_changes

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var TypeChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by author using the typed in author",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(4)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter changes to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter changes:")).
			Type("file01").
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 01"),
			)

		t.Views().Information().Content(Contains("Filtering by 'file01'"))

		t.Views().Status().
			Focus().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter changes to filter by")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter changes:")).
			Type("file02").
			Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("commit 02"),
			)

		t.Views().Information().Content(Contains("Filtering by 'file02'"))
	},
})
