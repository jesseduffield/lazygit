package components

import (
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/stretchr/testify/assert"
)

// this file is for testing our test code (meta, I know)

type fakeGuiDriver struct {
	failureMessage string
	pressedKeys    []string
}

var _ integrationTypes.GuiDriver = &fakeGuiDriver{}

func (self *fakeGuiDriver) PressKey(key string) {
	self.pressedKeys = append(self.pressedKeys, key)
}

func (self *fakeGuiDriver) Keys() config.KeybindingConfig {
	return config.KeybindingConfig{}
}

func (self *fakeGuiDriver) CurrentContext() types.Context {
	return nil
}

func (self *fakeGuiDriver) ContextForView(viewName string) types.Context {
	return nil
}

func (self *fakeGuiDriver) Fail(message string) {
	self.failureMessage = message
}

func (self *fakeGuiDriver) Log(message string) {
}

func (self *fakeGuiDriver) LogUI(message string) {
}

func (self *fakeGuiDriver) CheckedOutRef() *models.Branch {
	return nil
}

func (self *fakeGuiDriver) MainView() *gocui.View {
	return nil
}

func (self *fakeGuiDriver) SecondaryView() *gocui.View {
	return nil
}

func (self *fakeGuiDriver) View(viewName string) *gocui.View {
	return nil
}

func TestManualFailure(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(t *TestDriver, keys config.KeybindingConfig) {
			t.Fail("blah")
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.Equal(t, "blah", driver.failureMessage)
}

func TestSuccess(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(t *TestDriver, keys config.KeybindingConfig) {
			t.press("a")
			t.press("b")
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.EqualValues(t, []string{"a", "b"}, driver.pressedKeys)
	assert.Equal(t, "", driver.failureMessage)
}
