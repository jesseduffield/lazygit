package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GitignoreSpecialCharacters = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignore files with special characters in their names",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".gitignore", "")
		shell.CreateFile("#file", "")
		shell.CreateFile("file#abc", "")
		shell.CreateFile("!file", "")
		shell.CreateFile("file!abc", "")
		shell.CreateFile("abc*def", "")
		shell.CreateFile("abc_def", "")
		shell.CreateFile("file[x]", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		excludeFile := func(fileName string) {
			t.Views().Files().
				NavigateToLine(Contains(fileName)).
				Press(keys.Files.IgnoreFile)

			t.ExpectPopup().Menu().
				Title(Equals("Ignore or exclude file")).
				Select(Contains("Add to .gitignore")).
				Confirm()
		}

		t.Views().Files().
			Focus().
			Lines(
				Equals("▼ /"),
				Equals("  ?? !file"),
				Equals("  ?? #file"),
				Equals("  ?? .gitignore"),
				Equals("  ?? abc*def"),
				Equals("  ?? abc_def"),
				Equals("  ?? file!abc"),
				Equals("  ?? file#abc"),
				Equals("  ?? file[x]"),
			)

		excludeFile("#file")
		excludeFile("file#abc")
		excludeFile("!file")
		excludeFile("file!abc")
		excludeFile("abc*def")
		excludeFile("file[x]")

		t.Views().Files().
			Lines(
				Equals("▼ /"),
				Equals("  ?? .gitignore"),
				Equals("  ?? abc_def"),
			)

		t.FileSystem().FileContent(".gitignore", Equals("\\#file\nfile#abc\n\\!file\nfile!abc\nabc\\*def\nfile\\[x\\]\n"))
	},
})
