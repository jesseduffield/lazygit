package misc

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// When the new repo's .envrc is blocked, lazygit offers the user a popup to
// approve it without leaving the app. Confirming runs `direnv allow` and
// re-runs the load so the env reaches subprocesses immediately.
var DirenvApprovesEnvrc = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Approving a blocked .envrc from the in-app popup loads its env",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
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

		shell.CreateFile("../other/.envrc", "export LG_DIRENV_TEST=approved_value\n")

		// Fake direnv that flips behavior once `direnv allow` runs.
		// Before allow: export errors with the "blocked" signal,
		//               status reports allowed=1 (NotAllowed).
		// On allow:     create a sentinel and exit 0.
		// After allow:  export emits the loaded delta normally.
		shell.CreateFile("../bin/direnv", `#!/bin/sh
SENTINEL="$(dirname "$0")/.approved"
case "$1 $2" in
"allow "*)
    touch "$SENTINEL"
    exit 0
    ;;
"export json")
    if [ -f "$SENTINEL" ]; then
        echo '{"LG_DIRENV_TEST":"approved_value"}'
        echo "direnv: loading $PWD/.envrc" >&2
    else
        echo '{"LG_DIRENV_TEST":null}'
        echo "direnv: error $PWD/.envrc is blocked" >&2
        exit 1
    fi
    ;;
"status --json")
    if [ -f "$SENTINEL" ]; then
        printf '{"state":{"foundRC":{"allowed":0,"path":"%s/.envrc"}}}\n' "$PWD"
    else
        printf '{"state":{"foundRC":{"allowed":1,"path":"%s/.envrc"}}}\n' "$PWD"
    fi
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
			Content(Contains("export LG_DIRENV_TEST=approved_value")).
			Confirm()

		t.Views().Files().
			Focus().
			Press(config.Keybinding{"X"}).
			NavigateToLine(Contains("output.txt"))
		t.Views().Main().Content(Contains("VAR=approved_value"))
	},
})
