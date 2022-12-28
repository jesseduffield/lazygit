package components

import (
	"fmt"
	"os"
)

type FileSystem struct {
	*assertionHelper
}

// This does _not_ check the files panel, it actually checks the filesystem
func (self *FileSystem) PathPresent(path string) {
	self.assertWithRetries(func() (bool, string) {
		_, err := os.Stat(path)
		return err == nil, fmt.Sprintf("Expected path '%s' to exist, but it does not", path)
	})
}

// This does _not_ check the files panel, it actually checks the filesystem
func (self *FileSystem) PathNotPresent(path string) {
	self.assertWithRetries(func() (bool, string) {
		_, err := os.Stat(path)
		return os.IsNotExist(err), fmt.Sprintf("Expected path '%s' to not exist, but it does", path)
	})
}
