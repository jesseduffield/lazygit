package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartInGitDir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start lazygit in a repo's .git dir, and have it open the repo",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("blah", "original content\n")
		shell.Commit("initial commit")
		shell.UpdateFile("blah", "updated content\n")

		// this is where lazygit will start
		shell.Chdir(".git")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("initial commit"),
			)

		// we're in the work tree the .git belongs to, not in the .git itself
		t.Views().Files().
			IsFocused().
			Lines(
				Contains(" M blah"),
			)
	},
})
