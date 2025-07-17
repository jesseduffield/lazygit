package filter_by_path

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowDiffsForRenamedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter commits by file path for a file that was renamed, and verify that it shows the diffs correctly",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("oldFile", "a\nb\nc\n")
		shell.Commit("add old file")
		shell.UpdateFileAndAdd("oldFile", "x\nb\nc\n")
		shell.Commit("update old file")
		shell.CreateFileAndAdd("unrelatedFile", "content of unrelated file\n")
		shell.Commit("add unrelated file")
		shell.RenameFileInGit("oldFile", "newFile")
		shell.Commit("rename file")
		shell.UpdateFileAndAdd("newFile", "y\nb\nc\n")
		shell.UpdateFileAndAdd("unrelatedFile", "updated content of unrelated file\n")
		shell.Commit("update both files")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("update both files").IsSelected(),
				Contains("rename file"),
				Contains("add unrelated file"),
				Contains("update old file"),
				Contains("add old file"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  M newFile"),
				Equals("  M unrelatedFile"),
			).
			SelectNextItem().
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().Title(Equals("Filtering")).
			Select(Contains("Filter by 'newFile'")).Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("update both files").IsSelected(),
				Contains("rename file"),
				Contains("update old file"),
				Contains("add old file"),
			)

		t.Views().Main().ContainsLines(
			Equals("    update both files"),
			Equals("---"),
			Equals(" newFile | 2 +-"),
			Equals(" 1 file changed, 1 insertion(+), 1 deletion(-)"),
			Equals(""),
			Equals("diff --git a/newFile b/newFile"),
			Contains("index"),
			Equals("--- a/newFile"),
			Equals("+++ b/newFile"),
			Equals("@@ -1,3 +1,3 @@"),
			Equals("-x"),
			Equals("+y"),
			Equals(" b"),
			Equals(" c"),
		)

		t.Views().Commits().SelectNextItem()

		t.Views().Main().ContainsLines(
			Equals("    rename file"),
			Equals("---"),
			Equals(" oldFile => newFile | 0"),
			Equals(" 1 file changed, 0 insertions(+), 0 deletions(-)"),
			Equals(""),
			Equals("diff --git a/oldFile b/newFile"),
			Equals("similarity index 100%"),
			Equals("rename from oldFile"),
			Equals("rename to newFile"),
		)

		t.Views().Commits().SelectNextItem()

		t.Views().Main().ContainsLines(
			Equals("    update old file"),
			Equals("---"),
			Equals(" oldFile | 2 +-"),
			Equals(" 1 file changed, 1 insertion(+), 1 deletion(-)"),
			Equals(""),
			Equals("diff --git a/oldFile b/oldFile"),
			Contains("index"),
			Equals("--- a/oldFile"),
			Equals("+++ b/oldFile"),
			Equals("@@ -1,3 +1,3 @@"),
			Equals("-a"),
			Equals("+x"),
			Equals(" b"),
			Equals(" c"),
		)

		t.Views().Commits().
			Press(keys.Universal.RangeSelectUp).
			Press(keys.Universal.RangeSelectUp)

		t.Views().Main().ContainsLines(
			Contains("Showing diff for range"),
			Equals(""),
			Equals(" oldFile => newFile | 2 +-"),
			Equals(" 1 file changed, 1 insertion(+), 1 deletion(-)"),
			Equals(""),
			Equals("diff --git a/oldFile b/newFile"),
			Equals("similarity index 66%"),
			Equals("rename from oldFile"),
			Equals("rename to newFile"),
			Contains("index"),
			Equals("--- a/oldFile"),
			Equals("+++ b/newFile"),
			Equals("@@ -1,3 +1,3 @@"),
			Equals("-a"),
			Equals("+y"),
			Equals(" b"),
			Equals(" c"),
		)
	},
})
