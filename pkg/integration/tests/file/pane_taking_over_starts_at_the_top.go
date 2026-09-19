package file

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// changedLines is a file's worth of numbered lines, prefixed so that each file's diff
// can be told from the other's on screen.
func changedLines(prefix string) string {
	lines := make([]string, 40)
	for i := range lines {
		lines[i] = fmt.Sprintf("%s%02d", prefix, i+1)
	}
	return strings.Join(lines, "\n") + "\n"
}

var PaneTakingOverStartsAtTheTop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A pane taking the main section over shows its diff from the top, rather than at the offset it was left at",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", changedLines("one"))
		shell.CreateFileAndAdd("file2", changedLines("two"))
		shell.Commit("one")

		// One file's changes are unstaged and the other's are staged, so each is shown
		// in a pane of its own — and selecting one after the other hands the section
		// from one pane to the other. Both diffs are longer than the section, so either
		// pane can be scrolled.
		shell.UpdateFile("file1", changedLines("ONE"))
		shell.UpdateFileAndAdd("file2", changedLines("TWO"))
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("file1"))

		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsVisible().
			Title(Equals("Unstaged changes")).
			ScrollWheelDown().
			ScrollWheelDown().
			OriginYAtLeast(1)

		t.Views().Files().NavigateToLine(Contains("file2"))

		t.Views().Main().IsInvisible()
		t.Views().Secondary().
			IsVisible().
			Title(Equals("Staged changes")).
			Content(Contains("+TWO40")).
			OriginY(0).
			ScrollWheelDown().
			ScrollWheelDown().
			OriginYAtLeast(1)

		// Back to the pane that was left scrolled: what it is given is a diff the user
		// hasn't seen there, so it starts at the top like any other.
		t.Views().Files().NavigateToLine(Contains("file1"))

		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsVisible().
			Title(Equals("Unstaged changes")).
			Content(Contains("+ONE40")).
			OriginY(0)

		t.Views().Files().NavigateToLine(Contains("file2"))

		t.Views().Main().IsInvisible()
		t.Views().Secondary().
			IsVisible().
			Title(Equals("Staged changes")).
			Content(Contains("+TWO40")).
			OriginY(0)
	},
})
