package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToNewCommitInLastCommitOfStackedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to a new commit, in the last commit of a branch in the middle of a stack",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("commit 01").
			NewBranch("branch1").
			EmptyCommit("commit 02").
			CreateFileAndAdd("file1", "file1 content").
			CreateFileAndAdd("file2", "file2 content").
			Commit("commit 03").
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 05").IsSelected(),
				Contains("CI commit 04"),
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 03")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  A file1"),
				Equals("  A file2"),
			).
			SelectNextItem().
			PressPrimaryAction().
			PressEscape()

		t.Views().Information().Content(Contains("Building patch"))

		t.Common().SelectPatchOption(Contains("Move patch into new commit after the original commit"))

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("")).
			Type("new commit").Confirm()

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit 05"),
				Contains("CI commit 04"),
				Contains("CI * new commit").IsSelected(),
				Contains("CI commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			)
	},
})
