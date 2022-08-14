package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/integrii/flaggy"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/logs"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
	yaml "github.com/jesseduffield/yaml"
	"github.com/samber/lo"
)

const DEFAULT_VERSION = "unversioned"

var (
	commit      string
	version     = DEFAULT_VERSION
	date        string
	buildSource = "unknown"
)

func main() {
	updateBuildInfo()

	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	repoPath := ""
	flaggy.String(&repoPath, "p", "path", "Path of git repo. (equivalent to --work-tree=<path> --git-dir=<path>/.git/)")

	filterPath := ""
	flaggy.String(&filterPath, "f", "filter", "Path to filter on in `git log -- <path>`. When in filter mode, the commits, reflog, and stash are filtered based on the given path, and some operations are restricted")

	gitArg := ""
	flaggy.AddPositionalValue(&gitArg, "git-arg", 1, false, "Panel to focus upon opening lazygit. Accepted values (based on git terminology): status, branch, log, stash. Ignored if --filter arg is passed.")

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
	flaggy.String(&customConfig, "ucf", "use-config-file", "Comma separated list to custom config file(s)")

	flaggy.Parse()

	if os.Getenv("DEBUG") == "TRUE" {
		debuggingFlag = true
	}

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
		gitVersion := getGitVersionInfo()
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s, git version=%s\n", commit, date, buildSource, version, runtime.GOOS, runtime.GOARCH, gitVersion)
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
		logs.TailLogs()
		os.Exit(0)
	}

	if workTree != "" {
		if err := os.Chdir(workTree); err != nil {
			log.Fatal(err.Error())
		}
	}

	tempDir, err := os.MkdirTemp("", "lazygit-*")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer os.RemoveAll(tempDir)

	appConfig, err := config.NewAppConfig("lazygit", version, commit, date, buildSource, debuggingFlag, tempDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	common, err := app.NewCommon(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	if daemon.InDaemonMode() {
		daemon.Handle(common)
		return
	}

	parsedGitArg := parseGitArg(gitArg)

	app.Run(appConfig, common, types.NewStartArgs(filterPath, parsedGitArg))
}

func parseGitArg(gitArg string) types.GitArg {
	typedArg := types.GitArg(gitArg)

	// using switch so that linter catches when a new git arg value is defined but not handled here
	switch typedArg {
	case types.GitArgNone, types.GitArgStatus, types.GitArgBranch, types.GitArgLog, types.GitArgStash:
		return typedArg
	}

	permittedValues := []string{
		string(types.GitArgStatus),
		string(types.GitArgBranch),
		string(types.GitArgLog),
		string(types.GitArgStash),
	}

	log.Fatalf("Invalid git arg value: '%s'. Must be one of the following values: %s. e.g. 'lazygit status'. See 'lazygit --help'.",
		gitArg,
		strings.Join(permittedValues, ", "),
	)

	panic("unreachable")
}

func updateBuildInfo() {
	// if the version has already been set by build flags then we'll honour that.
	// chances are it's something like v0.31.0 which is more informative than a
	// commit hash.
	if version != DEFAULT_VERSION {
		return
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	revision, ok := lo.Find(buildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.revision"
	})
	if ok {
		commit = revision.Value
		// if lazygit was built from source we'll show the version as the
		// abbreviated commit hash
		version = utils.ShortSha(revision.Value)
	}

	// if version hasn't been set we assume that neither has the date
	time, ok := lo.Find(buildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.time"
	})
	if ok {
		date = time.Value
	}
}

func getGitVersionInfo() string {
	cmd := secureexec.Command("git", "--version")
	stdout, _ := cmd.Output()
	gitVersion := strings.Trim(strings.TrimPrefix(string(stdout), "git version "), " \r\n")
	return gitVersion
}
