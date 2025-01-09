package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResolveWithoutTrailingLf = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Regression test for resolving a merge conflict when the file doesn't have a trailing newline",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateFileAndAdd("file", "a\n\nno eol").
			Commit("initial commit").
			UpdateFileAndAdd("file", "a1\n\nno eol").
			Commit("commit on branch1").
			NewBranchFrom("branch2", "HEAD^").
			UpdateFileAndAdd("file", "a2\n\nno eol").
			Commit("commit on branch2").
			Checkout("branch1").
			RunCommandExpectError([]string{"git", "merge", "--no-edit", "branch2"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			SelectedLines(
				Contains("<<<<<<< HEAD"),
				Contains("a1"),
				Contains("======="),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue?")).
			Cancel()

		t.Views().Files().
			Focus().
			Lines(
				Contains("M  file").IsSelected(),
			)

		t.Views().Main().Content(Contains("-a1\n+a2\n").DoesNotContain("-no eol"))
	},
})
