package file

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PaneShownAgainStartsAtTheTop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A pane that was emptied while it wasn't shown starts at the top when it comes back, rather than where it was left",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.CreateFileAndAdd("file2", "one\n")
		shell.Commit("one")

		// More staged changes in file1 than fit in the pane they are shown in, so that
		// there is a position in it to lose, plus an unstaged change to give the file a
		// second pane.
		for i := range lines {
			lines[i] = strings.ToUpper(lines[i])
		}
		shell.UpdateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"unstaged\n")

		shell.UpdateFile("file2", "two\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("file1"))

		t.Views().Secondary().
			IsVisible().
			Title(Equals("Staged changes")).
			ScrollWheelDown().
			ScrollWheelDown().
			OriginYAtLeast(1)

		// A file with nothing staged leaves that pane with nothing to show, so it goes
		// away and is emptied.
		t.Views().Files().NavigateToLine(Contains("file2"))
		t.Views().Secondary().IsInvisible()

		t.Views().Files().NavigateToLine(Contains("file1"))
		t.Views().Secondary().
			IsVisible().
			Content(Contains("+LINE40")).
			OriginY(0)
	},
})
