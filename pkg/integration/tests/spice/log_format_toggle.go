package spice

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Note: This test requires git-spice (gs) to be installed.
var LogFormatToggle = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle log format between short and long in spice stacks view",
	ExtraCmdArgs: []string{},
	Skip:         true, // Requires git-spice to be installed
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("main")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.RunCommand([]string{"gs", "repo", "init", "--trunk", "main"})
		shell.RunCommand([]string{"gs", "branch", "create", "feature-1"})
		shell.CreateFileAndAdd("feature.txt", "feature content")
		shell.Commit("add feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().SpiceStacks().
			Focus().
			Press("V").
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Log Format")).
					Lines(
						Contains("Short"),
						Contains("Long"),
						Contains("Default"),
					).
					Select(Contains("Long")).
					Confirm()
			})
		// View should still be focused after changing format
		t.Views().SpiceStacks().IsFocused()
	},
})
