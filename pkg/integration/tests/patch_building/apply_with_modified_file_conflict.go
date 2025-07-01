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

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Equals("error: file1: does not match index")).
			Confirm()
	},
})
