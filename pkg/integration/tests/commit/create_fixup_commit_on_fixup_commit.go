package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateFixupCommitOnFixupCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a fixup commit on an existing fixup commit, verify that it prompts you to create it on the base commit",
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
		shell.EmptyCommit("fixup! branch2 commit 1")
		shell.CreateFileAndAdd("fixup-file", "fixup content")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ fixup! branch2 commit 1"),
				Contains("CI ◯ branch2 commit 1"),
				Contains("CI ◯ * branch1 commit 3"),
				Contains("CI ◯ branch1 commit 2"),
				Contains("CI ◯ branch1 commit 1"),
			).
			NavigateToLine(Contains("fixup! branch2 commit 1")).
			Press(keys.Commits.CreateFixupCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create fixup commit")).
					Select(Contains("fixup! commit")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Select base commit")).
					Select(Equals("b base commit")).
					Confirm()
			}).
			Lines(
				Contains("CI ◯ fixup! branch2 commit 1"),
				Contains("CI ◯ fixup! branch2 commit 1"),
				Contains("CI ◯ branch2 commit 1"),
				Contains("CI ◯ * branch1 commit 3"),
				Contains("CI ◯ branch1 commit 2"),
				Contains("CI ◯ branch1 commit 1"),
			)
	},
})
