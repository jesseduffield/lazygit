package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var HideEmptyStagedPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Hide the empty staged panel and show it as changes are staged",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.HideEmptyStagedPanel = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit("initial")
		shell.UpdateFile("file1", "one\ntwo\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().PressEnter()
		t.Views().Staging().IsFocused().ContainsLines(Contains("+two"))
		t.Views().StagingSecondary().IsInvisible()

		t.Views().Staging().Press(keys.Universal.TogglePanel).IsFocused().PressPrimaryAction()
		t.Views().StagingSecondary().IsVisible().IsFocused().ContainsLines(Contains("+two")).PressPrimaryAction()
		t.Views().Staging().IsFocused().ContainsLines(Contains("+two"))
		t.Views().StagingSecondary().IsInvisible()

		t.Views().Staging().PressPrimaryAction()
		t.Views().StagingSecondary().IsVisible().IsFocused().ContainsLines(Contains("+two"))
	},
})
