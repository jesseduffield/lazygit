package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitRewordMultiline = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reword commit with a multi-line commit message",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.GitAddAll()
		shell.Commit("first line")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first line"),
			).
			Press(keys.Commits.RenameCommit)

		t.ExpectPopup().RewordCommitPanel().AddNewline().Type("second line").Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(MatchesRegexp("first line\n\\s*second line"))
	},
})
