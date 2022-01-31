package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-errors/errors"
	"github.com/integrii/flaggy"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/env"
	yaml "github.com/jesseduffield/yaml"
)

var (
	commit      string
	version     = "unversioned"
	date        string
	buildSource = "unknown"
)

func main() {
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	repoPath := ""
	flaggy.String(&repoPath, "p", "path", "Path of git repo. (equivalent to --work-tree=<path> --git-dir=<path>/.git/)")

	filterPath := ""
	flaggy.String(&filterPath, "f", "filter", "Path to filter on in `git log -- <path>`. When in filter mode, the commits, reflog, and stash are filtered based on the given path, and some operations are restricted")

	dump := ""
	flaggy.AddPositionalValue(&dump, "gitargs", 1, false, "Todo file")
	flaggy.DefaultParser.PositionalFlags[0].Hidden = true

	versionFlag := false
	flaggy.Bool(&versionFlag, "v", "version", "Print the current version")

	debuggingFlag := false
	flaggy.Bool(&debuggingFlag, "d", "debug", "Run in debug mode with logging (see --logs flag below). Use the LOG_LEVEL env var to set the log level (debug/info/warn/error)")

	logFlag := false
	flaggy.Bool(&logFlag, "l", "logs", "Tail lazygit logs (intended to be used when `lazygit --debug` is called in a separate terminal tab)")

	configFlag := false
	flaggy.Bool(&configFlag, "c", "config", "Print the default config")

	configDirFlag := false
	flaggy.Bool(&configDirFlag, "cd", "print-config-dir", "Print the config directory")

	useConfigDir := ""
	flaggy.String(&useConfigDir, "ucd", "use-config-dir", "override default config directory with provided directory")

	workTree := ""
	flaggy.String(&workTree, "w", "work-tree", "equivalent of the --work-tree git argument")

	gitDir := ""
	flaggy.String(&gitDir, "g", "git-dir", "equivalent of the --git-dir git argument")

	customConfig := ""
	flaggy.String(&customConfig, "ucf", "use-config-file", "Comma seperated list to custom config file(s)")

	flaggy.Parse()

	if repoPath != "" {
		if workTree != "" || gitDir != "" {
			log.Fatal("--path option is incompatible with the --work-tree and --git-dir options")
		}

		absRepoPath, err := filepath.Abs(repoPath)
		if err != nil {
			log.Fatal(err)
		}
		workTree = absRepoPath
		gitDir = filepath.Join(absRepoPath, ".git")
	}

	if customConfig != "" {
		os.Setenv("LG_CONFIG_FILE", customConfig)
	}

	if useConfigDir != "" {
		os.Setenv("CONFIG_DIR", useConfigDir)
	}

	if workTree != "" {
		env.SetGitWorkTreeEnv(workTree)
	}

	if gitDir != "" {
		env.SetGitDirEnv(gitDir)
	}

	if versionFlag {
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s\n", commit, date, buildSource, version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if configFlag {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		err := encoder.Encode(config.GetDefaultConfig())
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s\n", buf.String())
		os.Exit(0)
	}

	if configDirFlag {
		fmt.Printf("%s\n", config.ConfigDir())
		os.Exit(0)
	}

	if logFlag {
		app.TailLogs()
		os.Exit(0)
	}

	if workTree != "" {
		if err := os.Chdir(workTree); err != nil {
			log.Fatal(err.Error())
		}
	}

	appConfig, err := config.NewAppConfig("lazygit", version, commit, date, buildSource, debuggingFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := app.NewApp(appConfig)

	if err == nil {
		err = app.Run(filterPath)
	}

	if err != nil {
		if errorMessage, known := app.KnownError(err); known {
			log.Fatal(errorMessage)
		}
		newErr := errors.Wrap(err, 0)
		stackTrace := newErr.ErrorStack()
		app.Log.Error(stackTrace)

		log.Fatal(fmt.Sprintf("%s: %s\n\n%s", app.Tr.ErrorOccurred, constants.Links.Issues, stackTrace))
	}
}
