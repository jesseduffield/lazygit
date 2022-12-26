package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Revert = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a commit",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.GitAddAll()
		shell.Commit("first commit")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(1)

		input.SwitchToCommitsWindow()

		assert.CurrentView().Name("commits").Lines(
			Contains("first commit"),
		)

		input.Press(keys.Commits.RevertCommit)
		input.AcceptConfirmation(Equals("Revert commit"), MatchesRegexp(`Are you sure you want to revert \w+?`))

		assert.CurrentView().Name("commits").
			Lines(
				Contains("Revert \"first commit\""),
				Contains("first commit"),
			).
			SelectedLineIdx(0)

		assert.MainView().Content(Contains("-myfile content"))
		assert.FileSystemPathNotPresent("myfile")
	},
})
