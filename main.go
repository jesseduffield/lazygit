package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

// ErrSubProcess is raised when we are running a subprocess
var (
	startTime     time.Time
	debugging     bool
	ErrSubprocess = errors.New("running subprocess")
	subprocess    *exec.Cmd
)

func homeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func devLog(objects ...interface{}) {
	localLog(color.FgWhite, homeDirectory()+"/go/src/github.com/jesseduffield/lazygit/development.log", objects...)
}

func colorLog(colour color.Attribute, objects ...interface{}) {
	localLog(colour, homeDirectory()+"/go/src/github.com/jesseduffield/lazygit/development.log", objects...)
}

func commandLog(objects ...interface{}) {
	localLog(color.FgWhite, homeDirectory()+"/go/src/github.com/jesseduffield/lazygit/commands.log", objects...)
}

func localLog(colour color.Attribute, path string, objects ...interface{}) {
	if !debugging {
		return
	}
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	for _, object := range objects {
		colorFunction := color.New(colour).SprintFunc()
		f.WriteString(colorFunction(fmt.Sprint(object)) + "\n")
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

func main() {
	debuggingPointer := flag.Bool("debug", false, "a boolean")
	flag.Parse()
	debugging = *debuggingPointer
	devLog("\n\n\n\n\n\n\n\n\n\n")
	startTime = time.Now()
	verifyInGitRepo()
	navigateToRepoRootDirectory()
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
