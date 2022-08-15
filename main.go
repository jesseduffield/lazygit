package main

import (
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

// These values may be set by the build script via the LDFLAGS argument
var (
	commit      string
	date        string
	version     string
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

	ldFlagsBuildInfo := &app.BuildInfo{
		Commit:      commit,
		Date:        date,
		Version:     version,
		BuildSource: buildSource,

	}

	app.Start(ldFlagsBuildInfo, nil)
}

func getGitVersionInfo() string {
	cmd := secureexec.Command("git", "--version")
	stdout, _ := cmd.Output()
	gitVersion := strings.Trim(strings.TrimPrefix(string(stdout), "git version "), " \r\n")
	return gitVersion
}
