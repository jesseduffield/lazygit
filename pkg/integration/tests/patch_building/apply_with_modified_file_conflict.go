package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyWithModifiedFileConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a custom patch, with a modified file in the working tree that conflicts with the patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "1\n2\n3\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "11\n2\n3\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
		shell.UpdateFile("file1", "111\n2\n3\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.NextItem).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file1").IsSelected(),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Secondary().Content(Contains("-1\n+11\n"))

		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		t.ExpectPopup().Confirmation().Title(Equals("Must stage files")).
			Content(Contains("Applying a patch to the index requires staging the unstaged files that are affected by the patch.")).
			Confirm()

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("Applied patch to 'file1' with conflicts.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Equals("UU file1").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			Lines(
				Equals("<<<<<<< ours"),
				Equals("111"),
				Equals("======="),
				Equals("11"),
				Equals(">>>>>>> theirs"),
				Equals("2"),
				Equals("3"),
			)
	},
})
