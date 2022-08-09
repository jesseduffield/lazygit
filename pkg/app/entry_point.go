package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/logs"
	"gopkg.in/yaml.v3"
)

type CliArgs struct {
	RepoPath           string
	FilterPath         string
	GitArg             string
	PrintVersionInfo   bool
	Debug              bool
	TailLogs           bool
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

func Start(cliArgs *CliArgs, buildInfo *BuildInfo, integrationTest integrationTypes.IntegrationTest) {
	if cliArgs.RepoPath != "" {
		if cliArgs.WorkTree != "" || cliArgs.GitDir != "" {
			log.Fatal("--path option is incompatible with the --work-tree and --git-dir options")
		}

		absRepoPath, err := filepath.Abs(cliArgs.RepoPath)
		if err != nil {
			log.Fatal(err)
		}
		cliArgs.WorkTree = absRepoPath
		cliArgs.GitDir = filepath.Join(absRepoPath, ".git")
	}

	if cliArgs.CustomConfigFile != "" {
		os.Setenv("LG_CONFIG_FILE", cliArgs.CustomConfigFile)
	}

	if cliArgs.UseConfigDir != "" {
		os.Setenv("CONFIG_DIR", cliArgs.UseConfigDir)
	}

	if cliArgs.WorkTree != "" {
		env.SetGitWorkTreeEnv(cliArgs.WorkTree)
	}

	if cliArgs.GitDir != "" {
		env.SetGitDirEnv(cliArgs.GitDir)
	}

	if cliArgs.PrintVersionInfo {
		fmt.Printf("commit=%s, build date=%s, build source=%s, version=%s, os=%s, arch=%s\n", buildInfo.Commit, buildInfo.Date, buildInfo.BuildSource, buildInfo.Version, runtime.GOOS, runtime.GOARCH)
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
		logs.TailLogs()
		os.Exit(0)
	}

	if cliArgs.WorkTree != "" {
		if err := os.Chdir(cliArgs.WorkTree); err != nil {
			log.Fatal(err.Error())
		}
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

	parsedGitArg := parseGitArg(cliArgs.GitArg)

	Run(appConfig, common, appTypes.NewStartArgs(cliArgs.FilterPath, parsedGitArg, integrationTest))
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
