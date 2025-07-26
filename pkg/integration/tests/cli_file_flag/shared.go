package cli_file_flag

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func createFileWithHistory(shell *Shell, filePath, initialContent string) {
	shell.CreateFileAndAdd(filePath, initialContent)
	shell.Commit("add " + filePath)
}

func updateFileWithHistory(shell *Shell, filePath, newContent, commitMessage string) {
	shell.UpdateFileAndAdd(filePath, newContent)
	shell.Commit(commitMessage)
}

func validateFileFilteringActive(t *TestDriver, filePath string) {
	t.Views().Information().Content(Contains("Filtering by '" + filePath + "'"))
	t.Views().Commits().IsFocused()
}