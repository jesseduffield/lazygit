package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyWithModifiedFileNoConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a custom patch, with a modified file in the working tree that does not conflict with the patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "1\n2\n3\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "1\n2\n3\n4\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
		shell.UpdateFile("file1", "11\n2\n3\n")
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

		t.Views().Secondary().Content(Contains("3\n+4"))

		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		t.ExpectPopup().Confirmation().Title(Equals("Must stage files")).
			Content(Contains("Applying a patch to the index requires staging the unstaged files that are affected by the patch.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Equals("M  file1").IsSelected(),
			)

		t.Views().Main().
			Content(Contains("-1\n+11\n 2\n 3\n+4"))
	},
})
