package git_commands

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ArchiveCommands struct {
	*GitCommon
}

func NewArchiveCommands(gitCommon *GitCommon) *ArchiveCommands {
	return &ArchiveCommands{
		GitCommon: gitCommon,
	}
}

func (self *ArchiveCommands) Archive(refName string, archiveName string, prefix string) error {
	cmdArgs := NewGitCmd("archive").ArgIf(prefix != "", "--prefix", prefix).Arg("-o", archiveName, refName)

	return self.cmd.New(cmdArgs.ToArgv()).DontLog().Run()
}

// GetValidArchiveFormats returns a slice of valid archive formats from git.
//
// For example, telling git how to handle archiving to .tar.bz2 using
// git config tar.tar.bz2.command="bzip2" will add "tar.bz2" to git's
// list of archival formats.
func (self *ArchiveCommands) GetValidArchiveFormats() ([]string, error) {
	cmdArgs := NewGitCmd("archive").Arg("-l").ToArgv()

	validExtensions, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return []string{""}, err
	}

	validExtensionArray := utils.SplitLines(validExtensions)

	for k, v := range validExtensionArray {
		validExtensionArray[k] = "." + v
	}

	return validExtensionArray, nil
}
