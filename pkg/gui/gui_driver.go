package gui

import (
	"fmt"
	"os"
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
	gui        *Gui
	isIdleChan chan struct{}
	toastChan  chan string
	headless   bool
}

var _ integrationTypes.GuiDriver = &GuiDriver{}

func (self *GuiDriver) PressKey(keyStr string) {
	self.CheckAllToastsAcknowledged()

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

	self.waitTillIdle()
}

func (self *GuiDriver) Click(x, y int) {
	self.CheckAllToastsAcknowledged()

	self.gui.g.ReplayedEvents.MouseEvents <- gocui.NewTcellMouseEventWrapper(
		tcell.NewEventMouse(x, y, tcell.ButtonPrimary, 0),
		0,
	)
	self.waitTillIdle()
	self.gui.g.ReplayedEvents.MouseEvents <- gocui.NewTcellMouseEventWrapper(
		tcell.NewEventMouse(x, y, tcell.ButtonNone, 0),
		0,
	)
	self.waitTillIdle()
}

// wait until lazygit is idle (i.e. all processing is done) before continuing
func (self *GuiDriver) waitTillIdle() {
	<-self.isIdleChan
}

func (self *GuiDriver) CheckAllToastsAcknowledged() {
	if t := self.NextToast(); t != nil {
		self.Fail("Toast not acknowledged: " + *t)
	}
}

func (self *GuiDriver) Keys() config.KeybindingConfig {
	return self.gui.Config.GetUserConfig().Keybinding
}

func (self *GuiDriver) CurrentContext() types.Context {
	return self.gui.c.Context().Current()
}

func (self *GuiDriver) ContextForView(viewName string) types.Context {
	context, ok := self.gui.helpers.View.ContextForView(viewName)
	if !ok {
		return nil
	}

	return context
}

func (self *GuiDriver) Fail(message string) {
	currentView := self.gui.g.CurrentView()

	// Check for unacknowledged toast: it may give us a hint as to why the test failed
	toastMessage := ""
	if t := self.NextToast(); t != nil {
		toastMessage = fmt.Sprintf("Unacknowledged toast message: %s\n", *t)
	}

	fullMessage := fmt.Sprintf(
		"%s\nFinal Lazygit state:\n%s\nUpon failure, focused view was '%s'.\n%sLog:\n%s", message,
		self.gui.g.Snapshot(),
		currentView.Name(),
		toastMessage,
		strings.Join(self.gui.GuiLog, "\n"),
	)

	self.gui.g.Close()
	// need to give the gui time to close
	time.Sleep(time.Millisecond * 100)
	_, err := fmt.Fprintln(os.Stderr, fullMessage)
	if err != nil {
		panic("Test failed. Failed writing to stderr")
	}
	panic("Test failed")
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

func (self *GuiDriver) SetCaption(caption string) {
	self.gui.setCaption(caption)
	self.waitTillIdle()
}

func (self *GuiDriver) SetCaptionPrefix(prefix string) {
	self.gui.setCaptionPrefix(prefix)
	self.waitTillIdle()
}

func (self *GuiDriver) NextToast() *string {
	select {
	case t := <-self.toastChan:
		return &t
	default:
		return nil
	}
}

func (self *GuiDriver) Headless() bool {
	return self.headless
}
