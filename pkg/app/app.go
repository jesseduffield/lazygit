package app

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

// App struct
type App struct {
	closers []io.Closer

	Config     config.AppConfigurer
	Log        *logrus.Logger
	OSCommand  *commands.OSCommand
	GitCommand *commands.GitCommand
	Gui        *gui.Gui
	Tr         *i18n.Localizer
}

func newLogger(config config.AppConfigurer) *logrus.Logger {
	log := logrus.New()
	if !config.GetDebug() {
		log.Out = ioutil.Discard
		return log
	}
	file, err := os.OpenFile("development.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file") // TODO: don't panic (also, remove this call to the `panic` function)
	}
	log.SetOutput(file)
	return log
}

// NewApp retruns a new applications
func NewApp(config config.AppConfigurer) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
	}
	var err error
	app.Log = newLogger(config)
	app.OSCommand, err = commands.NewOSCommand(app.Log)
	if err != nil {
		return app, err
	}

	app.Tr = i18n.NewLocalizer(app.Log)

	app.GitCommand, err = commands.NewGitCommand(app.Log, app.OSCommand)
	if err != nil {
		return app, err
	}
	app.Gui, err = gui.NewGui(app.Log, app.GitCommand, app.OSCommand, app.Tr, config)
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
