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
	cmdStr, editInTerminal := self.git.File.GetEditCmdStr(filename)
	return self.callEditor(cmdStr, editInTerminal)
}

func (self *FilesHelper) EditFileAtLine(filename string, lineNumber int) error {
	cmdStr, editInTerminal := self.git.File.GetEditAtLineCmdStr(filename, lineNumber)
	return self.callEditor(cmdStr, editInTerminal)
}

func (self *FilesHelper) EditFileAtLineAndWait(filename string, lineNumber int) error {
	cmdStr := self.git.File.GetEditAtLineAndWaitCmdStr(filename, lineNumber)

	// Always suspend, regardless of the value of the editInTerminal config,
	// since we want to prevent interacting with the UI until the editor
	// returns, even if the editor doesn't use the terminal
	return self.callEditor(cmdStr, true)
}

func (self *FilesHelper) callEditor(cmdStr string, editInTerminal bool) error {
	if editInTerminal {
		return self.c.RunSubprocessAndRefresh(
			self.os.Cmd.NewShell(cmdStr),
		)
	}

	return self.os.Cmd.NewShell(cmdStr).Run()
}

func (self *FilesHelper) OpenFile(filename string) error {
	self.c.LogAction(self.c.Tr.Actions.OpenFile)
	if err := self.os.OpenFile(filename); err != nil {
		return self.c.Error(err)
	}
	return nil
}
