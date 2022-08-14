package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IFilesHelper interface {
	EditFile(filename string) error
	EditFileAtLine(filename string, lineNumber int) error
	OpenFile(filename string) error
	OpenFileAtLine(filename string, lineNumber int) error
}

type FilesHelper struct {
	c   *types.HelperCommon
	git *commands.GitCommand
	os  *oscommands.OSCommand
}

func NewFilesHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
	os *oscommands.OSCommand,
) *FilesHelper {
	return &FilesHelper{
		c:   c,
		git: git,
		os:  os,
	}
}

var _ IFilesHelper = &FilesHelper{}

func (self *FilesHelper) EditFile(filename string) error {
	return self.EditFileAtLine(filename, 1)
}

func (self *FilesHelper) EditFileAtLine(filename string, lineNumber int) error {
	cmdStr, err := self.git.File.GetEditCmdStr(filename, lineNumber)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.EditFile)
	return self.c.RunSubprocessAndRefresh(
		self.os.Cmd.NewShell(cmdStr),
	)
}

func (self *FilesHelper) OpenFile(filename string) error {
	return self.OpenFileAtLine(filename, 1)
}

func (self *FilesHelper) OpenFileAtLine(filename string, lineNumber int) error {
	self.c.LogAction(self.c.Tr.Actions.OpenFile)
	if err := self.os.OpenFileAtLine(filename, lineNumber); err != nil {
		return self.c.Error(err)
	}
	return nil
}
