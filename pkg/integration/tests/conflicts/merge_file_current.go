package conflicts

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

func testDataCurrent() (original, current, incoming, final string) {
	original = `
1
2
3
4
5
6
`
	current = `
1
2
3
4
5a
6
`
	incoming = `
1b
2
3
4
5b
6
`
	final = `
1b
2
3
4
5a
6
`
	return original, current, incoming, final
}

var MergeFileCurrent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Conflicting file can be resolved to 'ours' (current changes) version via merge-file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		original, current, incoming, _ := testDataCurrent()
		shared.CreateMergeConflictFileForMergeFileTests(shell, original, current, incoming)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		_, _, _, expected := testDataCurrent()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file").IsSelected(),
			)

		t.GlobalPress(keys.Files.OpenMergeOptions)

		t.ExpectPopup().Menu().
			Title(Equals("Resolve merge conflicts")).
			Select(Contains("Use current changes")). // merge-file --ours
			Confirm()

		t.Common().ContinueOnConflictsResolved("merge")

		t.Views().Files().IsEmpty()

		t.FileSystem().FileContent("file", Equals(expected))
	},
})
