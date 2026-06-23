package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectedCommitAfterExternalCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep the same commit selected after an external commit is created",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "first content")
		shell.Commit("first commit")
		shell.UpdateFile("file", "second content")
		shell.GitAddAll()
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit"),
				Contains("first commit"),
			).
			NavigateToLine(Contains("first commit"))

		t.Views().Main().Content(Contains("+first content"))

		t.GlobalPress(keys.Universal.ExecuteShellCommand)
		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			Type("git commit --allow-empty -m 'external commit'").
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("external commit"),
				Contains("second commit"),
				Contains("first commit").IsSelected(),
			)

		t.Views().Main().Content(Contains("+first content"))
	},
})
