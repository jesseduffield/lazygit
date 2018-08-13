package app

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui"
)

// App struct
type App struct {
	closers []io.Closer

	Config     config.AppConfigurer
	Log        *logrus.Logger
	OSCommand  *commands.OSCommand
	GitCommand *commands.GitCommand
	Gui        *gocui.Gui
}

// NewApp retruns a new applications
func NewApp(config config.AppConfigurer) (*App, error) {
	app := &App{
		closers: []io.Closer{},
		Config:  config,
	}
	var err error
	app.Log = logrus.New()
	app.OSCommand, err = commands.NewOSCommand(app.Log)
	if err != nil {
		return nil, err
	}
	app.GitCommand, err = commands.NewGitCommand(app.Log, app.OSCommand)
	if err != nil {
		return nil, err
	}
	app.Gui, err = gui.NewGui(app.Log, app.GitCommand, config.GetVersion())
	if err != nil {
		return nil, err
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
