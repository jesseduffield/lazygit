package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// TODO: find out why we can't use assert.SelectedLine() on the staging/stagingSecondary views.

var Unstaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the unstaged files menu, staging a line and committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFile("myfile", "myfile content\nwith a second line").
			CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.Model().CommitCount(0)

		assert.Views().Current().Name("files").SelectedLine(Contains("myfile"))
		input.Enter()
		assert.Views().Current().Name("staging")
		assert.Views().ByName("stagingSecondary").Content(DoesNotContain("+myfile content"))
		// stage the first line
		input.PrimaryAction()
		assert.Views().ByName("staging").Content(DoesNotContain("+myfile content"))
		assert.Views().ByName("stagingSecondary").Content(Contains("+myfile content"))

		input.Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		input.CommitMessagePanel().Type(commitMessage).Confirm()

		assert.Model().CommitCount(1)
		assert.Model().HeadCommitMessage(Equals(commitMessage))
		assert.Views().Current().Name("staging")

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
