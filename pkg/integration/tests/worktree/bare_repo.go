package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BareRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in the worktree of a bare repo and do a rebase/bisect",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - .bare
		//  - repo (a worktree)
		//  - worktree2 (another worktree)
		//
		// The first repo is called 'repo' because that's the
		// directory that all lazygit tests start in

		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("blah", "blah")
		shell.Commit("initial commit")
		shell.EmptyCommit("commit two")
		shell.EmptyCommit("commit three")

		shell.RunCommand([]string{"git", "clone", "--bare", ".", "../.bare"})

		shell.DeleteFile(".git")

		shell.Chdir("..")

		// This is the dir we were just in (and the dir that lazygit starts in when the test runs)
		// We're going to replace it with a worktree
		shell.DeleteFile("repo")

		shell.RunCommand([]string{"git", "--git-dir", ".bare", "worktree", "add", "-b", "repo", "repo", "mybranch"})
		shell.RunCommand([]string{"git", "--git-dir", ".bare", "worktree", "add", "-b", "worktree2", "worktree2", "mybranch"})

		shell.Chdir("repo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("repo"),
				Contains("mybranch"),
				Contains("worktree2 (worktree)"),
			)

		// test that a rebase works fine
		// (rebase uses the git dir of the worktree so we're confirming that it points
		// to the right git dir)
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit three").IsSelected(),
				Contains("commit two"),
				Contains("initial commit"),
			).
			Press(keys.Commits.MoveDownCommit).
			Lines(
				Contains("commit two"),
				Contains("commit three").IsSelected(),
				Contains("initial commit"),
			).
			// test that bisect works fine (same logic as above)
			NavigateToLine(Contains("commit two")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Bisect")).
					Select(MatchesRegexp(`Mark .* as bad`)).
					Confirm()

				t.Views().Information().Content(Contains("Bisecting"))
			}).
			NavigateToLine(Contains("initial commit")).
			Press(keys.Commits.ViewBisectOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Bisect")).
					Select(MatchesRegexp(`Mark .* as good`)).
					Confirm()

				t.Views().Information().Content(Contains("Bisecting"))
			})

		// switch to other worktree
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo").IsSelected(),
				Contains("worktree2"),
			).
			NavigateToLine(Contains("worktree2")).
			Press(keys.Universal.Select).
			Lines(
				Contains("worktree2").IsSelected(),
				Contains("repo"),
			)
	},
})
