package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateFixupCommitInBranchStack = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a fixup commit in a stack of branches, verify that it is created at the end of the branch it belongs to",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch1")
		shell.EmptyCommit("branch1 commit 1")
		shell.EmptyCommit("branch1 commit 2")
		shell.EmptyCommit("branch1 commit 3")
		shell.NewBranch("branch2")
		shell.EmptyCommit("branch2 commit 1")
		shell.EmptyCommit("branch2 commit 2")
		shell.CreateFileAndAdd("fixup-file", "fixup content")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ branch2 commit 2"),
				Contains("CI ◯ branch2 commit 1"),
				Contains("CI ◯ * branch1 commit 3"),
				Contains("CI ◯ branch1 commit 2"),
				Contains("CI ◯ branch1 commit 1"),
			).
			NavigateToLine(Contains("branch1 commit 2")).
			Press(keys.Commits.CreateFixupCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create fixup commit")).
					Select(Contains("fixup! commit")).
					Confirm()
			}).
			Lines(
				Contains("CI ◯ branch2 commit 2"),
				Contains("CI ◯ branch2 commit 1"),
				Contains("CI ◯ * fixup! branch1 commit 2"),
				Contains("CI ◯ branch1 commit 3"),
				Contains("CI ◯ branch1 commit 2"),
				Contains("CI ◯ branch1 commit 1"),
			)
	},
})
