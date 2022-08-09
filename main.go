package main

import (
	"os"
	"runtime/debug"

	"github.com/integrii/flaggy"
	"github.com/jesseduffield/lazygit/pkg/app"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

const DEFAULT_VERSION = "unversioned"

// These values may be set by the build script.
// we'll overwrite them if they haven't been set by the build script and if Go itself has set corresponding values in the binary
var (
	commit      string
	version     = DEFAULT_VERSION
	date        string
	buildSource = "unknown"
)

func main() {
	cliArgs := parseCliArgsAndEnvVars()
	buildInfo := getBuildInfo()

	app.Start(cliArgs, buildInfo, nil)
}

func parseCliArgsAndEnvVars() *app.CliArgs {
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	repoPath := ""
	flaggy.String(&repoPath, "p", "path", "Path of git repo. (equivalent to --work-tree=<path> --git-dir=<path>/.git/)")

	filterPath := ""
	flaggy.String(&filterPath, "f", "filter", "Path to filter on in `git log -- <path>`. When in filter mode, the commits, reflog, and stash are filtered based on the given path, and some operations are restricted")

	gitArg := ""
	flaggy.AddPositionalValue(&gitArg, "git-arg", 1, false, "Panel to focus upon opening lazygit. Accepted values (based on git terminology): status, branch, log, stash. Ignored if --filter arg is passed.")

	printVersionInfo := false
	flaggy.Bool(&printVersionInfo, "v", "version", "Print the current version")

	debug := false
	flaggy.Bool(&debug, "d", "debug", "Run in debug mode with logging (see --logs flag below). Use the LOG_LEVEL env var to set the log level (debug/info/warn/error)")

	tailLogs := false
	flaggy.Bool(&tailLogs, "l", "logs", "Tail lazygit logs (intended to be used when `lazygit --debug` is called in a separate terminal tab)")

	printDefaultConfig := false
	flaggy.Bool(&printDefaultConfig, "c", "config", "Print the default config")

	printConfigDir := false
	flaggy.Bool(&printConfigDir, "cd", "print-config-dir", "Print the config directory")

	useConfigDir := ""
	flaggy.String(&useConfigDir, "ucd", "use-config-dir", "override default config directory with provided directory")

	workTree := ""
	flaggy.String(&workTree, "w", "work-tree", "equivalent of the --work-tree git argument")

	gitDir := ""
	flaggy.String(&gitDir, "g", "git-dir", "equivalent of the --git-dir git argument")

	customConfigFile := ""
	flaggy.String(&customConfigFile, "ucf", "use-config-file", "Comma separated list to custom config file(s)")

	flaggy.Parse()

	if os.Getenv("DEBUG") == "TRUE" {
		debug = true
	}

	return &app.CliArgs{
		RepoPath:           repoPath,
		FilterPath:         filterPath,
		GitArg:             gitArg,
		PrintVersionInfo:   printVersionInfo,
		Debug:              debug,
		TailLogs:           tailLogs,
		PrintDefaultConfig: printDefaultConfig,
		PrintConfigDir:     printConfigDir,
		UseConfigDir:       useConfigDir,
		WorkTree:           workTree,
		GitDir:             gitDir,
		CustomConfigFile:   customConfigFile,
	}
}

func getBuildInfo() *app.BuildInfo {
	buildInfo := &app.BuildInfo{
		Commit:      commit,
		Date:        date,
		Version:     version,
		BuildSource: buildSource,
	}

	// if the version has already been set by build flags then we'll honour that.
	// chances are it's something like v0.31.0 which is more informative than a
	// commit hash.
	if buildInfo.Version != DEFAULT_VERSION {
		return buildInfo
	}

	goBuildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return buildInfo
	}

	revision, ok := lo.Find(goBuildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.revision"
	})
	if ok {
		buildInfo.Commit = revision.Value
		// if lazygit was built from source we'll show the version as the
		// abbreviated commit hash
		buildInfo.Version = utils.ShortSha(revision.Value)
	}

	// if version hasn't been set we assume that neither has the date
	time, ok := lo.Find(goBuildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.time"
	})
	if ok {
		buildInfo.Date = time.Value
	}

	return buildInfo
}
