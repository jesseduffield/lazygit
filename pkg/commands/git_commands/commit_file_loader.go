package git_commands

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
	cmdArgs := NewGitCmd("diff").
		Config("diff.noprefix=false").
		Arg("--submodule").
		Arg("--no-ext-diff").
		Arg("--name-status").
		Arg("-z").
		Arg(fmt.Sprintf("--find-renames=%d%%", self.UserConfig().Git.RenameSimilarityThreshold)).
		ArgIf(reverse, "-R").
		Arg(from).
		Arg(to).
		ToArgv()

	filenames, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	return getCommitFilesFromFilenames(filenames), nil
}

// filenames string is something like "MM\x00file1\x00MU\x00file2\x00AA\x00file3\x00"
// so we need to split it by the null character and then map each status-name pair
// to a commit file. Renames (and copies) are special: their status is followed by
// two paths (the old one and the new one) rather than one, e.g.
// "R100\x00old\x00new\x00".
func getCommitFilesFromFilenames(filenames string) []*models.CommitFile {
	fields := strings.Split(strings.TrimRight(filenames, "\x00"), "\x00")
	if len(fields) == 1 {
		return []*models.CommitFile{}
	}

	commitFiles := make([]*models.CommitFile, 0, len(fields)/2)
	for i := 0; i < len(fields)-1; {
		changeStatus := fields[i]
		if changeStatus[0] == 'R' || changeStatus[0] == 'C' {
			// The status has a similarity score appended (e.g. "R100"); drop it
			// so the rest of the code only has to deal with a plain "R" or "C".
			commitFiles = append(commitFiles, &models.CommitFile{
				ChangeStatus: changeStatus[:1],
				PreviousPath: fields[i+1],
				Path:         fields[i+2],
			})
			i += 3
		} else {
			// typical result looks like 'A my_file' meaning my_file was added
			commitFiles = append(commitFiles, &models.CommitFile{
				ChangeStatus: changeStatus,
				Path:         fields[i+1],
			})
			i += 2
		}
	}

	return commitFiles
}
