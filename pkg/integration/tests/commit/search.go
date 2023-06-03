package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Search = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Search for a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.EmptyCommit("four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().
					Type("two").
					Confirm()

				t.Views().Search().Content(Contains("matches for 'two' (1 of 1)"))
			}).
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().
					Clear().
					Type("o").
					Confirm()

				t.Views().Search().Content(Contains("matches for 'o' (2 of 3)"))
			}).
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press("n").
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'o' (3 of 3)"))
			}).
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one").IsSelected(),
			).
			Press("n").
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'o' (1 of 3)"))
			}).
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press("n").
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'o' (2 of 3)"))
			}).
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two").IsSelected(),
				Contains("one"),
			).
			Press("N").
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'o' (1 of 3)"))
			}).
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			Press("N").
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'o' (3 of 3)"))
			}).
			Lines(
				Contains("four"),
				Contains("three"),
				Contains("two"),
				Contains("one").IsSelected(),
			)
	},
})
