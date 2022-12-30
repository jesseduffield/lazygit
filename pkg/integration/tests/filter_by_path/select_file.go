package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path, by finding file in UI and filtering on it",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		commonSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains(`only filterFile`).IsSelected(),
				Contains(`only otherFile`),
				Contains(`both files`),
			).
			PressEnter()

		// when you click into the commit itself, you see all files from that commit
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains(`filterFile`).IsSelected(),
			).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().Title(Equals("Filtering")).Select(Contains("filter by 'filterFile'")).Confirm()

		postFilterTest(t)
	},
})

func commonSetup(shell *Shell) {
	shell.CreateFileAndAdd("filterFile", "original filterFile content")
	shell.CreateFileAndAdd("otherFile", "original otherFile content")
	shell.Commit("both files")

	shell.UpdateFileAndAdd("otherFile", "new otherFile content")
	shell.Commit("only otherFile")

	shell.UpdateFileAndAdd("filterFile", "new filterFile content")
	shell.Commit("only filterFile")
}

func postFilterTest(t *TestDriver) {
	t.Views().Information().Content(Contains("filtering by 'filterFile'"))

	t.Views().Commits().
		IsFocused().
		Lines(
			Contains(`only filterFile`).IsSelected(),
			Contains(`both files`),
		).
		SelectNextItem().
		PressEnter()

	// we only show the filtered file's changes in the main view
	t.Views().Main().
		Content(Contains("filterFile").DoesNotContain("otherFile"))

	// when you click into the commit itself, you see all files from that commit
	t.Views().CommitFiles().
		IsFocused().
		Lines(
			Contains(`filterFile`),
			Contains(`otherFile`),
		)
}
