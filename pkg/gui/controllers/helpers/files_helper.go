package helpers

import (
	"path/filepath"

	"github.com/samber/lo"
)

type FilesHelper struct {
	c *HelperCommon
}

func NewFilesHelper(c *HelperCommon) *FilesHelper {
	return &FilesHelper{
		c: c,
	}
}

func (self *FilesHelper) EditFiles(filenames []string) error {
	absPaths := lo.Map(filenames, func(filename string, _ int) string {
		absPath, err := filepath.Abs(filename)
		if err != nil {
			return filename
		}
		return absPath
	})
	cmdStr, suspend := self.c.Git().File.GetEditCmdStr(absPaths)
	return self.callEditor(cmdStr, suspend)
}

func (self *FilesHelper) EditFileAtLine(filename string, lineNumber int) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	cmdStr, suspend := self.c.Git().File.GetEditAtLineCmdStr(absPath, lineNumber)
	return self.callEditor(cmdStr, suspend)
}

func (self *FilesHelper) EditFileAtLineAndWait(filename string, lineNumber int) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	cmdStr := self.c.Git().File.GetEditAtLineAndWaitCmdStr(absPath, lineNumber)

	// Always suspend, regardless of the value of the suspend config,
	// since we want to prevent interacting with the UI until the editor
	// returns, even if the editor doesn't use the terminal.
	//
	// And nothing is refreshed here; that is why this doesn't go through
	// callEditor. The editor was handed a patch we wrote for it, so the repo
	// hasn't changed when it returns; it changes when the caller applies what
	// came back. A refresh in between reads the state from before that, and
	// then races the caller's own refresh to publish it.
	_, err = self.c.RunSubprocess(
		self.c.OS().Cmd.NewShell(cmdStr, self.c.UserConfig().OS.ShellFunctionsFile))
	return err
}

func (self *FilesHelper) OpenDirInEditor(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	cmdStr, suspend := self.c.Git().File.GetOpenDirInEditorCmdStr(absPath)

	return self.callEditor(cmdStr, suspend)
}

func (self *FilesHelper) callEditor(cmdStr string, suspend bool) error {
	if suspend {
		return self.c.RunSubprocessAndRefresh(
			self.c.OS().Cmd.NewShell(cmdStr, self.c.UserConfig().OS.ShellFunctionsFile),
		)
	}

	return self.c.OS().Cmd.NewShell(cmdStr, self.c.UserConfig().OS.ShellFunctionsFile).Run()
}

func (self *FilesHelper) OpenFile(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.OpenFile)
	if err := self.c.OS().OpenFile(absPath); err != nil {
		return err
	}
	return nil
}
