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
// previous repo's env leaks into the new one. This test exercises the
// "skip approval" branch: the approval popup appears, the user cancels,
// and the previous repo's env is still gone.
var DirenvUnloadsOnBlockedEnvrc = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Blocked .envrc unloads the previous repo's env even if the user skips approval",
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

		shell.CreateFile("../other/.envrc", "export LG_DIRENV_TEST=from_envrc\n")

		shell.CreateFile("../bin/direnv", `#!/bin/sh
case "$1 $2" in
"export json")
    echo '{"LG_DIRENV_TEST":null}'
    echo "direnv: error $PWD/.envrc is blocked" >&2
    exit 1
    ;;
"status --json")
    printf '{"state":{"foundRC":{"allowed":1,"path":"%s/.envrc"}}}\n' "$PWD"
    ;;
esac
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

		t.ExpectPopup().Confirmation().
			Title(Equals("Approve .envrc?")).
			Content(Contains("export LG_DIRENV_TEST=from_envrc")).
			Cancel()

		t.Views().Files().
			Focus().
			Press(config.Keybinding{"X"}).
			NavigateToLine(Contains("output.txt"))
		t.Views().Main().Content(Contains("VAR=[]"))
	},
})
