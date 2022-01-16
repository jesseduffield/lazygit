package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
)

type FileHelper struct {
	c   *controllers.ControllerCommon
	git *commands.GitCommand
	os  *oscommands.OSCommand
}

func NewFileHelper(
	c *controllers.ControllerCommon,
	git *commands.GitCommand,
	os *oscommands.OSCommand,
) *FileHelper {
	return &FileHelper{
		c:   c,
		git: git,
		os:  os,
	}
}

var _ controllers.IFileHelper = &FileHelper{}

func (self *FileHelper) EditFile(filename string) error {
	return self.EditFileAtLine(filename, 1)
}

func (self *FileHelper) EditFileAtLine(filename string, lineNumber int) error {
	cmdStr, err := self.git.File.GetEditCmdStr(filename, lineNumber)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.EditFile)
	return self.c.RunSubprocessAndRefresh(
		self.os.Cmd.NewShell(cmdStr),
	)
}

func (self *FileHelper) OpenFile(filename string) error {
	self.c.LogAction(self.c.Tr.Actions.OpenFile)
	if err := self.os.OpenFile(filename); err != nil {
		return self.c.Error(err)
	}
	return nil
}
