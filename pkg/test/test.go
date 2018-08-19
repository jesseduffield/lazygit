package test

import (
	"errors"
	"os"
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// GenerateRepo generates a repo from test/repos and changes the directory to be
// inside the newly made repo
func GenerateRepo(filename string) error {
	testPath := utils.GetProjectRoot() + "/test/repos/"
	if err := os.Chdir(testPath); err != nil {
		return err
	}
	if output, err := exec.Command("bash", filename).CombinedOutput(); err != nil {
		return errors.New(string(output))
	}
	if err := os.Chdir(testPath + "repo"); err != nil {
		return err
	}
	return nil
}
