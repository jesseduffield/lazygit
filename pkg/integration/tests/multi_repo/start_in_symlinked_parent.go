package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartInSymlinkedParent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start in a symlink to a directory of repos and see the current repo marked",
	ExtraEnvVars: ceilingEnv,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setupWorkspace(shell)
		shell.Chdir("..")
		shell.RunCommand([]string{"ln", "-s", "workspace", "link"})
		shell.Chdir("link")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Repos().
			IsFocused().
			Lines(
				Contains("  * alpha").Contains("feature").IsSelected(),
				DoesNotContain("  *").Contains("beta"),
				DoesNotContain("  *").Contains("org/gamma"),
			)
	},
})
