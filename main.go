package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
)

var (
	commit  string
	version = "unversioned"
	date    string

	debuggingFlag = flag.Bool("debug", false, "a boolean")
	versionFlag   = flag.Bool("v", false, "Print the current version")
)

func projectPath(path string) string {
	gopath := os.Getenv("GOPATH")
	return filepath.FromSlash(gopath + "/src/github.com/jesseduffield/lazygit/" + path)
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

func main() {
	flag.Parse()
	if version == "unversioned" {
		version = fallbackVersion()
	}
	if *versionFlag {
		fmt.Printf("commit=%s, build date=%s, version=%s, os=%s, arch=%s\n", commit, date, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	appConfig, err := config.NewAppConfig("lazygit", version, commit, date, debuggingFlag)
	if err != nil {
		panic(err)
	}

	app, err := app.NewApp(appConfig)
	app.Log.Info(err)
	app.GitCommand.SetupGit()
	app.Gui.RunWithSubprocesses()
}
