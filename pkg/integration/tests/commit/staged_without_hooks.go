package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagedWithoutHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing without pre-commit hooks",
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

		// stage the file
		assert.CurrentViewName("files")
		assert.CurrentLine(Contains("myfile"))
		input.PrimaryAction()
		input.Enter()
		assert.CurrentViewName("stagingSecondary")
		// we start with both lines having been staged
		assert.ViewContent("stagingSecondary", Contains("+myfile content"))
		assert.ViewContent("stagingSecondary", Contains("+with a second line"))
		assert.ViewContent("staging", NotContains("+myfile content"))
		assert.ViewContent("staging", NotContains("+with a second line"))

		// unstage the selected line
		input.PrimaryAction()

		// the line should have been moved to the main view
		assert.ViewContent("stagingSecondary", NotContains("+myfile content"))
		assert.ViewContent("stagingSecondary", Contains("+with a second line"))
		assert.ViewContent("staging", Contains("+myfile content"))
		assert.ViewContent("staging", NotContains("+with a second line"))

		input.Press(keys.Files.CommitChangesWithoutHook)
		assert.InCommitMessagePanel()
		assert.CurrentViewContent(Contains("WIP"))
		commitMessage := ": my commit message"
		input.Type(commitMessage)
		input.Confirm()

		assert.CommitCount(1)
		assert.HeadCommitMessage(Equals("WIP" + commitMessage))
		assert.CurrentViewName("stagingSecondary")

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
