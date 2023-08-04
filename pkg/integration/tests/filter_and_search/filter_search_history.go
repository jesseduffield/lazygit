package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterSearchHistory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Navigating search history",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			// populate search history with some values
			FilterOrSearch("1").
			FilterOrSearch("2").
			FilterOrSearch("3").
			Press(keys.Universal.StartSearch).
			// clear initial search value
			Tap(func() {
				t.ExpectSearch().Clear()
			}).
			// test main search history functionality
			Tap(func() {
				t.Views().Search().
					Press(keys.Universal.PrevItem).
					Content(Contains("3")).
					Press(keys.Universal.PrevItem).
					Content(Contains("2")).
					Press(keys.Universal.PrevItem).
					Content(Contains("1")).
					Press(keys.Universal.PrevItem).
					Content(Contains("1")).
					Press(keys.Universal.NextItem).
					Content(Contains("2")).
					Press(keys.Universal.NextItem).
					Content(Contains("3")).
					Press(keys.Universal.NextItem).
					Content(Contains("")).
					Press(keys.Universal.NextItem).
					Content(Contains("")).
					Press(keys.Universal.PrevItem).
					Content(Contains("3")).
					PressEscape()
			}).
			// test that it resets after you enter and exit a search
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.Views().Search().
					Press(keys.Universal.PrevItem).
					Content(Contains("3")).
					PressEscape()
			})

		// test that the histories are separate for each view
		t.Views().Commits().
			Focus().
			FilterOrSearch("a").
			FilterOrSearch("b").
			FilterOrSearch("c").
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().Clear()
			}).
			Tap(func() {
				t.Views().Search().
					Press(keys.Universal.PrevItem).
					Content(Contains("c")).
					Press(keys.Universal.PrevItem).
					Content(Contains("b")).
					Press(keys.Universal.PrevItem).
					Content(Contains("a"))
			})
	},
})
