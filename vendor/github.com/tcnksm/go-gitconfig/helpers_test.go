package gitconfig

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func withGlobalGitConfigFile(content string) func() {
	tmpdir, err := ioutil.TempDir("", "go-gitconfig-test")
	if err != nil {
		panic(err)
	}

	tmpGitConfigFile := filepath.Join(tmpdir, ".gitconfig")

	ioutil.WriteFile(
		tmpGitConfigFile,
		[]byte(content),
		0777,
	)

	prevGitConfigEnv := os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)

	return func() {
		os.Setenv("HOME", prevGitConfigEnv)
	}
}

func includeGitConfigFile(content string) string {
	tmpdir, err := ioutil.TempDir("", "go-gitconfig-test")
	if err != nil {
		panic(err)
	}

	tmpGitIncludeConfigFile := filepath.Join(tmpdir, ".gitconfig.local")
	ioutil.WriteFile(
		tmpGitIncludeConfigFile,
		[]byte(content),
		0777,
	)

	return tmpGitIncludeConfigFile
}

func withLocalGitConfigFile(key string, value string) func() {
	var err error
	tmpdir, err := ioutil.TempDir(".", "go-gitconfig-test")
	if err != nil {
		panic(err)
	}

	prevDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	os.Chdir(tmpdir)

	gitInit := exec.Command("git", "init")
	gitInit.Stderr = ioutil.Discard
	if err = gitInit.Run(); err != nil {
		panic(err)
	}

	gitAddConfig := exec.Command("git", "config", "--local", key, value)
	gitAddConfig.Stderr = ioutil.Discard
	if err = gitAddConfig.Run(); err != nil {
		panic(err)
	}

	return func() {
		os.Chdir(prevDir)
		os.RemoveAll(tmpdir)
	}
}
