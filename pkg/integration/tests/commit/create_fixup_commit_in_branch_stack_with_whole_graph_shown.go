package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateFixupCommitInBranchStackWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a fixup commit in a stack of branches while the whole graph is shown and another branch's head sits between the fixed-up commit and the head of its branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.Order = "date-order"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithDate("base", "2026-01-01T00:00:00Z").
			NewBranch("branch1").
			EmptyCommitWithDate("branch1 commit 1", "2026-01-01T00:01:00Z").
			NewBranchFrom("other", "master").
			EmptyCommitWithDate("other commit", "2026-01-01T00:02:00Z").
			Checkout("branch1").
			EmptyCommitWithDate("branch1 commit 2", "2026-01-01T00:03:00Z").
			NewBranch("branch2").
			EmptyCommitWithDate("branch2 commit 1", "2026-01-01T00:04:00Z").
			CreateFileAndAdd("fixup-file", "fixup content").
			SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("branch2 commit 1"),
				Contains("branch1 commit 2"),
				Contains("other commit"),
				Contains("branch1 commit 1"),
				Contains("base"),
			).
			NavigateToLine(Contains("branch1 commit 1")).
			Press(keys.Commits.CreateFixupCommit).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Create fixup commit")).
					Select(Contains("fixup! commit")).
					Confirm()
			})

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch1")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("fixup! branch1 commit 1"),
				Contains("branch1 commit 2"),
				Contains("branch1 commit 1"),
				Contains("base"),
			)
	},
})
