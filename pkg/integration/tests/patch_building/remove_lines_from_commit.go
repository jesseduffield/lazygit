package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveLinesFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove specific lines from a commit using the 'd' shortcut in the patch building view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")

		shell.CreateFileAndAdd("file1", "1st line\n2nd line\n3rd line\n")
		shell.Commit("commit to remove from")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit to remove from").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file1").IsSelected(),
			).
			PressEnter()

		// Select the second line (+2nd line) and press 'd' to remove it
		t.Views().PatchBuilding().
			IsFocused().
			SelectNextItem().
			SelectedLines(
				Contains("+2nd line"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Remove lines from commit")).
			Content(Equals("Are you sure you want to remove the selected lines from this commit?")).
			Confirm()

		// After the rebase, we should be back at the commit files view
		// and the commit should now only contain the 1st and 3rd lines
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file1").IsSelected(),
			).
			PressEscape()

		t.Views().Main().ContainsLines(
			Equals("+1st line"),
			Equals("+3rd line"),
		)
	},
})
