package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func commonSetup(shell *Shell) {
	shell.CreateFileAndAdd("filterFile", "original filterFile content")
	shell.Commit("only filterFile")
	shell.CreateFileAndAdd("otherFile", "original otherFile content")
	shell.Commit("only otherFile")

	shell.UpdateFileAndAdd("otherFile", "new otherFile content")
	shell.UpdateFileAndAdd("filterFile", "new filterFile content")
	shell.Commit("both files")

	shell.EmptyCommit("none of the two")
}

func filterByFilterFile(t *TestDriver, keys config.KeybindingConfig) {
	t.Views().Commits().
		Focus().
		Lines(
			Contains(`none of the two`).IsSelected(),
			Contains(`both files`),
			Contains(`only otherFile`),
			Contains(`only filterFile`),
		).
		Press(keys.Universal.FilteringMenu)

	t.ExpectPopup().Menu().
		Title(Equals("Filtering")).
		Select(Contains("Enter path to filter by")).
		Confirm()
	t.ExpectPopup().Prompt().
		Title(Equals("Enter path:")).
		Type("filterF").
		SuggestionLines(Equals("filterFile")).
		ConfirmFirstSuggestion()
}

func postFilterTest(t *TestDriver) {
	t.Views().Information().Content(Contains("Filtering by 'filterFile'"))

	t.Views().Commits().
		IsFocused().
		Lines(
			Contains(`both files`).IsSelected(),
			Contains(`only filterFile`),
		)

	// we only show the filtered file's changes in the main view
	t.Views().Main().
		ContainsLines(
			Equals("    both files"),
			Equals("---"),
			Equals(" filterFile | 2 +-"),
			Equals(" 1 file changed, 1 insertion(+), 1 deletion(-)"),
		)

	t.Views().Commits().
		PressEnter()

	// when you click into the commit itself, you see all files from that commit
	t.Views().CommitFiles().
		IsFocused().
		Lines(
			Equals("▼ /"),
			Contains(`filterFile`),
			Contains(`otherFile`),
		)
}
