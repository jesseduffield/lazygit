package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/logs"
	"github.com/jesseduffield/lazygit/pkg/updates"
)

// App is the struct that's instantiated from within main.go and it manages
// bootstrapping and running the application.
type App struct {
	*common.Common
	closers   []io.Closer
	Config    config.AppConfigurer
	OSCommand *oscommands.OSCommand
	Gui       *gui.Gui
}

func Run(
	config config.AppConfigurer,
	common *common.Common,
	startArgs appTypes.StartArgs,
) {
	app, err := NewApp(config, startArgs.IntegrationTest, common)

	if err == nil {
		err = app.Run(startArgs)
	}

	if err != nil {
		if errorMessage, known := knownError(common.Tr, err); known {
			log.Fatal(errorMessage)
		}
		newErr := errors.Wrap(err, 0)
		stackTrace := newErr.ErrorStack()
		app.Log.Error(stackTrace)

		log.Fatalf("%s: %s\n\n%s", common.Tr.ErrorOccurred, constants.Links.Issues, stackTrace)
	}
}

func NewCommon(config config.AppConfigurer) (*common.Common, error) {
	userConfig := config.GetUserConfig()
	appState := config.GetAppState()
	log := newLogger(config)
	// Initialize with English for the time being; the real translation set for
	// the configured language will be read after reading the user config
	tr := i18n.EnglishTranslationSet()

	cmn := &common.Common{
		Log:      log,
		Tr:       tr,
		AppState: appState,
		Debug:    config.GetDebug(),
		Fs:       afero.NewOsFs(),
	}
	cmn.SetUserConfig(userConfig)
	return cmn, nil
}

func newLogger(cfg config.AppConfigurer) *logrus.Entry {
	if cfg.GetDebug() {
		logPath, err := config.LogPath()
		if err != nil {
			log.Fatal(err)
		}
		return logs.NewDevelopmentLogger(logPath)
	} else {
		return logs.NewProductionLogger()
	}
}

// NewApp bootstrap a new application
func NewApp(config config.AppConfigurer, test integrationTypes.IntegrationTest, common *common.Common) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
		Common:  common,
	}

	app.OSCommand = oscommands.NewOSCommand(common, config, oscommands.GetPlatform(), oscommands.NewNullGuiIO(app.Log))

	updater, err := updates.NewUpdater(common, config, app.OSCommand)
	if err != nil {
		return app, err
	}

	dirName, err := os.Getwd()
	if err != nil {
		return app, err
	}

	gitVersion, err := app.validateGitVersion()
	if err != nil {
		return app, err
	}

	// If we're not in a repo, GetRepoPaths will return an error. The error is moot for us
	// at this stage, since we'll try to init a new repo in setupRepo(), below
	repoPaths, err := git_commands.GetRepoPaths(app.OSCommand.Cmd, gitVersion)
	if err != nil {
		common.Log.Infof("Error getting repo paths: %v", err)
	}

	showRecentRepos, err := app.setupRepo(repoPaths)
	if err != nil {
		return app, err
	}

	// used for testing purposes
	if os.Getenv("SHOW_RECENT_REPOS") == "true" {
		showRecentRepos = true
	}

	app.Gui, err = gui.NewGui(common, config, gitVersion, updater, showRecentRepos, dirName, test)
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app *App) validateGitVersion() (*git_commands.GitVersion, error) {
	version, err := git_commands.GetGitVersion(app.OSCommand)
	// if we get an error anywhere here we'll show the same status
	minVersionError := errors.New(app.Tr.MinGitVersionError)
	if err != nil {
		return nil, minVersionError
	}

	if version.IsOlderThan(2, 20, 0) {
		return nil, minVersionError
	}

	return version, nil
}

func isDirectoryAGitRepository(dir string) (bool, error) {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return info != nil, err
}

func openRecentRepo(app *App) bool {
	for _, repoDir := range app.Config.GetAppState().RecentRepos {
		if isRepo, _ := isDirectoryAGitRepository(repoDir); isRepo {
			if err := os.Chdir(repoDir); err == nil {
				return true
			}
		}
	}

	return false
}

func (app *App) setupRepo(
	repoPaths *git_commands.RepoPaths,
) (bool, error) {
	if env.GetGitDirEnv() != "" {
		// we've been given the git dir directly. Skip setup
		return false, nil
	}

	// if we are not in a git repo, we ask if we want to `git init`
	if repoPaths == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return false, err
		}

		if isRepo, err := isDirectoryAGitRepository(cwd); isRepo {
			return false, err
		}

		var shouldInitRepo bool
		initialBranchArg := ""
		switch app.UserConfig().NotARepository {
		case "prompt":
			// Offer to initialize a new repository in current directory.
			fmt.Print(app.Tr.CreateRepo)
			response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			shouldInitRepo = (strings.Trim(response, " \r\n") == "y")
			if shouldInitRepo {
				// Ask for the initial branch name
				fmt.Print(app.Tr.InitialBranch)
				response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
				if trimmedResponse := strings.Trim(response, " \r\n"); len(trimmedResponse) > 0 {
					initialBranchArg += "--initial-branch=" + trimmedResponse
				}
			}
		case "create":
			shouldInitRepo = true
		case "skip":
			shouldInitRepo = false
		case "quit":
			fmt.Fprintln(os.Stderr, app.Tr.NotARepository)
			os.Exit(1)
		default:
			fmt.Fprintln(os.Stderr, app.Tr.IncorrectNotARepository)
			os.Exit(1)
		}

		if shouldInitRepo {
			args := []string{"git", "init"}
			if initialBranchArg != "" {
				args = append(args, initialBranchArg)
			}
			if err := app.OSCommand.Cmd.New(args).Run(); err != nil {
				return false, err
			}

			return false, nil
		}

		// check if we have a recent repo we can open
		for _, repoDir := range app.Config.GetAppState().RecentRepos {
			if isRepo, _ := isDirectoryAGitRepository(repoDir); isRepo {
				if err := os.Chdir(repoDir); err == nil {
					return true, nil
				}
			}
		}

		fmt.Fprintln(os.Stderr, app.Tr.NoRecentRepositories)
		os.Exit(1)
	}

	// Run this afterward so that the previous repo creation steps can run without this interfering
	if repoPaths.IsBareRepo() {

		fmt.Print(app.Tr.BareRepo)

		response, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		if shouldOpenRecent := strings.Trim(response, " \r\n") == "y"; !shouldOpenRecent {
			os.Exit(0)
		}

		if didOpenRepo := openRecentRepo(app); didOpenRepo {
			return true, nil
		}

		fmt.Println(app.Tr.NoRecentRepositories)
		os.Exit(1)
	}

	return false, nil
}

func (app *App) Run(startArgs appTypes.StartArgs) error {
	err := app.Gui.RunAndHandleError(startArgs)
	return err
}

// Close closes any resources
func (app *App) Close() error {
	for _, closer := range app.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}
