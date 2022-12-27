package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the diff between two commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("second commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("third commit")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsView()

		assert.Views().Current().Lines(
			Contains("third commit"),
			Contains("second commit"),
			Contains("first commit"),
		)

		input.Press(keys.Universal.DiffingMenu)
		input.Menu().Title(Equals("Diffing")).Select(MatchesRegexp(`diff \w+`)).Confirm()

		assert.NotInPopup()

		assert.Views().ByName("information").Content(Contains("showing output for: git diff"))

		input.NextItem()
		input.NextItem()
		assert.Views().Current().SelectedLine(Contains("first commit"))

		assert.Views().Main().Content(Contains("-second line\n-third line"))

		input.Press(keys.Universal.DiffingMenu)
		input.Menu().Title(Equals("Diffing")).Select(Contains("reverse diff direction")).Confirm()
		assert.NotInPopup()

		assert.Views().Main().Content(Contains("+second line\n+third line"))

		input.Enter()

		assert.Views().Current().Name("commitFiles").SelectedLine(Contains("file1"))
		assert.Views().Main().Content(Contains("+second line\n+third line"))
	},
})
