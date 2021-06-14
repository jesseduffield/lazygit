package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aybabtme/humanlog"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	closers []io.Closer

	Config        config.AppConfigurer
	Log           *logrus.Entry
	OSCommand     *oscommands.OSCommand
	GitCommand    *commands.GitCommand
	Gui           *gui.Gui
	Tr            *i18n.TranslationSet
	Updater       *updates.Updater // may only need this on the Gui
	ClientContext string
}

type errorMapping struct {
	originalError string
	newError      string
}

func newProductionLogger(config config.AppConfigurer) *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	log.SetLevel(logrus.ErrorLevel)
	return log
}

func getLogLevel() logrus.Level {
	strLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(strLevel)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func newDevelopmentLogger(configurer config.AppConfigurer) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(getLogLevel())
	logPath, err := config.LogPath()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file") // TODO: don't panic (also, remove this call to the `panic` function)
	}
	logger.SetOutput(file)
	return logger
}

func newLogger(config config.AppConfigurer) *logrus.Entry {
	var log *logrus.Logger
	if config.GetDebug() || os.Getenv("DEBUG") == "TRUE" {
		log = newDevelopmentLogger(config)
	} else {
		log = newProductionLogger(config)
	}

	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{}

	return log.WithFields(logrus.Fields{
		"debug":     config.GetDebug(),
		"version":   config.GetVersion(),
		"commit":    config.GetCommit(),
		"buildDate": config.GetBuildDate(),
	})
}

// NewApp bootstrap a new application
func NewApp(config config.AppConfigurer, filterPath string) (*App, error) {

	app := &App{
		closers: []io.Closer{},
		Config:  config,
	}
	var err error
	app.Log = newLogger(config)
	app.Tr = i18n.NewTranslationSet(app.Log)

	// if we are being called in 'demon' mode, we can just return here
	app.ClientContext = os.Getenv("LAZYGIT_CLIENT_COMMAND")
	if app.ClientContext != "" {
		return app, nil
	}

	app.OSCommand = oscommands.NewOSCommand(app.Log, config)

	app.Updater, err = updates.NewUpdater(app.Log, config, app.OSCommand, app.Tr)
	if err != nil {
		return app, err
	}

	showRecentRepos, err := app.setupRepo()
	if err != nil {
		return app, err
	}

	app.GitCommand, err = commands.NewGitCommand(app.Log, app.OSCommand, app.Tr, app.Config)
	if err != nil {
		return app, err
	}

	app.Gui, err = gui.NewGui(app.Log, app.GitCommand, app.OSCommand, app.Tr, config, app.Updater, filterPath, showRecentRepos)
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app *App) validateGitVersion() error {
	output, err := app.OSCommand.RunWithOutput(commands.GitVersionCmd())
	// if we get an error anywhere here we'll show the same status
	minVersionError := errors.New(app.Tr.MinGitVersionError)
	if err != nil {
		return minVersionError
	}

	if isGitVersionValid(output) {
		return nil
	}

	return minVersionError
}

func isGitVersionValid(versionStr string) bool {
	// output should be something like: 'git version 2.23.0 (blah)'
	re := regexp.MustCompile(`[^\d]+([\d\.]+)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) == 0 {
		return false
	}

	gitVersion := matches[1]
	majorVersion, err := strconv.Atoi(gitVersion[0:1])
	if err != nil {
		return false
	}
	if majorVersion < 2 {
		return false
	}

	return true
}

func (app *App) setupRepo() (bool, error) {
	if err := app.validateGitVersion(); err != nil {
		return false, err
	}

	if env.GetGitDirEnv() != "" {
		// we've been given the git dir directly. We'll verify this dir when initializing our GitCommand object
		return false, nil
	}

	// if we are not in a git repo, we ask if we want to `git init`
	if err := commands.VerifyInGitRepo(app.OSCommand); err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			return false, err
		}
		info, _ := os.Stat(filepath.Join(cwd, ".git"))
		if info != nil && info.IsDir() {
			return false, err // Current directory appears to be a git repository.
		}

		shouldInitRepo := true
		notARepository := app.Config.GetUserConfig().NotARepository
		if notARepository == "prompt" {
			// Offer to initialize a new repository in current directory.
			fmt.Print(app.Tr.CreateRepo)
			response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if strings.Trim(response, " \n") != "y" {
				shouldInitRepo = false
			}
		} else if notARepository == "skip" {
			shouldInitRepo = false
		}

		if !shouldInitRepo {
			// check if we have a recent repo we can open
			recentRepos := app.Config.GetAppState().RecentRepos
			if len(recentRepos) > 0 {
				var err error
				// try opening each repo in turn, in case any have been deleted
				for _, repoDir := range recentRepos {
					if err = os.Chdir(repoDir); err == nil {
						return true, nil
					}
				}
				return false, err
			}

			os.Exit(1)
		}
		if err := app.OSCommand.Run(commands.GitInitCmd()); err != nil {
			return false, err
		}
	}

	return false, nil
}

func (app *App) Run() error {
	if app.ClientContext == "INTERACTIVE_REBASE" {
		return app.Rebase()
	}

	if app.ClientContext == "EXIT_IMMEDIATELY" {
		os.Exit(0)
	}

	err := app.Gui.RunAndHandleError()
	return err
}

func gitDir() string {
	dir := env.GetGitDirEnv()
	if dir == "" {
		return ".git"
	}
	return dir
}

// Rebase contains logic for when we've been run in demon mode, meaning we've
// given lazygit as a command for git to call e.g. to edit a file
func (app *App) Rebase() error {
	app.Log.Info("Lazygit invoked as interactive rebase demon")
	app.Log.Info("args: ", os.Args)

	if strings.HasSuffix(os.Args[1], "git-rebase-todo") {
		if err := ioutil.WriteFile(os.Args[1], []byte(os.Getenv("LAZYGIT_REBASE_TODO")), 0644); err != nil {
			return err
		}

	} else if strings.HasSuffix(os.Args[1], filepath.Join(gitDir(), "COMMIT_EDITMSG")) { // TODO: test
		// if we are rebasing and squashing, we'll see a COMMIT_EDITMSG
		// but in this case we don't need to edit it, so we'll just return
	} else {
		app.Log.Info("Lazygit demon did not match on any use cases")
	}

	return nil
}

// Close closes any resources
func (app *App) Close() error {
	for _, closer := range app.closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// KnownError takes an error and tells us whether it's an error that we know about where we can print a nicely formatted version of it rather than panicking with a stack trace
func (app *App) KnownError(err error) (string, bool) {
	errorMessage := err.Error()

	knownErrorMessages := []string{app.Tr.MinGitVersionError}

	for _, message := range knownErrorMessages {
		if errorMessage == message {
			return message, true
		}
	}

	mappings := []errorMapping{
		{
			originalError: "fatal: not a git repository",
			newError:      app.Tr.NotARepository,
		},
	}

	for _, mapping := range mappings {
		if strings.Contains(errorMessage, mapping.originalError) {
			return mapping.newError, true
		}
	}
	return "", false
}

func TailLogs() {
	logFilePath, err := config.LogPath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tailing log file %s\n\n", logFilePath)

	_, err = os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Log file does not exist. Run `lazygit --debug` first to create the log file")
		}
		log.Fatal(err)
	}

	cmd := secureexec.Command("tail", "-f", logFilePath)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	opts := humanlog.DefaultOptions
	opts.Truncates = false
	if err := humanlog.Scanner(stdout, os.Stdout, opts); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
