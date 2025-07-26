package helpers

import (
	"os"
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
	// returns, even if the editor doesn't use the terminal
	return self.callEditor(cmdStr, true)
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

func (self *FilesHelper) EditFileAtRevision(filename string, commitHash string) error {
	// Get the file content from the specific commit
	content, err := self.c.Git().Commit.ShowFileContentCmdObj(commitHash, filename).RunWithOutput()
	if err != nil {
		return err
	}

	// Create a temporary file with the same extension for proper syntax highlighting
	ext := filepath.Ext(filename)
	basename := filepath.Base(filename)
	tempFile, err := os.CreateTemp("", basename+"_*"+ext)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Write the historical content to the temp file
	if _, err := tempFile.WriteString(content); err != nil {
		os.Remove(tempFile.Name())
		return err
	}

	// Get the editor command for the temp file
	cmdStr, suspend := self.c.Git().File.GetEditCmdStr([]string{tempFile.Name()})
	
	// Log the action
	self.c.LogAction(self.c.Tr.Actions.OpenFileAtRevision)

	return self.callEditor(cmdStr, suspend)
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
