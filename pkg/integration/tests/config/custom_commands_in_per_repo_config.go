package config

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomCommandsInPerRepoConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Custom commands in per-repo config add to the global ones instead of replacing them",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{otherRepo}

		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "global",
				Command: "printf 'global X' > file.txt",
			},
			{
				Key:     "Y",
				Context: "global",
				Command: "printf 'global Y' > file.txt",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CloneNonBare("other")
		shell.CreateFile("../other/.git/lazygit.yml", `
customCommands:
  - key: Y
    context: global
    command: printf 'local Y' > file.txt
  - key: Z
    context: global
    command: printf 'local Z' > file.txt`)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other").IsSelected(),
				Contains("Cancel"),
			).Confirm()
		t.Views().Status().Content(Contains("other → master"))

		t.GlobalPress("X")
		t.FileSystem().FileContent("../other/file.txt", Equals("global X"))

		t.GlobalPress("Y")
		t.FileSystem().FileContent("../other/file.txt", Equals("local Y"))

		t.GlobalPress("Z")
		t.FileSystem().FileContent("../other/file.txt", Equals("local Z"))
	},
})
