package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRangeOfLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a range of lines",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "1st\n2nd\n3rd\n4th\n5th\n6th\n")
		shell.Commit("Add file")
		shell.UpdateFile("myfile", "1st changed\n2nd changed\n3rd\n4th\n5th changed\n6th\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			PressEnter()

		t.Views().Staging().
			Content(
				Contains("-1st\n-2nd\n+1st changed\n+2nd changed\n 3rd\n 4th\n-5th\n+5th changed\n 6th"),
			).
			SelectedLine(Equals("-1st")).
			Press(keys.Main.ToggleDragSelect).
			SelectNextItem().
			SelectNextItem().
			SelectNextItem().
			SelectNextItem().
			PressPrimaryAction().
			Content(
				Contains(" 3rd\n 4th\n-5th\n+5th changed\n 6th"),
			).
			SelectedLine(Equals("-5th"))
	},
})
