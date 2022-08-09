package types

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
)

type Test interface {
	Run(GuiAdapter)
	SetupConfig(config *config.AppConfig)
}

// this is the interface through which our integration tests interact with the lazygit gui
type GuiAdapter interface {
	PressKey(string)
	Keys() config.KeybindingConfig
	CurrentContext() Context
	Model() *Model
	Fail(message string)
	// These two log methods are for the sake of debugging while testing. There's no need to actually
	// commit any logging.
	// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
	Log(message string)
	// logs in the actual UI (in the commands panel)
	LogUI(message string)
	CheckedOutRef() *models.Branch
}
