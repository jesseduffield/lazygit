package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenameSimilarityThresholdChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Change the rename similarity threshold while in the commits panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("add original")

		shell.DeleteFileAndAdd("original")
		shell.CreateFileAndAdd("renamed", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("change name and contents")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus()

		t.Views().Main().
			ContainsLines(
				Contains("2 files changed, 10 insertions(+), 5 deletions(-)"),
			)

		t.Views().Commits().
			Press(keys.Universal.DecreaseRenameSimilarityThreshold).
			Tap(func() {
				t.ExpectToast(Equals("Changed rename similarity threshold to 45%"))
			})

		t.Views().Main().
			ContainsLines(
				Contains("original => renamed"),
				Contains("1 file changed, 5 insertions(+)"),
			)
	},
})
