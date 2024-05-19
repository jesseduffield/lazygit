package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FindBaseCommitForFixupOnlyAddedLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Finds the base commit to create a fixup for, when all staged hunks have only added lines",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch").
			EmptyCommit("1st commit").
			CreateFileAndAdd("file1", "line A\nline B\nline C\n").
			Commit("2nd commit").
			UpdateFileAndAdd("file1", "line A\nline B changed\nline C\n").
			Commit("3rd commit").
			CreateFileAndAdd("file2", "line X\nline Y\nline Z\n").
			Commit("4th commit").
			UpdateFile("file1", "line A\nline B changed\nline B'\nline C\n").
			UpdateFile("file2", "line W\nline X\nline Y\nline Z\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("4th commit"),
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
					Contains("3rd commit").
					Contains("4th commit"),
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
				Contains("4th commit"),
				Contains("3rd commit").IsSelected(),
				Contains("2nd commit"),
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
				Contains("4th commit").IsSelected(),
				Contains("3rd commit"),
				Contains("2nd commit"),
				Contains("1st commit"),
			)
	},
})
