package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
)

var (
	commit      string
	version     = "unversioned"
	date        string
	buildSource = "unknown"

	configFlag    = flag.Bool("config", false, "Print the current default config")
	debuggingFlag = flag.Bool("debug", false, "a boolean")
	versionFlag   = flag.Bool("v", false, "Print the current version")
)

func projectPath(path string) string {
	gopath := os.Getenv("GOPATH")
	return filepath.FromSlash(gopath + "/src/github.com/jesseduffield/lazygit/" + path)
}

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s\n", commit, date, buildSource, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if *configFlag {
		fmt.Printf("%s\n", config.GetDefaultConfig())
		os.Exit(0)
	}
	appConfig, err := config.NewAppConfig("lazygit", version, commit, date, buildSource, debuggingFlag)
	if err != nil {
		panic(err)
	}

	app, err := app.NewApp(appConfig)
	if err != nil {
		app.Log.Error(err.Error())
		panic(err)
	}

	app.Gui.RunWithSubprocesses()
}
