package cli_file_flag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Basic = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open file history using --file flag",
	ExtraCmdArgs: []string{"--file=src/main.go"},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		// Create directory structure
		shell.RunCommand([]string{"mkdir", "-p", "src"})
		
		// Create file with initial content
		shell.CreateFileAndAdd("src/main.go", "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}")
		shell.Commit("initial main.go")
		
		// Create unrelated file
		shell.CreateFileAndAdd("other.txt", "unrelated content")
		shell.Commit("add unrelated file")
		
		// Update main.go
		shell.UpdateFileAndAdd("src/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}")
		shell.Commit("update main.go with fmt")
		
		// Another unrelated commit
		shell.UpdateFileAndAdd("other.txt", "more unrelated content")
		shell.Commit("update other file")
		
		// Another main.go update
		shell.UpdateFileAndAdd("src/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, lazygit!\")\n}")
		shell.Commit("update main.go message")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Should open directly to commits view with filtering active
		t.Views().Information().Content(Contains("Filtering by 'src/main.go'"))
		
		// Should be focused on commits view
		t.Views().Commits().IsFocused()
		
		// Should show only commits that touch src/main.go
		t.Views().Commits().
			Lines(
				Contains("update main.go message").IsSelected(),
				Contains("update main.go with fmt"),
				Contains("initial main.go"),
			)
		
		// Should NOT show unrelated commits
		t.Views().Commits().Content(DoesNotContain("add unrelated file"))
		t.Views().Commits().Content(DoesNotContain("update other file"))
		
		// Main view should show the diff for the selected commit
		t.Views().Main().
			ContainsLines(
				Contains("update main.go message"),
				Contains("src/main.go"),
			)
	},
})