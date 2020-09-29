package commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
)

// GetFilesInDiff get the specified commit files
func (c *GitCommand) GetFilesInDiff(from string, to string, reverse bool, patchManager *patch.PatchManager) ([]*models.CommitFile, error) {
	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	filenames, err := c.OSCommand.RunCommandWithOutput("git diff --submodule --no-ext-diff --name-status %s %s %s", reverseFlag, from, to)
	if err != nil {
		return nil, err
	}

	return c.getCommitFilesFromFilenames(filenames, to, patchManager), nil
}

// filenames string is something like "file1\nfile2\nfile3"
func (c *GitCommand) getCommitFilesFromFilenames(filenames string, parent string, patchManager *patch.PatchManager) []*models.CommitFile {
	commitFiles := make([]*models.CommitFile, 0)

	for _, line := range strings.Split(strings.TrimRight(filenames, "\n"), "\n") {
		// typical result looks like 'A my_file' meaning my_file was added
		if line == "" {
			continue
		}
		changeStatus := line[0:1]
		name := line[2:]
		status := patch.UNSELECTED
		if patchManager != nil && patchManager.To == parent {
			status = patchManager.GetFileStatus(name)
		}

		commitFiles = append(commitFiles, &models.CommitFile{
			Parent:       parent,
			Name:         name,
			ChangeStatus: changeStatus,
			PatchStatus:  status,
		})
	}

	return commitFiles
}
