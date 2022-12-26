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

		input.SwitchToCommitsView()

		assert.CurrentView().Lines(
			Contains("first commit"),
		)

		input.Press(keys.Commits.RevertCommit)
		input.InConfirm().
			Title(Equals("Revert commit")).
			Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
			Confirm()

		assert.CurrentView().Name("commits").
			Lines(
				Contains("Revert \"first commit\"").IsSelected(),
				Contains("first commit"),
			)

		assert.MainView().Content(Contains("-myfile content"))
		assert.FileSystemPathNotPresent("myfile")
	},
})
