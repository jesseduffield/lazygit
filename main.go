package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"

	"github.com/jesseduffield/gocui"
	git "gopkg.in/src-d/go-git.v4"
)

// ErrSubProcess is raised when we are running a subprocess
var (
	ErrSubprocess = errors.New("running subprocess")
	subprocess    *exec.Cmd

	commit  string
	version = "unversioned"

	date          string
	debuggingFlag = flag.Bool("debug", false, "a boolean")
	versionFlag   = flag.Bool("v", false, "Print the current version")

	w *git.Worktree
	r *git.Repository
)

func homeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func projectPath(path string) string {
	gopath := os.Getenv("GOPATH")
	return filepath.FromSlash(gopath + "/src/github.com/jesseduffield/lazygit/" + path)
}

func devLog(objects ...interface{}) {
	localLog("development.log", objects...)
}

func objectLog(object interface{}) {
	if !*debuggingFlag {
		return
	}
	str := spew.Sdump(object)
	localLog("development.log", str)
}

func commandLog(objects ...interface{}) {
	localLog("commands.log", objects...)
}

func localLog(path string, objects ...interface{}) {
	if !*debuggingFlag {
		return
	}
	f, err := os.OpenFile(projectPath(path), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	log.SetOutput(f)
	for _, object := range objects {
		log.Println(fmt.Sprint(object))
	}
}

func navigateToRepoRootDirectory() {
	_, err := os.Stat(".git")
	for os.IsNotExist(err) {
		devLog("going up a directory to find the root")
		os.Chdir("..")
		_, err = os.Stat(".git")
	}
}

// when building the binary, `version` is set as a compile-time variable, along
// with `date` and `commit`. If this program has been opened directly via go,
// we will populate the `version` with VERSION in the lazygit root directory
func fallbackVersion() string {
	path := projectPath("VERSION")
	byteVersion, err := ioutil.ReadFile(path)
	if err != nil {
		return "unversioned"
	}
	return string(byteVersion)
}

func setupWorktree() {
	var err error
	r, err = git.PlainOpen(".")
	if err != nil {
		panic(err)
	}

	w, err = r.Worktree()
	if err != nil {
		panic(err)
	}
}

func main() {
	devLog("\n\n\n\n\n\n\n\n\n\n")
	flag.Parse()
	if version == "unversioned" {
		version = fallbackVersion()
	}
	if *versionFlag {
		fmt.Printf("commit=%s, build date=%s, version=%s\n", commit, date, version)
		os.Exit(0)
	}
	verifyInGitRepo()
	navigateToRepoRootDirectory()
	setupWorktree()
	for {
		if err := run(); err != nil {
			if err == gocui.ErrQuit {
				break
			} else if err == ErrSubprocess {
				subprocess.Run()
			} else {
				log.Panicln(err)
			}
		}
	}
}
