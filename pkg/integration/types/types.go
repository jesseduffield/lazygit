package types

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// these interfaces are used by the gui package so that it knows what it needs
// to provide to a test in order for the test to run.

type IntegrationTest interface {
	Run(GuiDriver)
	SetupConfig(config *config.AppConfig)
}

// this is the interface through which our integration tests interact with the lazygit gui
type GuiDriver interface {
	PressKey(string)
	Keys() config.KeybindingConfig
	CurrentContext() types.Context
	Model() *types.Model
	Fail(message string)
	// These two log methods are for the sake of debugging while testing. There's no need to actually
	// commit any logging.
	// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
	Log(message string)
	// logs in the actual UI (in the commands panel)
	LogUI(message string)
	CheckedOutRef() *models.Branch
	// the view that appears to the right of the side panel
	MainView() *gocui.View
	// the other view that sometimes appears to the right of the side panel
	// e.g. when we're showing both staged and unstaged changes
	SecondaryView() *gocui.View
	View(viewName string) *gocui.View
}
