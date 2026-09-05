package filter_and_search

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// longFileWithThreeMatches is long enough that a render of its diff stops well short
// of the end, with two of the three matches for the search below the point it stops at.
func longFileWithThreeMatches() string {
	lines := make([]string, 0, 2000)
	for i := range 2000 {
		switch i {
		case 100:
			lines = append(lines, "NEEDLE first")
		case 1000:
			lines = append(lines, "NEEDLE middle")
		case 1900:
			lines = append(lines, "NEEDLE last")
		default:
			lines = append(lines, fmt.Sprintf("line %d", i))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

var SearchALongDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Search a diff that is longer than a single render of it reads",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A small window, so that a render stops well short of 2000 lines.
	Width:       120,
	Height:      30,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("one")

		shell.UpdateFile("file1", longFileWithThreeMatches())
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// All three matches are counted: opening the prompt reads the whole diff
		// first, however much of it the render had got to.
		t.Views().Main().
			IsFocused().
			FilterOrSearch("NEEDLE")

		t.Views().Search().Content(Contains("matches for 'NEEDLE' (1 of 3)"))
	},
})
