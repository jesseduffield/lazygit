package spice

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Note: This test requires git-spice (gs) to be installed.
// See: https://github.com/abhinav/git-spice
var BasicDisplay = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify the spice stacks view can be focused and displays content",
	ExtraCmdArgs: []string{},
	Skip:         true, // Requires git-spice to be installed
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("main")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		// Initialize git-spice
		shell.RunCommand([]string{"gs", "repo", "init", "--trunk", "main"})
		// Create a feature branch
		shell.RunCommand([]string{"gs", "branch", "create", "feature-1"})
		shell.CreateFileAndAdd("feature.txt", "feature content")
		shell.Commit("add feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().SpiceStacks().
			Focus().
			Lines(
				Contains("feature-1"),
				Contains("main"),
			)
	},
})
