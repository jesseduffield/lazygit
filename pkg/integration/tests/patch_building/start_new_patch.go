package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartNewPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Attempt to add a file from another commit to a patch, then agree to start a new patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.Commit("first commit")

		shell.CreateFileAndAdd("file2", "file2 content")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file2").IsSelected(),
			).
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				t.Views().Secondary().Content(Contains("file2"))
			}).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("first commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Contains("Discard patch")).
					Content(Contains("You can only build a patch from one commit/stash-entry at a time. Discard current patch?")).
					Confirm()

				t.Views().Secondary().Content(Contains("file1").DoesNotContain("file2"))
			})
	},
})
