package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/integrii/flaggy"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/logs/tail"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type cliArgs struct {
	RepoPath           string
	FilterPath         string
	GitArg             string
	PrintVersionInfo   bool
	Debug              bool
	TailLogs           bool
	Profile            bool
	PrintDefaultConfig bool
	PrintConfigDir     bool
	UseConfigDir       string
	WorkTree           string
	GitDir             string
	CustomConfigFile   string
}

type BuildInfo struct {
	Commit      string
	Date        string
	Version     string
	BuildSource string
}

func Start(buildInfo *BuildInfo, integrationTest integrationTypes.IntegrationTest) {
	cliArgs := parseCliArgsAndEnvVars()
	mergeBuildInfo(buildInfo)

	if cliArgs.RepoPath != "" {
		if cliArgs.WorkTree != "" || cliArgs.GitDir != "" {
			log.Fatal("--path option is incompatible with the --work-tree and --git-dir options")
		}

		absRepoPath, err := filepath.Abs(cliArgs.RepoPath)
		if err != nil {
			log.Fatal(err)
		}

		if isRepo, err := isDirectoryAGitRepository(absRepoPath); err != nil || !isRepo {
			log.Fatal(absRepoPath + " is not a valid git repository.")
		}

		cliArgs.GitDir = filepath.Join(absRepoPath, ".git")
		err = os.Chdir(absRepoPath)
		if err != nil {
			log.Fatalf("Failed to change directory to %s: %v", absRepoPath, err)
		}
	} else if cliArgs.WorkTree != "" {
		env.SetWorkTreeEnv(cliArgs.WorkTree)

		if err := os.Chdir(cliArgs.WorkTree); err != nil {
			log.Fatalf("Failed to change directory to %s: %v", cliArgs.WorkTree, err)
		}
	}

	if cliArgs.CustomConfigFile != "" {
		os.Setenv("LG_CONFIG_FILE", cliArgs.CustomConfigFile)
	}

	if cliArgs.UseConfigDir != "" {
		os.Setenv("CONFIG_DIR", cliArgs.UseConfigDir)
	}

	if cliArgs.GitDir != "" {
		env.SetGitDirEnv(cliArgs.GitDir)
	}

	if cliArgs.PrintVersionInfo {
		gitVersion := getGitVersionInfo()
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s, git version=%s\n", buildInfo.Commit, buildInfo.Date, buildInfo.BuildSource, buildInfo.Version, runtime.GOOS, runtime.GOARCH, gitVersion)
		os.Exit(0)
	}

	if cliArgs.PrintDefaultConfig {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		err := encoder.Encode(config.GetDefaultConfig())
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s\n", buf.String())
		os.Exit(0)
	}

	if cliArgs.PrintConfigDir {
		fmt.Printf("%s\n", config.ConfigDir())
		os.Exit(0)
	}

	if cliArgs.TailLogs {
		logPath, err := config.LogPath()
		if err != nil {
			log.Fatal(err.Error())
		}

		tail.TailLogs(logPath)
		os.Exit(0)
	}

	tempDir, err := os.MkdirTemp("", "lazygit-*")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer os.RemoveAll(tempDir)

	appConfig, err := config.NewAppConfig("lazygit", buildInfo.Version, buildInfo.Commit, buildInfo.Date, buildInfo.BuildSource, cliArgs.Debug, tempDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	if integrationTest != nil {
		integrationTest.SetupConfig(appConfig)
	}

	common, err := NewCommon(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	if daemon.InDaemonMode() {
		daemon.Handle(common)
		return
	}

	if cliArgs.Profile {
		go func() {
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Fatal(err)
			}
		}()
	}

	parsedGitArg := parseGitArg(cliArgs.GitArg)

	Run(appConfig, common, appTypes.NewStartArgs(cliArgs.FilterPath, parsedGitArg, integrationTest))
}

func parseCliArgsAndEnvVars() *cliArgs {
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

	profile := false
	flaggy.Bool(&profile, "", "profile", "Start the profiler and serve it on http port 6060. See CONTRIBUTING.md for more info.")

	printDefaultConfig := false
	flaggy.Bool(&printDefaultConfig, "c", "config", "Print the default config")

	printConfigDir := false
	flaggy.Bool(&printConfigDir, "cd", "print-config-dir", "Print the config directory")

	useConfigDir := ""
	flaggy.String(&useConfigDir, "ucd", "use-config-dir", "override default config directory with provided directory")

	workTree := os.Getenv("GIT_WORK_TREE")
	flaggy.String(&workTree, "w", "work-tree", "equivalent of the --work-tree git argument")

	gitDir := os.Getenv("GIT_DIR")
	flaggy.String(&gitDir, "g", "git-dir", "equivalent of the --git-dir git argument")

	customConfigFile := ""
	flaggy.String(&customConfigFile, "ucf", "use-config-file", "Comma separated list to custom config file(s)")

	flaggy.Parse()

	if os.Getenv("DEBUG") == "TRUE" {
		debug = true
	}

	return &cliArgs{
		RepoPath:           repoPath,
		FilterPath:         filterPath,
		GitArg:             gitArg,
		PrintVersionInfo:   printVersionInfo,
		Debug:              debug,
		TailLogs:           tailLogs,
		Profile:            profile,
		PrintDefaultConfig: printDefaultConfig,
		PrintConfigDir:     printConfigDir,
		UseConfigDir:       useConfigDir,
		WorkTree:           workTree,
		GitDir:             gitDir,
		CustomConfigFile:   customConfigFile,
	}
}

func parseGitArg(gitArg string) appTypes.GitArg {
	typedArg := appTypes.GitArg(gitArg)

	// using switch so that linter catches when a new git arg value is defined but not handled here
	switch typedArg {
	case appTypes.GitArgNone, appTypes.GitArgStatus, appTypes.GitArgBranch, appTypes.GitArgLog, appTypes.GitArgStash:
		return typedArg
	}

	permittedValues := []string{
		string(appTypes.GitArgStatus),
		string(appTypes.GitArgBranch),
		string(appTypes.GitArgLog),
		string(appTypes.GitArgStash),
	}

	log.Fatalf("Invalid git arg value: '%s'. Must be one of the following values: %s. e.g. 'lazygit status'. See 'lazygit --help'.",
		gitArg,
		strings.Join(permittedValues, ", "),
	)

	panic("unreachable")
}

// the buildInfo struct we get passed in is based on what's baked into the lazygit
// binary via the LDFLAGS argument. Some lazygit distributions will make use of these
// arguments and some will not. Go recently started baking in build info
// into the binary by default e.g. the git commit hash. So in this function
// we merge the two together, giving priority to the stuff set by LDFLAGS.
// Note: this mutates the argument passed in
func mergeBuildInfo(buildInfo *BuildInfo) {
	// if the version has already been set by build flags then we'll honour that.
	// chances are it's something like v0.31.0 which is more informative than a
	// commit hash.
	if buildInfo.Version != "" {
		return
	}

	buildInfo.Version = "unversioned"

	goBuildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	revision, ok := lo.Find(goBuildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.revision"
	})
	if ok {
		buildInfo.Commit = revision.Value
		// if lazygit was built from source we'll show the version as the
		// abbreviated commit hash
		buildInfo.Version = utils.ShortHash(revision.Value)
	}

	// if version hasn't been set we assume that neither has the date
	time, ok := lo.Find(goBuildInfo.Settings, func(setting debug.BuildSetting) bool {
		return setting.Key == "vcs.time"
	})
	if ok {
		buildInfo.Date = time.Value
	}
}

func getGitVersionInfo() string {
	cmd := exec.Command("git", "--version")
	stdout, _ := cmd.Output()
	gitVersion := strings.Trim(strings.TrimPrefix(string(stdout), "git version "), " \r\n")
	return gitVersion
}
