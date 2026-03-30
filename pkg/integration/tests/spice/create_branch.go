package spice

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Note: This test requires git-spice (gs) to be installed.
var CreateBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a new branch from the spice stacks view using n key",
	ExtraCmdArgs: []string{},
	Skip:         true, // Requires git-spice to be installed
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("main")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.RunCommand([]string{"gs", "repo", "init", "--trunk", "main"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().SpiceStacks().
			Focus().
			Lines(
				Contains("main"),
			).
			Press("n").
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Branch name:")).
					Type("new-feature").
					Confirm()
			}).
			Lines(
				Contains("new-feature").IsSelected(),
				Contains("main"),
			)
	},
})
