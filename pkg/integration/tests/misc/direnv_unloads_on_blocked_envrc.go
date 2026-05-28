package misc

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Real direnv exits non-zero when the destination .envrc isn't authorized,
// but it still emits a valid JSON delta on stdout that unloads vars from
// the previously-active .envrc. We have to apply that delta anyway, or the
// previous repo's env leaks into the new one. The fake direnv here mimics
// that behavior; the test also asserts that the user gets an error popup
// (the command log alone is easy to miss).
var DirenvUnloadsOnBlockedEnvrc = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Blocked .envrc unloads the previous repo's env and shows an error popup",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"PATH": "{{actualPath}}/bin:" + os.Getenv("PATH"),
		// Simulates a var that the previous repo's .envrc would have set.
		"LG_DIRENV_TEST": "from_previous_repo",
	},
	SetupConfig: func(cfg *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{otherRepo}
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     config.Keybinding{"X"},
				Context: "files",
				Command: `echo "VAR=[$LG_DIRENV_TEST]" > output.txt`,
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.CloneNonBare("other")

		shell.CreateFile("../bin/direnv", `#!/bin/sh
echo '{"LG_DIRENV_TEST":null}'
echo "direnv: error /repo/.envrc is blocked. Run 'direnv allow' to approve its content" >&2
exit 1
`)
		shell.MakeExecutable("../bin/direnv")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other").IsSelected(),
				Contains("Cancel"),
			).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("is blocked")).
			Confirm()

		// If unload worked, $LG_DIRENV_TEST is empty in the custom command.
		t.Views().Files().
			Focus().
			Press(config.Keybinding{"X"}).
			Lines(
				Contains("output.txt").IsSelected(),
			)
		t.Views().Main().Content(Contains("VAR=[]"))
	},
})
