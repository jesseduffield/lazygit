package app

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/heroku/rollrus"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	closers []io.Closer

	Config     config.AppConfigurer
	Log        *logrus.Entry
	OSCommand  *commands.OSCommand
	GitCommand *commands.GitCommand
	Gui        *gui.Gui
	Tr         *i18n.Localizer
	Updater    *updates.Updater // may only need this on the Gui
}

func newProductionLogger(config config.AppConfigurer) *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log
}

func newDevelopmentLogger() *logrus.Logger {
	log := logrus.New()
	file, err := os.OpenFile("development.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file") // TODO: don't panic (also, remove this call to the `panic` function)
	}
	log.SetOutput(file)
	return log
}

func newLogger(config config.AppConfigurer) *logrus.Entry {
	var log *logrus.Logger
	environment := "production"
	if config.GetDebug() {
		environment = "development"
		log = newDevelopmentLogger()
	} else {
		log = newProductionLogger(config)
	}
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

// Setup bootstrap a new application
func Setup(config config.AppConfigurer) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
	}
	var err error
	app.Log = newLogger(config)
	app.OSCommand = commands.NewOSCommand(app.Log, config)

	app.Tr = i18n.NewLocalizer(app.Log)

	app.Updater, err = updates.NewUpdater(app.Log, config, app.OSCommand, app.Tr)
	if err != nil {
		return app, err
	}
	app.GitCommand, err = commands.NewGitCommand(app.Log, app.OSCommand, app.Tr)
	if err != nil {
		return app, err
	}
	app.Gui, err = gui.NewGui(app.Log, app.GitCommand, app.OSCommand, app.Tr, config, app.Updater)
	if err != nil {
		return app, err
	}
	return app, nil
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
