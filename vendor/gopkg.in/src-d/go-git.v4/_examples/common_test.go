package examples

import (
	"flag"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var examplesTest = flag.Bool("examples", false, "run the examples tests")

var defaultURL = "https://github.com/git-fixtures/basic.git"

var args = map[string][]string{
	"branch":      {defaultURL, tempFolder()},
	"checkout":    {defaultURL, tempFolder(), "35e85108805c84807bc66a02d91535e1e24b38b9"},
	"clone":       {defaultURL, tempFolder()},
	"context":     {defaultURL, tempFolder()},
	"commit":      {cloneRepository(defaultURL, tempFolder())},
	"custom_http": {defaultURL},
	"open":        {cloneRepository(defaultURL, tempFolder())},
	"progress":    {defaultURL, tempFolder()},
	"push":        {setEmptyRemote(cloneRepository(defaultURL, tempFolder()))},
	"revision":    {cloneRepository(defaultURL, tempFolder()), "master~2^"},
	"showcase":    {defaultURL, tempFolder()},
	"tag":         {cloneRepository(defaultURL, tempFolder())},
	"pull":        {createRepositoryWithRemote(tempFolder(), defaultURL)},
}

var ignored = map[string]bool{}

var tempFolders = []string{}

func TestExamples(t *testing.T) {
	flag.Parse()
	if !*examplesTest && os.Getenv("CI") == "" {
		t.Skip("skipping examples tests, pass --examples to execute it")
		return
	}

	defer deleteTempFolders()

	examples, err := filepath.Glob(examplesFolder())
	if err != nil {
		t.Errorf("error finding tests: %s", err)
	}

	for _, example := range examples {
		_, name := filepath.Split(filepath.Dir(example))

		if ignored[name] {
			continue
		}

		t.Run(name, func(t *testing.T) {
			testExample(t, name, example)
		})
	}
}

func tempFolder() string {
	path, err := ioutil.TempDir("", "")
	CheckIfError(err)

	tempFolders = append(tempFolders, path)
	return path
}

func packageFolder() string {
	return filepath.Join(
		build.Default.GOPATH,
		"src", "gopkg.in/src-d/go-git.v4",
	)
}

func examplesFolder() string {
	return filepath.Join(
		packageFolder(),
		"_examples", "*", "main.go",
	)
}

func cloneRepository(url, folder string) string {
	cmd := exec.Command("git", "clone", url, folder)
	err := cmd.Run()
	CheckIfError(err)

	return folder
}

func createBareRepository(dir string) string {
	return createRepository(dir, true)
}

func createRepository(dir string, isBare bool) string {
	var cmd *exec.Cmd
	if isBare {
		cmd = exec.Command("git", "init", "--bare", dir)
	} else {
		cmd = exec.Command("git", "init", dir)
	}
	err := cmd.Run()
	CheckIfError(err)

	return dir
}

func createRepositoryWithRemote(local, remote string) string {
	createRepository(local, false)
	addRemote(local, remote)
	return local
}

func setEmptyRemote(dir string) string {
	remote := createBareRepository(tempFolder())
	setRemote(dir, remote)
	return dir
}

func setRemote(local, remote string) {
	cmd := exec.Command("git", "remote", "set-url", "origin", remote)
	cmd.Dir = local
	err := cmd.Run()
	CheckIfError(err)
}

func addRemote(local, remote string) {
	cmd := exec.Command("git", "remote", "add", "origin", remote)
	cmd.Dir = local
	err := cmd.Run()
	CheckIfError(err)
}

func testExample(t *testing.T, name, example string) {
	cmd := exec.Command("go", append([]string{
		"run", filepath.Join(example),
	}, args[name]...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Errorf("error running cmd %q", err)
	}
}

func deleteTempFolders() {
	for _, folder := range tempFolders {
		err := os.RemoveAll(folder)
		CheckIfError(err)
	}
}
