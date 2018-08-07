package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/fatih/color"
)

var (
	startTime        time.Time
	debugging        bool
	Rev              string
	builddate        string
	debuggingPointer = flag.Bool("debug", false, "a boolean")
	versionFlag      = flag.Bool("v", false, "Print the current version")
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
	startTime = time.Now()
	debugging = *debuggingPointer
	devLog("\n\n\n\n\n\n\n\n\n\n")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("rev=%s, build date=%s", Rev, builddate)
		os.Exit(0)
	}
	verifyInGitRepo()
	navigateToRepoRootDirectory()
	run()
}
