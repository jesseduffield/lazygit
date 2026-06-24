package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ReloadSidePanels = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Editing the side panel config and refocusing the window re-applies the layout live, keeping the focused panel focused",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
		// Start with worktrees promoted to its own panel.
		shell.CreateFile(".git/lazygit.yml", `
gui:
  sidePanels:
    - [status]
    - [files, submodules]
    - [worktrees]
    - [branches, remotes, tags]
    - [commits, reflog]
    - [stash]`)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Worktrees is its own panel in the third position.
		t.Views().Files().IsFocused().
			Press(keys.Universal.JumpToBlock[2])
		t.Views().Worktrees().IsFocused()

		// Demote worktrees back into the files panel, then refocus the window to
		// trigger a live reload of the changed config.
		t.Shell().UpdateFile(".git/lazygit.yml", `
gui:
  sidePanels:
    - [status]
    - [files, worktrees, submodules]
    - [branches, remotes, tags]
    - [commits, reflog]
    - [stash]`)
		t.FocusIn()

		// Worktrees is now a tab of the files panel. It stays focused, and is shown
		// in front rather than being hidden behind the files tab (which would leave
		// the panel looking unfocused).
		t.Views().Worktrees().IsActiveTab().IsFocused().
			Press(keys.Universal.PrevTab)
		t.Views().Files().IsActiveTab().IsFocused()
	},
})
