package spice

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Note: This test requires git-spice (gs) to be installed.
var CheckoutBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a branch from the spice stacks view",
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
		// Go back to main
		shell.Checkout("main")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().SpiceStacks().
			Focus().
			Lines(
				Contains("feature-1"),
				Contains("main").IsSelected(),
			).
			NavigateToLine(Contains("feature-1")).
			Press(keys.Universal.Select)

		t.Views().Branches().
			Lines(
				Contains("feature-1").IsSelected(),
				Contains("main"),
			)
	},
})
