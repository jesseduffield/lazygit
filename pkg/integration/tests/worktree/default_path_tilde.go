package worktree

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DefaultPathTilde = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A leading ~ in the worktree.defaultPath config is expanded to the home directory",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Worktree.DefaultPath = "~/my-worktrees"
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
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

				t.ExpectPopup().Prompt().
					Title(Equals("New branch and worktree name")).
					Type("newbranch").
					Confirm()

				// The default path's "~" is expanded to an absolute home-directory
				// path; without expansion it would stay a literal "~" resolved
				// against the repo, so the candidate would still contain a "~".
				home, _ := os.UserHomeDir()
				t.ExpectPopup().Menu().
					Title(Equals("Worktree location")).
					ContainsLines(
						Contains(filepath.Join(home, "my-worktrees", "newbranch")).DoesNotContain("~"),
					).
					Cancel()
			})
	},
})
