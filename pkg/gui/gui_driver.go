package gui

import (
	"fmt"
	"strings"
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
type GuiDriver struct {
	gui *Gui
}

var _ integrationTypes.GuiDriver = &GuiDriver{}

func (self *GuiDriver) PressKey(keyStr string) {
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

func (self *GuiDriver) Keys() config.KeybindingConfig {
	return self.gui.Config.GetUserConfig().Keybinding
}

func (self *GuiDriver) CurrentContext() types.Context {
	return self.gui.c.CurrentContext()
}

func (self *GuiDriver) Fail(message string) {
	self.gui.g.Close()
	// need to give the gui time to close
	time.Sleep(time.Millisecond * 100)
	panic(fmt.Sprintf("%s\nLog:\n%s", message, strings.Join(self.gui.CmdLog, "\n")))
}

// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
func (self *GuiDriver) Log(message string) {
	self.gui.c.Log.Warn(message)
}

// logs in the actual UI (in the commands panel)
func (self *GuiDriver) LogUI(message string) {
	self.gui.c.LogAction(message)
}

func (self *GuiDriver) CheckedOutRef() *models.Branch {
	return self.gui.helpers.Refs.GetCheckedOutRef()
}

func (self *GuiDriver) MainView() *gocui.View {
	return self.gui.mainView()
}

func (self *GuiDriver) SecondaryView() *gocui.View {
	return self.gui.secondaryView()
}

func (self *GuiDriver) View(viewName string) *gocui.View {
	view, err := self.gui.g.View(viewName)
	if err != nil {
		panic(err)
	}
	return view
}
