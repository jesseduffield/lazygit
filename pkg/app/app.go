package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/heroku/rollrus"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	closers []io.Closer

	Config        config.AppConfigurer
	Log           *logrus.Entry
	OSCommand     *commands.OSCommand
	GitCommand    *commands.GitCommand
	Gui           *gui.Gui
	Tr            *i18n.Localizer
	Updater       *updates.Updater // may only need this on the Gui
	ClientContext string
}

func newProductionLogger(config config.AppConfigurer) *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log
}

func globalConfigDir() string {
	configDirs := configdir.New("jesseduffield", "lazygit")
	configDir := configDirs.QueryFolders(configdir.Global)[0]
	return configDir.Path
}

func newDevelopmentLogger(config config.AppConfigurer) *logrus.Logger {
	log := logrus.New()
	file, err := os.OpenFile(filepath.Join(globalConfigDir(), "development.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file") // TODO: don't panic (also, remove this call to the `panic` function)
	}
	log.SetOutput(file)
	return log
}

func newLogger(config config.AppConfigurer) *logrus.Entry {
	var log *logrus.Logger
	environment := "production"
	if config.GetDebug() || os.Getenv("DEBUG") == "TRUE" {
		environment = "development"
		log = newDevelopmentLogger(config)
	} else {
		log = newProductionLogger(config)
	}

	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{}

	if config.GetUserConfig().GetString("reporting") == "on" {
		// this isn't really a secret token: it only has permission to push new rollbar items
		hook := rollrus.NewHook("23432119147a4367abf7c0de2aa99a2d", environment)
		log.Hooks.Add(hook)
	}
	return log.WithFields(logrus.Fields{
		"debug":     config.GetDebug(),
		"version":   config.GetVersion(),
		"commit":    config.GetCommit(),
		"buildDate": config.GetBuildDate(),
	})
}

// NewApp bootstrap a new application
func NewApp(config config.AppConfigurer) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
	}
	var err error
	app.Log = newLogger(config)
	app.Tr = i18n.NewLocalizer(app.Log)

	// if we are being called in 'demon' mode, we can just return here
	app.ClientContext = os.Getenv("LAZYGIT_CLIENT_COMMAND")
	if app.ClientContext != "" {
		return app, nil
	}

	app.OSCommand = commands.NewOSCommand(app.Log, config)

	app.Updater, err = updates.NewUpdater(app.Log, config, app.OSCommand, app.Tr)
	if err != nil {
		return app, err
	}
	app.GitCommand, err = commands.NewGitCommand(app.Log, app.OSCommand, app.Tr, app.Config)
	if err != nil {
		if strings.Contains(err.Error(), "Not a git repository") {
			fmt.Println("Not in a git repository. Use `git init` to create a new one")
			os.Exit(1)
		}
		return app, err
	}
	app.Gui, err = gui.NewGui(app.Log, app.GitCommand, app.OSCommand, app.Tr, config, app.Updater)
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app *App) Run() error {
	if app.ClientContext == "INTERACTIVE_REBASE" {
		return app.Rebase()
	}

	if app.ClientContext == "EXIT_IMMEDIATELY" {
		os.Exit(0)
	}

	return app.Gui.RunWithSubprocesses()
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

	} else if strings.HasSuffix(os.Args[1], ".git/COMMIT_EDITMSG") {
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
