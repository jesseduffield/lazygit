package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RerenderTheSearchedMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A refresh renders the focused main view again even while it is being searched",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nNEEDLE\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			FilterOrSearch("NEEDLE").
			Content(Contains("+NEEDLE")).
			Tap(func() {
				t.Shell().UpdateFile("file1", "one\nOTHER\nthree\n")
			}).
			Press(keys.Universal.Refresh).
			Content(Contains("+OTHER")).
			Content(DoesNotContain("+NEEDLE"))

		t.Views().Search().Content(Contains("No matches for 'NEEDLE'"))
	},
})
