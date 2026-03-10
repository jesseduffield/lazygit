package spice

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Note: This test requires git-spice (gs) to be installed.
var DeleteBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a branch from the spice stacks view via the menu",
	ExtraCmdArgs: []string{},
	Skip:         true, // Requires git-spice to be installed
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("main")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.RunCommand([]string{"gs", "repo", "init", "--trunk", "main"})
		shell.RunCommand([]string{"gs", "branch", "create", "feature-to-delete"})
		shell.CreateFileAndAdd("feature.txt", "feature content")
		shell.Commit("add feature")
		// Go back to main before deleting
		shell.Checkout("main")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().SpiceStacks().
			Focus().
			NavigateToLine(Contains("feature-to-delete")).
			Press("S").
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Stack Operations")).
					Select(Contains("Delete branch")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Delete branch")).
					Content(Contains("Are you sure you want to delete")).
					Confirm()
			}).
			Lines(
				Contains("main"),
			)
	},
})
