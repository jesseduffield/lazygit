package ui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RangeSelectWithAutoscroll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Keep scrolling while creating a range selection at the panel edge",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(40)
		fileContent := "base\n"
		shell.CreateFileAndAdd("file1", fileContent)
		for i := 1; i <= 40; i++ {
			fileContent += fmt.Sprintf("line %d\n", i)
		}
		shell.UpdateFile("file1", fileContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus()
		t.Views().Commits().
			ClickAndHold(1, 0).
			MouseMoveToBottom(1).
			OriginYAtLeast(3).
			SelectedLineIdxAtLeast(3).
			MouseRelease()

		t.Views().Files().
			Focus().
			PressEnter()
		t.Views().Staging().
			ClickAndHold(1, 6).
			MouseMoveToBottom(1).
			OriginYAtLeast(3).
			SelectedLineIdxAtLeast(9).
			MouseRelease()
	},
})
