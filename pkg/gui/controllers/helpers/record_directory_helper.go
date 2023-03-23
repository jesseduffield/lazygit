package helpers

import (
	"os"
)

type RecordDirectoryHelper struct {
	c *HelperCommon
}

func NewRecordDirectoryHelper(c *HelperCommon) *RecordDirectoryHelper {
	return &RecordDirectoryHelper{
		c: c,
	}
}

// when a user runs lazygit with the LAZYGIT_NEW_DIR_FILE env variable defined
// we will write the current directory to that file on exit so that their
// shell can then change to that directory. That means you don't get kicked
// back to the directory that you started with.
func (self *RecordDirectoryHelper) RecordCurrentDirectory() error {
	// determine current directory, set it in LAZYGIT_NEW_DIR_FILE
	dirName, err := os.Getwd()
	if err != nil {
		return err
	}
	return self.RecordDirectory(dirName)
}

func (self *RecordDirectoryHelper) RecordDirectory(dirName string) error {
	newDirFilePath := os.Getenv("LAZYGIT_NEW_DIR_FILE")
	if newDirFilePath == "" {
		return nil
	}
	return self.c.OS().CreateFileWithContent(newDirFilePath, dirName)
}
