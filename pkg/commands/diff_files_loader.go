package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type DiffFilesLoader struct {
	*MgrCtx
}

func NewDiffFilesLoader(mgrCtx *MgrCtx) *DiffFilesLoader {
	return &DiffFilesLoader{
		MgrCtx: mgrCtx,
	}
}

// Load gets the specified files from a diff
func (c *DiffFilesLoader) Load(from string, to string, reverse bool) ([]*models.CommitFile, error) {
	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	filenames, err := c.RunWithOutput(
		BuildGitCmdObjFromStr(
			fmt.Sprintf("diff --submodule --no-ext-diff --name-status -z --no-renames %s %s %s", reverseFlag, from, to),
		),
	)
	if err != nil {
		return nil, err
	}

	return c.getCommitFilesFromFilenames(filenames), nil
}

// filenames string is something like "file1\nfile2\nfile3"
func (c *DiffFilesLoader) getCommitFilesFromFilenames(filenames string) []*models.CommitFile {
	commitFiles := make([]*models.CommitFile, 0)

	lines := strings.Split(strings.TrimRight(filenames, "\x00"), "\x00")
	n := len(lines)
	for i := 0; i < n-1; i += 2 {
		// typical result looks like 'A my_file' meaning my_file was added
		changeStatus := lines[i]
		name := lines[i+1]

		commitFiles = append(commitFiles, &models.CommitFile{
			Name:         name,
			ChangeStatus: changeStatus,
		})
	}

	return commitFiles
}
