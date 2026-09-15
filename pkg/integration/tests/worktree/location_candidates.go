package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LocationCandidates = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The location menu offers the parents of existing worktrees and the configured default path, and sanitizes the typed name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Worktree.DefaultPath = "../config-worktrees"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		// a pre-existing linked worktree, so its parent directory is offered as
		// a candidate location alongside the configured default path
		shell.RunCommand([]string{"git", "worktree", "add", "-b", "existing", "../manual-worktrees/existing"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("mybranch")).
			Press(keys.Universal.NewWorktree).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("New worktree")).
					Select(Contains("New branch and worktree from 'mybranch'")).
					Confirm()

				// the space is sanitized to a dash so it's a valid branch (and
				// directory) name
				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					Type("new feature").
					Confirm()

				// the parent of the existing worktree comes first, then the
				// configured default path; both target the sanitized name
				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					ContainsLines(
						Contains("manual-worktrees").Contains("new-feature"),
						Contains("config-worktrees").Contains("new-feature"),
					).
					Cancel()
			})
	},
})
