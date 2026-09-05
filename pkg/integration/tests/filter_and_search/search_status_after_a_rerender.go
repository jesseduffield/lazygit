package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchStatusAfterARerender = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The search status counts the matches in a diff that has been rendered again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1",
			"line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nNEEDLE\nline 11\nline 12\nline 13\nline 14\n")
		shell.Commit("one")

		// Four lines above NEEDLE, so that it is context at a context size of 4 but
		// not at 3.
		shell.UpdateFile("file1",
			"line 1\nline 2\nline 3\nline 4\nline 5\nchanged\nline 7\nline 8\nline 9\nNEEDLE\nline 11\nline 12\nline 13\nline 14\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("NEEDLE")).
			FilterOrSearch("NEEDLE")

		t.Views().Search().Content(Contains("No matches for 'NEEDLE'"))

		// A wider context brings NEEDLE into the diff, and the search counts it.
		t.Views().Main().
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			Content(Contains("NEEDLE"))

		t.Views().Search().Content(Contains("matches for 'NEEDLE' (1 of 1)"))
	},
})
