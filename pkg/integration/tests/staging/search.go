package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Search = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the search feature in the staging panel",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file1", "one\ntwo\nthree\nfour\nfive")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().
					Type("four").
					Confirm()

				t.Views().Search().Content(Contains("matches for 'four' (1 of 1)"))
			}).
			SelectedLine(Contains("+four")). // stage the line
			PressPrimaryAction().
			Content(DoesNotContain("+four")).
			Tap(func() {
				t.Views().StagingSecondary().
					Content(Contains("+four"))
			})
	},
})
