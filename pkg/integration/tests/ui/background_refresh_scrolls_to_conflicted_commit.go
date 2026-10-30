package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BackgroundRefreshScrollsToConflictedCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A background refresh that selects a newly conflicted commit scrolls it into view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file.txt", "master content")
		shell.Commit("master commit")
		shell.NewBranchFrom("branch", "master^")
		shell.CreateFileAndAdd("file.txt", "branch content")
		shell.Commit("this will conflict")
		shell.CreateNCommits(30)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("commit-30")).
			Tap(func() {
				t.Shell().RunCommandExpectError([]string{"git", "rebase", "master"})
				t.RefreshInBackground()
			}).
			// The conflicted commit, below the 30 pending todos and their header
			SelectedLineIdx(31).
			/* EXPECTED:
			SelectedLineIsVisible().
			SelectedLine(Contains("<-- CONFLICT --- this will conflict"))
			ACTUAL: */
			OriginY(0)
	},
})
