package gui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

// this gives our integration test a way of interacting with the gui for sending keypresses
// and reading state.
type GuiAdapterImpl struct {
	gui *Gui
}

var _ integrationTypes.GuiAdapter = &GuiAdapterImpl{}

func (self *GuiAdapterImpl) PressKey(keyStr string) {
	key := keybindings.GetKey(keyStr)

	var r rune
	var tcellKey tcell.Key
	switch v := key.(type) {
	case rune:
		r = v
		tcellKey = tcell.KeyRune
	case gocui.Key:
		tcellKey = tcell.Key(v)
	}

	self.gui.g.ReplayedEvents.Keys <- gocui.NewTcellKeyEventWrapper(
		tcell.NewEventKey(tcellKey, r, tcell.ModNone),
		0,
	)
}

func (self *GuiAdapterImpl) Keys() config.KeybindingConfig {
	return self.gui.Config.GetUserConfig().Keybinding
}

func (self *GuiAdapterImpl) CurrentContext() types.Context {
	return self.gui.c.CurrentContext()
}

func (self *GuiAdapterImpl) Model() *types.Model {
	return self.gui.State.Model
}

func (self *GuiAdapterImpl) Fail(message string) {
	self.gui.g.Close()
	// need to give the gui time to close
	time.Sleep(time.Millisecond * 100)
	panic(message)
}

// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
func (self *GuiAdapterImpl) Log(message string) {
	self.gui.c.Log.Warn(message)
}

// logs in the actual UI (in the commands panel)
func (self *GuiAdapterImpl) LogUI(message string) {
	self.gui.c.LogAction(message)
}

func (self *GuiAdapterImpl) CheckedOutRef() *models.Branch {
	return self.gui.helpers.Refs.GetCheckedOutRef()
}
