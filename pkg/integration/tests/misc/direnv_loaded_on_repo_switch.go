package misc

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Verifies that when the user switches repos from inside lazygit, env vars
// that direnv would load for the target repo are applied to subprocesses
// (custom commands, git hooks, etc.). The test puts a fake `direnv` binary
// on PATH so it works regardless of whether the host has real direnv
// installed.
var DirenvLoadedOnRepoSwitch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching repos applies direnv-loaded env vars to subprocesses",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		// Prepend a dir under the test fixture to PATH so our fake direnv
		// wins lookup. The placeholder is resolved at run time.
		"PATH": "{{actualPath}}/bin:" + os.Getenv("PATH"),
	},
	SetupConfig: func(cfg *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{otherRepo}
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     config.Keybinding{"X"},
				Context: "files",
				Command: `echo "VAR=$LG_DIRENV_TEST" > output.txt`,
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.CloneNonBare("other")

		// Fake direnv: echoes a fixed JSON delta on stdout (set
		// LG_DIRENV_TEST) and a "loading" line on stderr, exactly as
		// real direnv would after authorizing an .envrc.
		shell.CreateFile("../bin/direnv", `#!/bin/sh
echo '{"LG_DIRENV_TEST":"from_direnv"}'
echo "direnv: loading .envrc" >&2
`)
		shell.MakeExecutable("../bin/direnv")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Switch to the "other" repo via the recent-repos menu.
		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other").IsSelected(),
				Contains("Cancel"),
			).
			Confirm()

		// Run the custom command; if direnv loading worked, $LG_DIRENV_TEST
		// reaches the subprocess and ends up in output.txt.
		t.Views().Files().
			Focus().
			Press(config.Keybinding{"X"}).
			Lines(
				Contains("output.txt").IsSelected(),
			)
		t.Views().Main().Content(Contains("VAR=from_direnv"))
	},
})
