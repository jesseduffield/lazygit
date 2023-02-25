package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetWithEscape = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset a custom patch with the escape keybinding",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("building patch"))
			}).
			PressEscape()

		// hitting escape at the top level will reset the patch
		t.Views().Commits().
			IsFocused().
			PressEscape()

		t.Views().Information().Content(DoesNotContain("building patch"))
	},
})
