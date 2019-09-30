package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-errors/errors"
	"github.com/integrii/flaggy"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
)

var (
	commit      string
	version     = "unversioned"
	date        string
	buildSource = "unknown"
)

func projectPath(path string) string {
	gopath := os.Getenv("GOPATH")
	return filepath.FromSlash(gopath + "/src/github.com/jesseduffield/lazygit/" + path)
}

func main() {
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	repoPath := "."
	flaggy.String(&repoPath, "p", "path", "Path of git repo")

	dump := ""
	flaggy.AddPositionalValue(&dump, "gitargs", 1, false, "Todo file")
	flaggy.DefaultParser.PositionalFlags[0].Hidden = true

	versionFlag := false
	flaggy.Bool(&versionFlag, "v", "version", "Print the current version")

	debuggingFlag := false
	flaggy.Bool(&debuggingFlag, "d", "debug", "Run in debug mode with logging")

	configFlag := false
	flaggy.Bool(&configFlag, "c", "config", "Print the current default config")

	flaggy.Parse()

	if versionFlag {
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s\n", commit, date, buildSource, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if configFlag {
		fmt.Printf("%s\n", config.GetDefaultConfig())
		os.Exit(0)
	}

	if repoPath != "." {
		if err := os.Chdir(repoPath); err != nil {
			log.Fatal(err.Error())
		}
	}

	appConfig, err := config.NewAppConfig("lazygit", version, commit, date, buildSource, debuggingFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := app.NewApp(appConfig)

	if err == nil {
		err = app.Run()
	}

	if err != nil {
		if errorMessage, known := app.KnownError(err); known {
			log.Fatal(errorMessage)
		}
		newErr := errors.Wrap(err, 0)
		stackTrace := newErr.ErrorStack()
		app.Log.Error(stackTrace)

		log.Fatal(fmt.Sprintf("%s\n\n%s", app.Tr.SLocalize("ErrorOccurred"), stackTrace))
	}
}
