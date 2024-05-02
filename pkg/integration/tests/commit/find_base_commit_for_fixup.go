package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch").
			EmptyCommit("1st commit").
			CreateFileAndAdd("file1", "file1 content\n").
			Commit("2nd commit").
			CreateFileAndAdd("file2", "file2 content\n").
			Commit("3rd commit").
			UpdateFile("file1", "file1 changed content").
			UpdateFile("file2", "file2 changed content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("3rd commit"),
				Contains("2nd commit"),
				Contains("1st commit"),
			)

		// Two changes from different commits: this fails
		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(
				Contains("Multiple base commits found").
					Contains("2nd commit").
					Contains("3rd commit"),
			).
			Confirm()

		// Stage only one of the files: this succeeds
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("file1")).
			PressPrimaryAction().
			Press(keys.Files.FindBaseCommitForFixup)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("3rd commit"),
				Contains("2nd commit").IsSelected(),
				Contains("1st commit"),
			).
			Press(keys.Commits.AmendToCommit)

		t.ExpectPopup().Confirmation().
			Title(Equals("Amend commit")).
			Content(Contains("Are you sure you want to amend this commit with your staged files?")).
			Confirm()

		// Now only the other file is modified (and unstaged); this works now
		t.Views().Files().
			Focus().
			Press(keys.Files.FindBaseCommitForFixup)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("3rd commit").IsSelected(),
				Contains("2nd commit"),
				Contains("1st commit"),
			)
	},
})
