package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Staged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFile("myfile", "myfile content\nwith a second line").
			CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(0)

		assert.CurrentView().Name("files")
		assert.CurrentView().SelectedLine(Contains("myfile"))
		// stage the file
		input.PrimaryAction()
		input.Enter()
		assert.CurrentView().Name("stagingSecondary")
		// we start with both lines having been staged
		assert.View("stagingSecondary").Content(Contains("+myfile content"))
		assert.View("stagingSecondary").Content(Contains("+with a second line"))
		assert.View("staging").Content(DoesNotContain("+myfile content"))
		assert.View("staging").Content(DoesNotContain("+with a second line"))

		// unstage the selected line
		input.PrimaryAction()

		// the line should have been moved to the main view
		assert.View("stagingSecondary").Content(DoesNotContain("+myfile content"))
		assert.View("stagingSecondary").Content(Contains("+with a second line"))
		assert.View("staging").Content(Contains("+myfile content"))
		assert.View("staging").Content(DoesNotContain("+with a second line"))

		input.Press(keys.Files.CommitChanges)
		commitMessage := "my commit message"
		input.Type(commitMessage)
		input.Confirm()

		assert.CommitCount(1)
		assert.HeadCommitMessage(Equals(commitMessage))
		assert.CurrentWindowName("stagingSecondary")

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
