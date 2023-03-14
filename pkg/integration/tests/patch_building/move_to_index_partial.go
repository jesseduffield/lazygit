package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToIndexPartial = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to the index. This is different from the MoveToIndex test in that we're only selecting a partial patch from a file",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "first line2\nsecond line\nthird line2\n")
		shell.Commit("second commit")

		shell.CreateFileAndAdd("file2", "file1 content")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("second commit"),
				Contains("first commit"),
			).
			NavigateToLine(Contains("second commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			ContainsLines(
				Contains(`-first line`).IsSelected(),
				Contains(`+first line2`),
				Contains(` second line`),
				Contains(`-third line`),
				Contains(`+third line2`),
			).
			PressPrimaryAction().
			SelectNextItem().
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("building patch"))

				t.Views().PatchBuildingSecondary().
					ContainsLines(
						Contains(`-first line`),
						Contains(`+first line2`),
						Contains(` second line`),
						Contains(` third line`),
					)

				t.Common().SelectPatchOption(Contains("move patch out into index"))

				t.Views().Files().
					Lines(
						Contains("M").Contains("file1"),
					)
			})

		// Focus is automatically returned to the commit files panel. Arguably it shouldn't be.
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1"),
			)

		t.Views().Main().
			ContainsLines(
				Contains(` first line`),
				Contains(` second line`),
				Contains(`-third line`),
				Contains(`+third line2`),
			)

		t.Views().Files().
			Focus()

		t.Views().Main().
			ContainsLines(
				Contains(`-first line`),
				Contains(`+first line2`),
				Contains(` second line`),
				Contains(` third line2`),
			)
	},
})
