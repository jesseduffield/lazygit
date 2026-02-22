package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the {{ .SelectedCommits }} template variable for non-contiguous selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(5)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "commits",
				Command: `echo "{{ range .SelectedCommits }}{{ .Name }} {{ end }}" > file.txt`,
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().
			Lines(
				Contains("commit 05").IsSelected(),
				Contains("commit 04"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			// Mark commit 05 (currently selected)
			Press("z").
			// Move down to commit 03 and mark it
			NavigateToLine(Contains("commit 03")).
			Press("z").
			// Move down to commit 01 and mark it
			NavigateToLine(Contains("commit 01")).
			Press("z")

		// Run the custom command which should output all marked commits
		t.GlobalPress("X")
		// The commits should be in index order (05, 03, 01)
		t.FileSystem().FileContent("file.txt", Equals("commit 05 commit 03 commit 01 \n"))
	},
})
