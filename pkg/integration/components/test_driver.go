package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type TestDriver struct {
	gui          integrationTypes.GuiDriver
	keys         config.KeybindingConfig
	pushKeyDelay int
	*assertionHelper
	shell *Shell
}

func NewTestDriver(gui integrationTypes.GuiDriver, shell *Shell, keys config.KeybindingConfig, pushKeyDelay int) *TestDriver {
	return &TestDriver{
		gui:             gui,
		keys:            keys,
		pushKeyDelay:    pushKeyDelay,
		assertionHelper: &assertionHelper{gui: gui},
		shell:           shell,
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *TestDriver) press(keyStr string) {
	self.Wait(self.pushKeyDelay)

	self.gui.PressKey(keyStr)
}

func (self *TestDriver) typeContent(content string) {
	for _, char := range content {
		self.press(string(char))
	}
}

func (self *TestDriver) Actions() *Actions {
	return &Actions{t: self}
}

// for when you want to allow lazygit to process something before continuing
func (self *TestDriver) Wait(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
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

// this will look for a list item in the current panel and if it finds it, it will
// enter the keypresses required to navigate to it.
// The test will fail if:
// - the user is not in a list item
// - no list item is found containing the given text
// - multiple list items are found containing the given text in the initial page of items
//
// NOTE: this currently assumes that ViewBufferLines returns all the lines that can be accessed.
// If this changes in future, we'll need to update this code to first attempt to find the item
// in the current page and failing that, jump to the top of the view and iterate through all of it,
// looking for the item.
func (self *TestDriver) navigateToListItem(matcher *matcher) {
	self.inListContext()

	currentContext := self.gui.CurrentContext().(types.IListContext)

	view := currentContext.GetView()

	var matchIndex int

	self.assertWithRetries(func() (bool, string) {
		matchIndex = -1
		var matches []string
		lines := view.ViewBufferLines()
		// first we look for a duplicate on the current screen. We won't bother looking beyond that though.
		for i, line := range lines {
			ok, _ := matcher.test(line)
			if ok {
				matches = append(matches, line)
				matchIndex = i
			}
		}
		if len(matches) > 1 {
			return false, fmt.Sprintf("Found %d matches for `%s`, expected only a single match. Matching lines:\n%s", len(matches), matcher.name(), strings.Join(matches, "\n"))
		} else if len(matches) == 0 {
			return false, fmt.Sprintf("Could not find item matching: %s. Lines:\n%s", matcher.name(), strings.Join(lines, "\n"))
		} else {
			return true, ""
		}
	})

	selectedLineIdx := view.SelectedLineIdx()
	if selectedLineIdx == matchIndex {
		self.Views().current().SelectedLine(matcher)
		return
	}
	if selectedLineIdx < matchIndex {
		for i := selectedLineIdx; i < matchIndex; i++ {
			self.Views().current().SelectNextItem()
		}
		self.Views().current().SelectedLine(matcher)
		return
	} else {
		for i := selectedLineIdx; i > matchIndex; i-- {
			self.Views().current().SelectPreviousItem()
		}
		self.Views().current().SelectedLine(matcher)
		return
	}
}

func (self *TestDriver) inListContext() {
	self.assertWithRetries(func() (bool, string) {
		currentContext := self.gui.CurrentContext()
		_, ok := currentContext.(types.IListContext)
		return ok, fmt.Sprintf("Expected current context to be a list context, but got %s", currentContext.GetKey())
	})
}

// for making assertions on lazygit views
func (self *TestDriver) Views() *Views {
	return &Views{t: self}
}

// for interacting with popups
func (self *TestDriver) ExpectPopup() *Popup {
	return &Popup{t: self}
}

func (self *TestDriver) ExpectToast(matcher *matcher) {
	self.Views().AppStatus().Content(matcher)
}

func (self *TestDriver) ExpectClipboard(matcher *matcher) {
	self.assertWithRetries(func() (bool, string) {
		text, err := clipboard.ReadAll()
		if err != nil {
			return false, "Error occured when reading from clipboard: " + err.Error()
		}
		ok, _ := matcher.test(text)
		return ok, fmt.Sprintf("Expected clipboard to match %s, but got %s", matcher.name(), text)
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
