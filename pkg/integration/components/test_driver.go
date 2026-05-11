package components

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/config"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type TestDriver struct {
	gui        integrationTypes.GuiDriver
	keys       config.KeybindingConfig
	inputDelay int
	*assertionHelper
	shell *Shell
}

func NewTestDriver(gui integrationTypes.GuiDriver, shell *Shell, keys config.KeybindingConfig, inputDelay int) *TestDriver {
	return &TestDriver{
		gui:             gui,
		keys:            keys,
		inputDelay:      inputDelay,
		assertionHelper: &assertionHelper{gui: gui},
		shell:           shell,
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *TestDriver) press(keyStr string) {
	self.SetCaption(fmt.Sprintf("Pressing %s", keyStr))
	self.gui.PressKey(keyStr)
	self.Wait(self.inputDelay)
}

// for use when typing or navigating, because in demos we want that to happen
// faster
func (self *TestDriver) pressFast(keyStr string) {
	self.SetCaption("")
	self.gui.PressKey(keyStr)
	self.Wait(self.inputDelay / 5)
}

func (self *TestDriver) click(x, y int) {
	self.SetCaption(fmt.Sprintf("Clicking %d, %d", x, y))
	self.gui.Click(x, y)
	self.Wait(self.inputDelay)
}

// Should only be used in specific cases where you're doing something weird!
// E.g. invoking a global keybinding from within a popup.
// You probably shouldn't use this function, and should instead go through a view like t.Views().Commit().Focus().Press(...)
func (self *TestDriver) GlobalPress(keyStr string) {
	self.press(keyStr)
}

func (self *TestDriver) typeContent(content string) {
	for _, char := range content {
		self.pressFast(string(char))
	}
}

func (self *TestDriver) Common() *Common {
	return &Common{t: self}
}

// for when you want to allow lazygit to process something before continuing
func (self *TestDriver) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

func (self *TestDriver) SetCaption(caption string) {
	self.gui.SetCaption(caption)
}

func (self *TestDriver) SetCaptionPrefix(prefix string) {
	self.gui.SetCaptionPrefix(prefix)
}

func (self *TestDriver) LogUI(message string) {
	self.gui.LogUI(message)
}

func (self *TestDriver) Log(message string) {
	self.gui.LogUI(message)
}

// allows the user to run shell commands during the test to emulate background activity
func (self *TestDriver) Shell() *Shell {
	return self.shell
}

// for making assertions on lazygit views
func (self *TestDriver) Views() *Views {
	return &Views{t: self}
}

// for interacting with popups
func (self *TestDriver) ExpectPopup() *Popup {
	return &Popup{t: self}
}

func (self *TestDriver) ExpectToast(matcher *TextMatcher) *TestDriver {
	t := self.gui.NextToast()
	if t == nil {
		self.gui.Fail("Expected toast, but didn't get one")
	} else {
		self.matchString(matcher, "Unexpected toast message",
			func() string {
				return *t
			},
		)
	}

	return self
}

func (self *TestDriver) ExpectClipboard(matcher *TextMatcher) {
	self.assertWithRetries(func() (bool, string) {
		text, err := clipboard.ReadAll()
		if err != nil {
			return false, "Error occurred when reading from clipboard: " + err.Error()
		}
		ok, _ := matcher.test(text)
		return ok, fmt.Sprintf("Expected clipboard to match %s, but got %s", matcher.name(), text)
	})
}

func (self *TestDriver) ExpectSearch() *SearchDriver {
	self.inSearch()

	return &SearchDriver{t: self}
}

func (self *TestDriver) inSearch() {
	self.assertWithRetries(func() (bool, string) {
		currentView := self.gui.CurrentContext().GetView()
		return currentView.Name() == "search", "Expected search prompt to be focused"
	})
}

// for making assertions through git itself
func (self *TestDriver) Git() *Git {
	return &Git{assertionHelper: self.assertionHelper, shell: self.shell}
}

// for making assertions on the file system
func (self *TestDriver) FileSystem() *FileSystem {
	return &FileSystem{assertionHelper: self.assertionHelper}
}

// for when you just want to fail the test yourself.
// This runs callbacks to ensure we render the error after closing the gui.
func (self *TestDriver) Fail(message string) {
	self.assertionHelper.fail(message)
}
