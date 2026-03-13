package misc

import (
	"os"
	"runtime"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitWarning = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify lazygit handles git warnings without crashing",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"PATH": "{{actualRepoPath}}/bin:" + os.Getenv("PATH"),
	},
	Skip:        runtime.GOOS == "windows", // Shell wrapper won't work on Windows
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Create bin directory for our git wrapper
		shell.CreateDir("bin")

		// Create a git wrapper that outputs the warning then calls real git
		shell.CreateFile("bin/git", `#!/bin/sh
echo "warning: unhandled Platform key FamilyDisplayName" >&2
exec /usr/bin/git "$@"
`)
		shell.MakeExecutable("bin/git")

		// Create a simple repo with a commit
		shell.CreateFileAndAdd("file.txt", "content").
			Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Simply verify lazygit opens and shows the files view without crashing
		// The git wrapper outputs the warning on every git command
		t.Views().Files().IsFocused()

		// Navigate to branches to trigger more git commands
		t.Views().Branches().Focus()
		t.Views().Branches().Lines(
			Contains("master"),
		)

		// Navigate to commits to trigger git log
		t.Views().Commits().Focus()
		t.Views().Commits().Lines(
			Contains("initial commit"),
		)
	},
})
