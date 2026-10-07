package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SkipBrokenRepos = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start in a directory whose first repos are broken and enter the first usable one",
	ExtraEnvVars: ceilingEnv,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setupWorkspace(shell)
		shell.CreateDir("aaa-broken/.git")
		shell.CreateFile("aab-stale/.git", "gitdir: ../nowhere\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Repos().
			IsFocused().
			Lines(
				DoesNotContain("  *").Contains("aaa-broken"),
				DoesNotContain("  *").Contains("aab-stale"),
				Contains("  * alpha").Contains("feature"),
				DoesNotContain("  *").Contains("beta"),
				DoesNotContain("  *").Contains("org/gamma"),
			)
	},
})
