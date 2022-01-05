package loaders

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type CommitFileLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewCommitFileLoader(common *common.Common, cmd oscommands.ICmdObjBuilder) *CommitFileLoader {
	return &CommitFileLoader{
		Common: common,
		cmd:    cmd,
	}
}

// GetFilesInDiff get the specified commit files
func (self *CommitFileLoader) GetFilesInDiff(from string, to string, reverse bool) ([]*models.CommitFile, error) {
	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	filenames, err := self.cmd.New(fmt.Sprintf("git diff --submodule --no-ext-diff --name-status -z --no-renames %s %s %s", reverseFlag, from, to)).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	return self.getCommitFilesFromFilenames(filenames), nil
}

// filenames string is something like "file1\nfile2\nfile3"
func (self *CommitFileLoader) getCommitFilesFromFilenames(filenames string) []*models.CommitFile {
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
