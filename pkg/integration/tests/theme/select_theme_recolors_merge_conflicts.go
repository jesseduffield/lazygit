package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var SelectThemeRecolorsMergeConflicts = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme that sets the default text color recolors a merge conflict, both where the files panel shows it and in the merge conflicts view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    defaultFgColor:\n      - '#ff00ff'\n")
		writeThemeFile(cfg, "blue", "gui:\n  theme:\n    defaultFgColor:\n      - '#0000ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			)

		t.Views().MergeConflicts().
			Content(Contains("<<<<<<< HEAD")).
			DoesNotContainColoredText("#ff00ff", "This")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().MergeConflicts().
			ContainsColoredText("#ff00ff", "This")

		t.Views().Files().
			IsFocused().
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("First Change"),
				Contains("======="),
			).
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("blue")).
			Confirm()

		t.ExpectToast(Equals("Theme: blue"))

		t.Views().MergeConflicts().
			IsFocused().
			ContainsColoredText("#0000ff", "This")
	},
})
