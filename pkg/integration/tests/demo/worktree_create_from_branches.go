package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var WorktreeCreateFromBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a worktree from the branches view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(cfg *config.AppConfig) {
		// No idea why I had to use version 2: it should be using my own computer's
		// font and the one iterm uses is version 3.
		cfg.UserConfig.Gui.NerdFontsVersion = "2"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/user-authentication")
		shell.EmptyCommit("Add user authentication feature")
		shell.EmptyCommit("Fix local session storage")
		shell.CreateFile("src/authentication.go", "package main")
		shell.CreateFile("src/shims.go", "package main")
		shell.CreateFile("src/session.go", "package main")
		shell.EmptyCommit("Stop using shims")
		shell.UpdateFile("src/authentication.go", "package authentication")
		shell.UpdateFileAndAdd("src/shims.go", "// removing for now")
		shell.UpdateFile("src/session.go", "package session")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Create a worktree from a branch")
		t.Wait(1000)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("master")).
			Wait(500).
			Press(keys.Worktrees.ViewWorktreeOptions).
			Tap(func() {
				t.Wait(500)

				t.ExpectPopup().Menu().
					Title(Equals("Worktree")).
					Select(Contains("Create worktree from master").DoesNotContain("detached")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("New worktree path")).
					Type("../hotfix").
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("New branch name")).
					Type("hotfix/db-on-fire").
					Confirm()
			})
	},
})
