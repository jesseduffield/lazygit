package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewWorktreePickerRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "From the worktrees panel, picking a remote branch creates a new local tracking branch and worktree; remote branches whose local branch already exists are filtered out",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranch("feature")
		shell.NewBranch("existing")
		shell.CloneIntoRemote("origin")
		shell.Checkout("master")
		// "feature" now exists only on the remote; "existing" stays as a local
		// branch (not checked out) that also has a remote counterpart
		shell.RunCommand([]string{"git", "branch", "-D", "feature"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Press(keys.Universal.New).
			Tap(func() {
				// master is checked out, so neither it nor origin/master is
				// offered; "existing" already has a local branch, so
				// origin/existing is left out too (you'd reach it via the local
				// entry); origin/feature has no local branch, so it's offered
				t.ExpectPopup().Prompt().
					Title(Equals("New worktree for branch")).
					SuggestionLines(
						Contains("existing"),
						Contains("origin/feature"),
					).
					Type("origin/feature").
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					Confirm()
			}).
			IsFocused().
			Lines(
				Contains("feature").IsSelected(),
				Contains("(main worktree)"),
			)

		// the new worktree is on a local branch that tracks the remote one (the
		// ✓ confirms tracking is set up)
		t.Views().Branches().
			Focus().
			ContainsLines(
				Contains("feature").Contains("✓").IsSelected(),
			)
	},
})
