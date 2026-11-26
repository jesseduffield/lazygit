package file

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagedChangesInLowerPane = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A file's staged changes are shown in the lower pane whether or not it also has unstaged ones",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("both", "one\n")
		shell.CreateFileAndAdd("indexOnly", "one\n")
		shell.CreateFileAndAdd("worktreeOnly", "one\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("both", "one\nstaged\n")
		shell.UpdateFile("both", "one\nstaged\nunstaged\n")
		// More staged lines than fit in the pane, so that it can be scrolled.
		staged := make([]string, 40)
		for i := range staged {
			staged[i] = fmt.Sprintf("staged%02d", i+1)
		}
		shell.UpdateFileAndAdd("indexOnly", "one\n"+strings.Join(staged, "\n")+"\n")
		shell.UpdateFile("worktreeOnly", "one\nunstaged\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("both"))

		// With changes on both sides, each side has its own pane.
		t.Views().Main().
			Title(Equals("Unstaged changes")).
			Content(Contains("+unstaged"))
		t.Views().Secondary().
			Title(Equals("Staged changes")).
			Content(Contains("+staged"))

		// With nothing unstaged, the staged side keeps its pane, which then has the
		// whole space to itself.
		t.Views().Files().NavigateToLine(Contains("indexOnly"))

		t.Views().Main().IsInvisible()
		t.Views().Secondary().
			IsVisible().
			Title(Equals("Staged changes")).
			Content(Contains("+staged01")).
			// The key that focuses the diff wears its label, wherever the diff is.
			TitlePrefix(Equals("[0]")).
			OriginY(0)

		// And the keys for scrolling the diff scroll the pane it is in.
		t.GlobalPress(keys.Universal.ScrollDownMain)
		t.Views().Secondary().OriginYAtLeast(1)
		t.GlobalPress(keys.Universal.ScrollUpMain)
		t.Views().Secondary().OriginY(0)

		// And focusing the diff focuses the pane it is in.
		t.Views().Files().Press(keys.Universal.FocusMainView)
		t.Views().Secondary().IsFocused()
		t.Views().Secondary().PressEscape()

		// With nothing staged, only the upper pane is shown.
		t.Views().Files().
			IsFocused().
			NavigateToLine(Contains("worktreeOnly"))

		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsVisible().
			Title(Equals("Unstaged changes")).
			Content(Contains("+unstaged")).
			TitlePrefix(Equals("[0]"))
	},
})
