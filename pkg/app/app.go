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

	"github.com/jesseduffield/generics/slices"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
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
	Updater   *updates.Updater // may only need this on the Gui
}

func Run(
	config config.AppConfigurer,
	common *common.Common,
	startArgs appTypes.StartArgs,
) {
	app, err := NewApp(config, common)

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

	var err error
	log := newLogger(config)
	tr, err := i18n.NewTranslationSetFromConfig(log, userConfig.Gui.Language)
	if err != nil {
		return nil, err
	}

	return &common.Common{
		Log:        log,
		Tr:         tr,
		UserConfig: userConfig,
		Debug:      config.GetDebug(),
	}, nil
}

// NewApp bootstrap a new application
func NewApp(config config.AppConfigurer, common *common.Common) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
		Common:  common,
	}

	app.OSCommand = oscommands.NewOSCommand(common, config, oscommands.GetPlatform(), oscommands.NewNullGuiIO(app.Log))

	var err error
	app.Updater, err = updates.NewUpdater(common, config, app.OSCommand)
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

	showRecentRepos, err := app.setupRepo()
	if err != nil {
		return app, err
	}

	app.Gui, err = gui.NewGui(common, config, gitVersion, app.Updater, showRecentRepos, dirName)
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

	if version.IsOlderThan(2, 0, 0) {
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

func (app *App) setupRepo() (bool, error) {
	if env.GetGitDirEnv() != "" {
		// we've been given the git dir directly. We'll verify this dir when initializing our Git object
		return false, nil
	}

	// if we are not in a git repo, we ask if we want to `git init`
	if err := commands.VerifyInGitRepo(app.OSCommand); err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			return false, err
		}

		if isRepo, err := isDirectoryAGitRepository(cwd); isRepo {
			return false, err
		}

		var shouldInitRepo bool
		initialBranchArg := ""
		switch app.UserConfig.NotARepository {
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
					initialBranchArg += "--initial-branch=" + app.OSCommand.Quote(trimmedResponse)
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
			if err := app.OSCommand.Cmd.New("git init " + initialBranchArg).Run(); err != nil {
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
	if isBare, err := git_commands.IsBareRepo(app.OSCommand); isBare {
		if err != nil {
			return false, err
		}

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
	return slices.TryForEach(app.closers, func(closer io.Closer) error {
		return closer.Close()
	})
}
